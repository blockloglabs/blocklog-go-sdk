package blocklog

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/blockloglabs/blocklog-go-sdk/api"
	"github.com/blockloglabs/blocklog-go-sdk/config"
	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/blockloglabs/blocklog-go-sdk/transport"
	"github.com/google/uuid"
)

type Client struct {
	config    *config.Config
	transport *transport.Transport
	retry     *transport.RetryPolicy

	// API clients
	Decisions        *api.DecisionsClient
	Incidents        *api.IncidentsClient
	Approvals        *api.ApprovalsClient
	Traces           *api.TracesClient
	Replay           *api.ReplayClient
	Compliance       *api.ComplianceClient
	Verify           *api.VerifyClient
	Teams            *api.TeamsClient
	Auth             *api.AuthClient
	Executions       *api.ExecutionsClient
	ExecutionGateway *api.ExecutionGatewayClient
	Receipts         *api.ReceiptVerificationClient

	// Aliases
	Forensics *api.ReplayClient
	HITL      *api.ApprovalsClient

	// Event processing
	eventBuffer    *EventBuffer
	eventProcessor *EventProcessor
	mu             sync.Mutex
}

type ClientOptions struct {
	Config    *config.Config
	Transport *transport.Transport
	Retry     *transport.RetryPolicy
}

func NewClient(opts *ClientOptions) (*Client, error) {
	cfg := opts.Config
	if cfg == nil {
		cfg = config.NewConfig().FromEnv()
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	tr := opts.Transport
	if tr == nil {
		tr = transport.NewTransport(cfg.Endpoint, cfg.APIKey, cfg.AccessToken, cfg.Timeout, cfg.Debug)
	}

	retry := opts.Retry
	if retry == nil {
		retry = transport.NewRetryPolicy(cfg.RetryCount)
	}

	base := api.NewBaseClient(tr, retry)

	client := &Client{
		config:           cfg,
		transport:        tr,
		retry:            retry,
		Decisions:        api.NewDecisionsClient(base),
		Incidents:        api.NewIncidentsClient(base),
		Approvals:        api.NewApprovalsClient(base),
		Traces:           api.NewTracesClient(base),
		Replay:           api.NewReplayClient(base),
		Compliance:       api.NewComplianceClient(base),
		Verify:           api.NewVerifyClient(base),
		Teams:            api.NewTeamsClient(base),
		Auth:             api.NewAuthClient(base),
		Executions:       api.NewExecutionsClient(base),
		ExecutionGateway: api.NewExecutionGatewayClient(base),
		Receipts:         api.NewReceiptVerificationClient(base),
	}

	client.Forensics = client.Replay
	client.HITL = client.Approvals

	client.eventBuffer = NewEventBuffer(cfg.BatchSize)
	client.eventProcessor = NewEventProcessor(cfg, tr, client.eventBuffer)

	return client, nil
}

func (c *Client) SetAccessToken(token string) *Client {
	c.transport.SetAccessToken(token)
	c.config.AccessToken = token
	return c
}

func (c *Client) Config() *config.Config {
	return c.config
}

func (c *Client) Transport() *transport.Transport {
	return c.transport
}

func (c *Client) Event(ctx context.Context, eventType string, payload map[string]interface{}, options map[string]interface{}) (*models.IngestResponse, error) {
	opts := make(map[string]interface{})
	for k, v := range options {
		opts[k] = v
	}
	opts["immediate"] = true

	c.mu.Lock()
	defer c.mu.Unlock()

	return c.eventProcessor.ProcessEvent(ctx, eventType, payload, opts)
}

func (c *Client) Enqueue(ctx context.Context, eventType string, payload map[string]interface{}, options map[string]interface{}) (*models.IngestResponse, error) {
	opts := make(map[string]interface{})
	for k, v := range options {
		opts[k] = v
	}
	opts["noAutoFlush"] = true

	c.mu.Lock()
	defer c.mu.Unlock()

	return c.eventProcessor.ProcessEvent(ctx, eventType, payload, opts)
}

func (c *Client) Flush(ctx context.Context) (*models.IngestResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	result, err := c.eventProcessor.Flush(ctx)
	return result, err
}

func (c *Client) Shutdown(ctx context.Context) error {
	if _, err := c.Flush(ctx); err != nil {
		return err
	}

	c.eventProcessor.Shutdown()
	return nil
}

func (c *Client) Health(ctx context.Context) (*HealthStatus, error) {
	queueDepth := c.eventBuffer.Len()
	return &HealthStatus{
		Healthy:        queueDepth == 0,
		QueueDepth:     queueDepth,
		PendingEvents:  queueDepth,
		TransportReady: true,
	}, nil
}

type HealthStatus struct {
	Healthy        bool `json:"healthy"`
	QueueDepth     int  `json:"queue_depth"`
	PendingEvents  int  `json:"pending_events"`
	TransportReady bool `json:"transport_ready"`
}

type EventBuffer struct {
	batchSize int
	events    []*EventEnvelope
	mu        sync.Mutex
}

type EventEnvelope struct {
	EventType      string                 `json:"event_type"`
	Payload        map[string]interface{} `json:"payload"`
	Source         string                 `json:"source,omitempty"`
	Timestamp      time.Time              `json:"timestamp,omitempty"`
	IdempotencyKey string                 `json:"idempotency_key,omitempty"`
	TraceID        *uuid.UUID             `json:"trace_id,omitempty"`
	SessionID      *uuid.UUID             `json:"session_id,omitempty"`
	WorkflowID     *uuid.UUID             `json:"workflow_id,omitempty"`
	ParentEventID  *uuid.UUID             `json:"parent_event_id,omitempty"`
	RootEventID    *uuid.UUID             `json:"root_event_id,omitempty"`
	SpanID         string                 `json:"span_id,omitempty"`
	AttemptNo      int                    `json:"attempt_no,omitempty"`
	CausalityType  string                 `json:"causality_type,omitempty"`
	SchemaVersion  string                 `json:"schema_version,omitempty"`
	EventVersion   string                 `json:"event_version,omitempty"`
	AgentType      string                 `json:"agent_type,omitempty"`
	AgentID        string                 `json:"agent_id,omitempty"`
	AgentMetadata  map[string]interface{} `json:"agent_metadata,omitempty"`
}

func NewEventBuffer(batchSize int) *EventBuffer {
	return &EventBuffer{
		batchSize: batchSize,
		events:    make([]*EventEnvelope, 0, batchSize),
	}
}

func (b *EventBuffer) Add(event *EventEnvelope) []*EventEnvelope {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.events = append(b.events, event)
	if len(b.events) >= b.batchSize {
		return b.flushLocked()
	}
	return nil
}

func (b *EventBuffer) Flush() []*EventEnvelope {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.flushLocked()
}

func (b *EventBuffer) flushLocked() []*EventEnvelope {
	if len(b.events) == 0 {
		return nil
	}
	batch := make([]*EventEnvelope, len(b.events))
	copy(batch, b.events)
	b.events = b.events[:0]
	return batch
}

func (b *EventBuffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.events)
}

