package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/blockloglabs/blocklog-go-sdk/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestDecisionsClient(t *testing.T, handler http.HandlerFunc) *DecisionsClient {
	server := httptest.NewServer(handler)
	t.Cleanup(func() { server.Close() })

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	return NewDecisionsClient(base)
}

func TestDecisionsClient_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/decisions", r.URL.Path)

		var req models.DecisionCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, "TRADE_EXECUTE", req.DecisionType)
		assert.Equal(t, "trading-bot", *req.Agent)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(models.Decision{
			ID:           "dec_123",
			DecisionType: req.DecisionType,
			Agent:        req.Agent,
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewDecisionsClient(base)

	ctx := context.Background()
	req := &models.DecisionCreateRequest{
		DecisionType: "TRADE_EXECUTE",
		Agent:        strPtr("trading-bot"),
	}
	result, err := client.Create(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, "dec_123", result.ID)
	assert.Equal(t, "TRADE_EXECUTE", result.DecisionType)
}

func TestDecisionsClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v1/decisions/dec_123", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.Decision{
			ID:           "dec_123",
			DecisionType: "TRADE_EXECUTE",
			Agent:        strPtr("trading-bot"),
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewDecisionsClient(base)

	ctx := context.Background()
	result, err := client.Get(ctx, "dec_123")

	require.NoError(t, err)
	assert.Equal(t, "dec_123", result.ID)
}

func TestDecisionsClient_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v1/decisions", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]models.Decision{
			{ID: "dec_1", DecisionType: "TYPE_A"},
			{ID: "dec_2", DecisionType: "TYPE_B"},
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewDecisionsClient(base)

	ctx := context.Background()
	result, err := client.List(ctx)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "dec_1", result[0].ID)
	assert.Equal(t, "dec_2", result[1].ID)
}

func TestDecisionsClient_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/api/v1/decisions/dec_123", r.URL.Path)

		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, "updated", req["status"])

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.Decision{
			ID:           "dec_123",
			DecisionType: "TRADE_EXECUTE",
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewDecisionsClient(base)

	ctx := context.Background()
	result, err := client.Update(ctx, "dec_123", map[string]interface{}{"status": "updated"})

	require.NoError(t, err)
	assert.Equal(t, "dec_123", result.ID)
}

func TestDecisionsClient_Verify(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v1/decisions/dec_123/verify", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.VerifyResponse{
			Verified:  true,
			Signature: "sig_123",
			Hash:      "hash_123",
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewDecisionsClient(base)

	ctx := context.Background()
	result, err := client.Verify(ctx, "dec_123")

	require.NoError(t, err)
	assert.True(t, result.Verified)
	assert.Equal(t, "sig_123", result.Signature)
}

func TestDecisionsClient_Timeline(t *testing.T) {
	now := time.Now()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v1/decisions/dec_123/timeline", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]models.DecisionTimelineItem{
			{ID: "1", EventType: "CREATED", Timestamp: now, Payload: map[string]interface{}{}},
			{ID: "2", EventType: "UPDATED", Timestamp: now.Add(time.Minute), Payload: map[string]interface{}{}},
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewDecisionsClient(base)

	ctx := context.Background()
	result, err := client.Timeline(ctx, "dec_123")

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "CREATED", result[0].EventType)
}

func TestDecisionsClient_Evidence(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v1/decisions/dec_123/evidence", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.EvidenceBundle{
			DecisionID: "dec_123",
			Evidence:   map[string]interface{}{"logs": []map[string]interface{}{{"id": "log_1"}}},
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewDecisionsClient(base)

	ctx := context.Background()
	result, err := client.Evidence(ctx, "dec_123")

	require.NoError(t, err)
	assert.Equal(t, "dec_123", result.DecisionID)
	assert.NotNil(t, result.Evidence)
}

func TestDecisionsClient_Replay(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v1/decisions/dec_123/replay", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"replay_id": "replay_123",
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewDecisionsClient(base)

	ctx := context.Background()
	result, err := client.Replay(ctx, "dec_123")

	require.NoError(t, err)
	assert.Equal(t, "replay_123", result["replay_id"])
}

func TestDecisionsClient_Search(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v1/decisions", r.URL.Path)
		assert.Equal(t, "TYPE_A", r.URL.Query().Get("decision_type"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]models.Decision{
			{ID: "dec_1", DecisionType: "TYPE_A"},
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewDecisionsClient(base)

	ctx := context.Background()
	result, err := client.Search(ctx, map[string]interface{}{"decision_type": "TYPE_A"})

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "TYPE_A", result[0].DecisionType)
}

func strPtr(s string) *string {
	return &s
}
