package blocklog

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blockloglabs/blocklog-go-sdk/api"
	"github.com/blockloglabs/blocklog-go-sdk/config"
	"github.com/blockloglabs/blocklog-go-sdk/transport"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.Decisions)
	assert.NotNil(t, client.Incidents)
	assert.NotNil(t, client.Approvals)
	assert.NotNil(t, client.Traces)
	assert.NotNil(t, client.Replay)
	assert.NotNil(t, client.Compliance)
	assert.NotNil(t, client.Verify)
	assert.NotNil(t, client.Teams)
	assert.NotNil(t, client.Auth)
	assert.NotNil(t, client.Executions)
	assert.NotNil(t, client.ExecutionGateway)
	assert.NotNil(t, client.Receipts)
	assert.NotNil(t, client.Forensics)
	assert.NotNil(t, client.HITL)
	assert.Equal(t, client.Replay, client.Forensics)
	assert.Equal(t, client.Approvals, client.HITL)
}

func TestNewClient_WithCustomTransport(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"

	customTransport := transport.NewTransport("https://custom.example.com", "custom-key", "", 30*time.Second, true)

	client, err := NewClient(&ClientOptions{
		Config:    cfg,
		Transport: customTransport,
	})
	require.NoError(t, err)
	assert.Equal(t, customTransport, client.Transport())
}

func TestNewClient_WithCustomRetry(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"

	customRetry := transport.NewRetryPolicy(5)

	client, err := NewClient(&ClientOptions{
		Config: cfg,
		Retry:  customRetry,
	})
	require.NoError(t, err)
	assert.Equal(t, customRetry, client.retry)
}

func TestNewClient_FromEnv(t *testing.T) {
	// This test would require setting env vars, skip for now
	t.Skip("Requires environment variables")
}

func TestClient_Config(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"
	cfg.Endpoint = "https://test.example.com"
	cfg.Timeout = 30 * time.Second

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	returnedCfg := client.Config()
	assert.Equal(t, cfg.APIKey, returnedCfg.APIKey)
	assert.Equal(t, cfg.Endpoint, returnedCfg.Endpoint)
	assert.Equal(t, cfg.Timeout, returnedCfg.Timeout)
}

func TestClient_Transport(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	tr := client.Transport()
	assert.NotNil(t, tr)
}

func TestClient_SetAccessToken(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	client.SetAccessToken("new-access-token")
	assert.Equal(t, "new-access-token", client.Config().AccessToken)
}

func TestClient_Event(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/logs/batch", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ingested": 1, "log_ids": ["log_123"]}`))
	}))
	defer server.Close()

	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"
	cfg.Endpoint = server.URL

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	ctx := context.Background()
	resp, err := client.Event(ctx, "TEST_EVENT", map[string]interface{}{"key": "value"}, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Ingested)
	assert.Len(t, resp.LogIDs, 1)
}

func TestClient_Enqueue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/logs/batch", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ingested": 1, "log_ids": ["log_123"]}`))
	}))
	defer server.Close()

	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"
	cfg.Endpoint = server.URL
	cfg.BatchSize = 10

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	ctx := context.Background()
	resp, err := client.Enqueue(ctx, "TEST_EVENT", map[string]interface{}{"key": "value"}, nil)
	require.NoError(t, err)
	// Enqueue doesn't flush immediately unless batch is full
	assert.Equal(t, 0, resp.Ingested)
}

func TestClient_Flush(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/logs/batch", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ingested": 5, "log_ids": ["log_1", "log_2", "log_3", "log_4", "log_5"]}`))
	}))
	defer server.Close()

	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"
	cfg.Endpoint = server.URL
	cfg.BatchSize = 10

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	ctx := context.Background()
	// Enqueue some events
	for i := 0; i < 5; i++ {
		client.Enqueue(ctx, "TEST_EVENT", map[string]interface{}{"index": i}, nil)
	}

	resp, err := client.Flush(ctx)
	require.NoError(t, err)
	assert.Equal(t, 5, resp.Ingested)
	assert.Len(t, resp.LogIDs, 5)
}

func TestClient_Shutdown(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ingested": 0, "log_ids": []}`))
	}))
	defer server.Close()

	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"
	cfg.Endpoint = server.URL

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	ctx := context.Background()
	err = client.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestClient_Health(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	ctx := context.Background()
	health, err := client.Health(ctx)
	require.NoError(t, err)
	assert.NotNil(t, health)
	assert.True(t, health.Healthy)
	assert.Equal(t, 0, health.QueueDepth)
	assert.Equal(t, 0, health.PendingEvents)
	assert.True(t, health.TransportReady)
}

func TestClient_Health_WithPendingEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ingested": 0, "log_ids": []}`))
	}))
	defer server.Close()

	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"
	cfg.Endpoint = server.URL
	cfg.BatchSize = 100

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	ctx := context.Background()
	client.Enqueue(ctx, "TEST_EVENT", map[string]interface{}{"key": "value"}, nil)

	health, err := client.Health(ctx)
	require.NoError(t, err)
	assert.False(t, health.Healthy)
	assert.Equal(t, 1, health.QueueDepth)
	assert.Equal(t, 1, health.PendingEvents)
}

