package models

import (
	"time"

	"github.com/google/uuid"
)

// EventType identifies a canonical agent execution event. The string values
// are intentionally extensible so new backend event types do not require an
// SDK release.
type EventType string

const (
	EventTypeExecution     EventType = "execution"
	EventTypeLLMCall       EventType = "llm_call"
	EventTypeToolCall      EventType = "tool_call"
	EventTypeDecision      EventType = "decision"
	EventTypeInput         EventType = "input"
	EventTypeOutput        EventType = "output"
	EventTypeSideEffect    EventType = "side_effect"
	EventTypeApproval      EventType = "approval"
	EventTypeAuthorization EventType = "authorization"
)

// AgentEvent is the portable event envelope shared by Blocklog SDKs. Payload
// remains intentionally open-ended: the backend owns validation for each event
// type and future event fields can be added without breaking clients.
type AgentEvent struct {
	ID            string                 `json:"id,omitempty"`
	ExecutionID   string                 `json:"execution_id,omitempty"`
	EventType     EventType              `json:"event_type"`
	Timestamp     time.Time              `json:"timestamp"`
	Sequence      int64                  `json:"sequence,omitempty"`
	ParentEventID string                 `json:"parent_event_id,omitempty"`
	Status        string                 `json:"status,omitempty"`
	Payload       map[string]interface{} `json:"payload,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Error         *EventError            `json:"error,omitempty"`
	Provenance    map[string]interface{} `json:"provenance,omitempty"`
}

// EventError describes a recorded execution failure without exposing local
// credentials or transport implementation details.
type EventError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

// ToolCall describes a tool invocation payload for an AgentEvent.
type ToolCall struct {
	Name     string                 `json:"name"`
	Input    map[string]interface{} `json:"input,omitempty"`
	Output   map[string]interface{} `json:"output,omitempty"`
	Status   string                 `json:"status,omitempty"`
	Error    *EventError            `json:"error,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// LLMCall describes model request and response metadata for an AgentEvent.
// Applications should avoid placing raw sensitive prompts in Metadata.
type LLMCall struct {
	Provider string                 `json:"provider,omitempty"`
	Model    string                 `json:"model,omitempty"`
	Input    map[string]interface{} `json:"input,omitempty"`
	Output   map[string]interface{} `json:"output,omitempty"`
	Usage    map[string]interface{} `json:"usage,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ExecutionReceipt is a backend-issued verification boundary. Receipt fields
// are kept opaque until receipt issuance is supported by the configured API.
type ExecutionReceipt struct {
	ExecutionID string                 `json:"execution_id"`
	ReceiptID   string                 `json:"receipt_id,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
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

type SessionContext struct {
	TraceID    uuid.UUID  `json:"trace_id"`
	SessionID  uuid.UUID  `json:"session_id"`
	WorkflowID *uuid.UUID `json:"workflow_id,omitempty"`
	AgentID    string     `json:"agent_id,omitempty"`
	Source     string     `json:"source,omitempty"`
}

type IngestResponse struct {
	Ingested int      `json:"ingested"`
	LogIDs   []string `json:"log_ids"`
}