type EventProcessor struct {
	config    *config.Config
	transport *transport.Transport
	buffer    *EventBuffer
	retry     *transport.RetryPolicy
}

func NewEventProcessor(cfg *config.Config, tr *transport.Transport, buffer *EventBuffer) *EventProcessor {
	return &EventProcessor{
		config:    cfg,
		transport: tr,
		buffer:    buffer,
		retry:     transport.NewRetryPolicy(cfg.RetryCount),
	}
}

func (p *EventProcessor) ProcessEvent(ctx context.Context, eventType string, payload map[string]interface{}, options map[string]interface{}) (*models.IngestResponse, error) {
	envelope := &EventEnvelope{
		EventType:      eventType,
		Payload:        payload,
		Source:         getStringOption(options, "source", "go-sdk"),
		Timestamp:      time.Now().UTC(),
		IdempotencyKey: getStringOption(options, "idempotency_key", ""),
		AttemptNo:      getIntOption(options, "attempt_no", 1),
		CausalityType:  getStringOption(options, "causality_type", ""),
		SchemaVersion:  "1.0",
		EventVersion:   "1.0",
	}

	if traceID := getStringOption(options, "trace_id", ""); traceID != "" {
		if id, err := uuid.Parse(traceID); err == nil {
			envelope.TraceID = &id
		}
	}
	if sessionID := getStringOption(options, "session_id", ""); sessionID != "" {
		if id, err := uuid.Parse(sessionID); err == nil {
			envelope.SessionID = &id
		}
	}
	if workflowID := getStringOption(options, "workflow_id", ""); workflowID != "" {
		if id, err := uuid.Parse(workflowID); err == nil {
			envelope.WorkflowID = &id
		}
	}
	if parentEventID := getStringOption(options, "parent_event_id", ""); parentEventID != "" {
		if id, err := uuid.Parse(parentEventID); err == nil {
			envelope.ParentEventID = &id
		}
	}
	if rootEventID := getStringOption(options, "root_event_id", ""); rootEventID != "" {
		if id, err := uuid.Parse(rootEventID); err == nil {
			envelope.RootEventID = &id
		}
	}
	if spanID := getStringOption(options, "span_id", ""); spanID != "" {
		envelope.SpanID = spanID
	}
	if agentType := getStringOption(options, "agent_type", ""); agentType != "" {
		envelope.AgentType = agentType
	}
	if agentID := getStringOption(options, "agent_id", ""); agentID != "" {
		envelope.AgentID = agentID
	}
	if agentMeta := getMapOption(options, "agent_metadata"); agentMeta != nil {
		envelope.AgentMetadata = agentMeta
	}

	if envelope.IdempotencyKey == "" {
		envelope.IdempotencyKey = generateIdempotencyKey(envelope)
	}

	batch := p.buffer.Add(envelope)
	if batch != nil {
		return p.flushBatch(ctx, batch)
	}

	immediate := getBoolOption(options, "immediate", false)
	if immediate {
		batch = p.buffer.Flush()
		if batch != nil {
			return p.flushBatch(ctx, batch)
		}
	}

	return &models.IngestResponse{Ingested: 0, LogIDs: []string{}}, nil
}