func TestInit(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"

	client, err := Init(cfg)
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, client, GetGlobalClient())
}

func TestInitFromEnv(t *testing.T) {
	t.Skip("Requires environment variables")
}

func TestGetGlobalClient_Nil(t *testing.T) {
	// Reset global client
	globalClient = nil

	client := GetGlobalClient()
	assert.Nil(t, client)
}

func TestEvent_NotInitialized(t *testing.T) {
	globalClient = nil
	defer func() { globalClient = nil }()

	ctx := context.Background()
	_, err := Event(ctx, "TEST", nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")
}

func TestEnqueue_NotInitialized(t *testing.T) {
	globalClient = nil
	defer func() { globalClient = nil }()

	ctx := context.Background()
	_, err := Enqueue(ctx, "TEST", nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")
}

func TestFlush_NotInitialized(t *testing.T) {
	globalClient = nil
	defer func() { globalClient = nil }()

	ctx := context.Background()
	_, err := Flush(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")
}

func TestShutdown_NotInitialized(t *testing.T) {
	globalClient = nil
	defer func() { globalClient = nil }()

	ctx := context.Background()
	err := Shutdown(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")
}

func TestHealth_NotInitialized(t *testing.T) {
	globalClient = nil
	defer func() { globalClient = nil }()

	ctx := context.Background()
	_, err := Health(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not initialized")
}

func TestInitError(t *testing.T) {
	err := ErrNotInitialized
	assert.Contains(t, err.Error(), "not initialized")
}

func TestClient_DecisionsClient(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	assert.IsType(t, &api.DecisionsClient{}, client.Decisions)
}

func TestClient_IncidentsClient(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	assert.IsType(t, &api.IncidentsClient{}, client.Incidents)
}

func TestClient_ApprovalsClient(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	assert.IsType(t, &api.ApprovalsClient{}, client.Approvals)
}

func TestClient_WithOptions(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"

	customTransport := transport.NewTransport("https://custom.example.com", "custom-key", "", 30*time.Second, true)
	customRetry := transport.NewRetryPolicy(5)

	client, err := NewClient(&ClientOptions{
		Config:    cfg,
		Transport: customTransport,
		Retry:     customRetry,
	})
	require.NoError(t, err)

	assert.Equal(t, customTransport, client.Transport())
	assert.Equal(t, customRetry, client.retry)
}

func TestClient_WithContextOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ingested": 1, "log_ids": ["log_123"]}`))
	}))
	defer server.Close()

	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"
	cfg.Endpoint = server.URL

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	ctx := context.Background()
	// Test with various options
	resp, err := client.Event(ctx, "TEST_EVENT", map[string]interface{}{"key": "value"}, map[string]interface{}{
		"source":          "test-source",
		"idempotency_key": "test-key-123",
		"trace_id":        uuid.New().String(),
		"session_id":      uuid.New().String(),
		"workflow_id":     uuid.New().String(),
		"parent_event_id": uuid.New().String(),
		"root_event_id":   uuid.New().String(),
		"span_id":         "span-123",
		"agent_type":      "test-agent",
		"agent_id":        "agent-123",
		"agent_metadata":  map[string]interface{}{"version": "1.0"},
	})
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Ingested)
}

func TestClient_MultipleFlushes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ingested": 1, "log_ids": ["log_123"]}`))
	}))
	defer server.Close()

	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"
	cfg.Endpoint = server.URL
	cfg.BatchSize = 2

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	ctx := context.Background()

	// First flush - should be empty
	resp1, err := client.Flush(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, resp1.Ingested)

	// Add one event
	client.Enqueue(ctx, "EVENT_1", map[string]interface{}{"data": "1"}, nil)

	// Second flush - should have 1 event
	resp2, err := client.Flush(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, resp2.Ingested)

	// Add two events (batch size is 2, so should auto-flush)
	client.Enqueue(ctx, "EVENT_2", map[string]interface{}{"data": "2"}, nil)
	client.Enqueue(ctx, "EVENT_3", map[string]interface{}{"data": "3"}, nil)

	// Third flush - should be empty (already flushed)
	resp3, err := client.Flush(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, resp3.Ingested)
}

func TestClient_EventBuffer(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		var reqBody map[string]interface{}
		json.NewDecoder(r.Body).Decode(&reqBody)
		logs := reqBody["logs"].([]interface{})

		if requestCount == 1 {
			// First request (auto-flush of 2 events)
			assert.Len(t, logs, 2)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ingested": 2, "log_ids": ["log_1", "log_2"]}`))
		} else {
			// Second request (flush of 1 event)
			assert.Len(t, logs, 1)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ingested": 1, "log_ids": ["log_3"]}`))
		}
	}))
	defer server.Close()

	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"
	cfg.Endpoint = server.URL
	cfg.BatchSize = 2

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	ctx := context.Background()

	// Add first event - buffer not full
	client.Enqueue(ctx, "EVENT_1", map[string]interface{}{"data": "1"}, nil)

	// Add second event - buffer full, should auto-flush
	client.Enqueue(ctx, "EVENT_2", map[string]interface{}{"data": "2"}, nil)

	// Add third event - new buffer
	client.Enqueue(ctx, "EVENT_3", map[string]interface{}{"data": "3"}, nil)

	// Flush remaining
	resp, err := client.Flush(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Ingested)
}

func TestClient_ImmediateEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ingested": 1, "log_ids": ["log_123"]}`))
	}))
	defer server.Close()

	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"
	cfg.Endpoint = server.URL

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	ctx := context.Background()

	// Event with immediate=true should flush immediately
	resp, err := client.Event(ctx, "IMMEDIATE_EVENT", map[string]interface{}{"data": "test"}, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Ingested)
}

func TestClient_Concurrency(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate some latency
		time.Sleep(5 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ingested": 1, "log_ids": ["log_123"]}`))
	}))
	defer server.Close()

	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"
	cfg.Endpoint = server.URL

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	ctx := context.Background()

	// Fire multiple concurrent events
	done := make(chan error, 10)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			_, err := client.Event(ctx, "CONCURRENT_EVENT", map[string]interface{}{"index": idx}, nil)
			done <- err
		}(i)
	}

	// Wait for all to complete
	for i := 0; i < 10; i++ {
		err := <-done
		assert.NoError(t, err)
	}
}

func TestClient_EventProcessor_SerializeEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ingested": 1, "log_ids": ["log_123"]}`))
	}))
	defer server.Close()

	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"
	cfg.Endpoint = server.URL

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	ctx := context.Background()

	// Test that all fields are properly serialized
	traceID := uuid.New()
	sessionID := uuid.New()
	workflowID := uuid.New()
	parentEventID := uuid.New()
	rootEventID := uuid.New()

	resp, err := client.Event(ctx, "FULL_EVENT", map[string]interface{}{"data": "test"}, map[string]interface{}{
		"source":          "test-source",
		"idempotency_key": "test-key-123",
		"trace_id":        traceID.String(),
		"session_id":      sessionID.String(),
		"workflow_id":     workflowID.String(),
		"parent_event_id": parentEventID.String(),
		"root_event_id":   rootEventID.String(),
		"span_id":         "span-123",
		"agent_type":      "test-agent",
		"agent_id":        "agent-123",
		"agent_metadata":  map[string]interface{}{"version": "1.0"},
		"causality_type":  "cause",
		"attempt_no":      1,
	})
	require.NoError(t, err)
	assert.Equal(t, 1, resp.Ingested)
}

func TestClient_APIClientsNotNil(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	assert.NotNil(t, client.Decisions)
	assert.NotNil(t, client.Incidents)
	assert.NotNil(t, client.Approvals)
	assert.NotNil(t, client.Traces)
	assert.NotNil(t, client.Replay)
	assert.NotNil(t, client.Compliance)
	assert.NotNil(t, client.Verify)
	assert.NotNil(t, client.Teams)
	assert.NotNil(t, client.Auth)
	assert.NotNil(t, client.Executions)
	assert.NotNil(t, client.ExecutionGateway)
	assert.NotNil(t, client.Receipts)
}

func TestClient_ConfigValidation(t *testing.T) {
	cfg := config.NewConfig()
	// No API key or access token
	_, err := NewClient(&ClientOptions{Config: cfg})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing API key")
}

func TestClient_DefaultValues(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	assert.Equal(t, 100, client.config.BatchSize)
	assert.Equal(t, 2*time.Second, client.config.FlushInterval)
	assert.Equal(t, 10*time.Second, client.config.Timeout)
	assert.Equal(t, 3, client.config.RetryCount)
	assert.True(t, client.config.EnableSigning)
	assert.True(t, client.config.EnableCompression)
	assert.False(t, client.config.Debug)
	assert.Equal(t, "ed25519", client.config.SigningAlg)
}

func TestClient_CustomDefaultValues(t *testing.T) {
	cfg := config.NewConfig()
	cfg.APIKey = "test-api-key"
	cfg.BatchSize = 50
	cfg.FlushInterval = 5 * time.Second
	cfg.Timeout = 30 * time.Second
	cfg.RetryCount = 5
	cfg.EnableSigning = false
	cfg.EnableCompression = false
	cfg.Debug = true
	cfg.SigningAlg = "custom-alg"

	client, err := NewClient(&ClientOptions{Config: cfg})
	require.NoError(t, err)

	assert.Equal(t, 50, client.config.BatchSize)
	assert.Equal(t, 5*time.Second, client.config.FlushInterval)
	assert.Equal(t, 30*time.Second, client.config.Timeout)
	assert.Equal(t, 5, client.config.RetryCount)
	assert.False(t, client.config.EnableSigning)
	assert.False(t, client.config.EnableCompression)
	assert.True(t, client.config.Debug)
	assert.Equal(t, "custom-alg", client.config.SigningAlg)
}