func (p *EventProcessor) Flush(ctx context.Context) (*models.IngestResponse, error) {
	batch := p.buffer.Flush()
	if batch == nil {
		return &models.IngestResponse{Ingested: 0, LogIDs: []string{}}, nil
	}
	return p.flushBatch(ctx, batch)
}

func (p *EventProcessor) flushBatch(ctx context.Context, batch []*EventEnvelope) (*models.IngestResponse, error) {
	if len(batch) == 0 {
		return &models.IngestResponse{Ingested: 0, LogIDs: []string{}}, nil
	}

	logs := make([]map[string]interface{}, len(batch))
	for i, e := range batch {
		logs[i] = p.serializeEvent(e)
	}

	payload := map[string]interface{}{"logs": logs}
	data, err := p.retry.Run(func() ([]byte, error) {
		return p.transport.Request(ctx, "POST", "/logs/batch", &transport.RequestOptions{JSON: payload})
	})
	if err != nil {
		return nil, err
	}

	var result models.IngestResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *EventProcessor) serializeEvent(e *EventEnvelope) map[string]interface{} {
	return map[string]interface{}{
		"event_type":      e.EventType,
		"timestamp":       e.Timestamp.UTC().Format(time.RFC3339Nano),
		"data":            e.Payload,
		"source":          e.Source,
		"idempotency_key": e.IdempotencyKey,
		"trace_id":        e.TraceID,
		"session_id":      e.SessionID,
		"workflow_id":     e.WorkflowID,
		"parent_event_id": e.ParentEventID,
		"root_event_id":   e.RootEventID,
		"span_id":         e.SpanID,
		"attempt_no":      e.AttemptNo,
		"causality_type":  e.CausalityType,
		"schema_version":  e.SchemaVersion,
		"event_version":   e.EventVersion,
		"agent_id":        e.AgentID,
		"agent_metadata":  e.AgentMetadata,
	}
}

func (p *EventProcessor) Shutdown() {
	p.buffer.Flush()
}

func generateIdempotencyKey(e *EventEnvelope) string {
	// Simple hash for idempotency key
	data := fmt.Sprintf("%s:%s:%v:%v:%v", e.EventType, e.Source, e.TraceID, e.SessionID, e.Payload)
	// In a real implementation, use a proper hash
	return fmt.Sprintf("blk_%x", len(data))
}

func getStringOption(options map[string]interface{}, key, defaultValue string) string {
	if v, ok := options[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultValue
}

func getIntOption(options map[string]interface{}, key string, defaultValue int) int {
	if v, ok := options[key]; ok {
		if i, ok := v.(int); ok {
			return i
		}
		if f, ok := v.(float64); ok {
			return int(f)
		}
	}
	return defaultValue
}

func getBoolOption(options map[string]interface{}, key string, defaultValue bool) bool {
	if v, ok := options[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return defaultValue
}

func getMapOption(options map[string]interface{}, key string) map[string]interface{} {
	if v, ok := options[key]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			return m
		}
	}
	return nil
}
