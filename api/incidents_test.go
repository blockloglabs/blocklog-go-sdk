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

func TestIncidentsClient_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/incidents", r.URL.Path)

		var req models.IncidentCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, "Test Incident", req.Title)
		assert.Equal(t, models.IncidentSeverityHigh, req.Severity)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(models.Incident{
			ID:          "inc_123",
			Title:       req.Title,
			Severity:    req.Severity,
			Status:      models.IncidentStatusOpen,
			Description: req.Description,
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewIncidentsClient(base)

	ctx := context.Background()
	req := &models.IncidentCreateRequest{
		Title:       "Test Incident",
		Severity:    models.IncidentSeverityHigh,
		Description: strPtr("Test description"),
	}
	result, err := client.Create(ctx, req)

	require.NoError(t, err)
	assert.Equal(t, "inc_123", result.ID)
	assert.Equal(t, models.IncidentSeverityHigh, result.Severity)
	assert.Equal(t, models.IncidentStatusOpen, result.Status)
}

func TestIncidentsClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v1/incidents/inc_123", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.Incident{
			ID:       "inc_123",
			Title:    "Test Incident",
			Severity: models.IncidentSeverityHigh,
			Status:   models.IncidentStatusOpen,
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewIncidentsClient(base)

	ctx := context.Background()
	result, err := client.Get(ctx, "inc_123")

	require.NoError(t, err)
	assert.Equal(t, "inc_123", result.ID)
}

func TestIncidentsClient_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v1/incidents", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]models.Incident{
			{ID: "inc_1", Title: "Incident 1", Severity: models.IncidentSeverityHigh},
			{ID: "inc_2", Title: "Incident 2", Severity: models.IncidentSeverityMedium},
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewIncidentsClient(base)

	ctx := context.Background()
	result, err := client.List(ctx)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestIncidentsClient_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "/api/v1/incidents/inc_123", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.Incident{
			ID:     "inc_123",
			Title:  "Updated Incident",
			Status: models.IncidentStatusResolved,
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewIncidentsClient(base)

	ctx := context.Background()
	status := models.IncidentStatusResolved
	result, err := client.Update(ctx, "inc_123", &models.IncidentUpdateRequest{
		Status: &status,
	})

	require.NoError(t, err)
	assert.Equal(t, "inc_123", result.ID)
	assert.Equal(t, models.IncidentStatusResolved, result.Status)
}

func TestIncidentsClient_Assign(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/incidents/inc_123/assign", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.Incident{
			ID:       "inc_123",
			Title:    "Test Incident",
			Assignee: strPtr("analyst@company.com"),
			Status:   models.IncidentStatusInvestigating,
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewIncidentsClient(base)

	ctx := context.Background()
	result, err := client.Assign(ctx, "inc_123", &models.IncidentAssignRequest{
		Assignee: "analyst@company.com",
		Notes:    strPtr("Investigating"),
	})

	require.NoError(t, err)
	assert.Equal(t, "analyst@company.com", *result.Assignee)
	assert.Equal(t, models.IncidentStatusInvestigating, result.Status)
}

func TestIncidentsClient_Annotate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/incidents/inc_123/annotations", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.IncidentAnnotation{
			ID:         "ann_123",
			IncidentID: "inc_123",
			Text:       "Found root cause",
			Author:     strPtr("analyst"),
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewIncidentsClient(base)

	ctx := context.Background()
	result, err := client.Annotate(ctx, "inc_123", &models.IncidentAnnotationRequest{
		Text:   "Found root cause",
		Author: strPtr("analyst"),
	})

	require.NoError(t, err)
	assert.Equal(t, "ann_123", result.ID)
	assert.Equal(t, "Found root cause", result.Text)
}

func TestIncidentsClient_Resolve(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/incidents/inc_123/resolve", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.Incident{
			ID:                "inc_123",
			Title:             "Test Incident",
			Status:            models.IncidentStatusResolved,
			ResolutionSummary: strPtr("Fixed"),
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewIncidentsClient(base)

	ctx := context.Background()
	result, err := client.Resolve(ctx, "inc_123", &models.IncidentResolveRequest{
		Summary: "Fixed",
	})

	require.NoError(t, err)
	assert.Equal(t, models.IncidentStatusResolved, result.Status)
}

func TestIncidentsClient_Close(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/incidents/inc_123/close", r.URL.Path)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.Incident{
			ID:     "inc_123",
			Title:  "Test Incident",
			Status: models.IncidentStatusClosed,
		})
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewIncidentsClient(base)

	ctx := context.Background()
	result, err := client.Close(ctx, "inc_123", &models.IncidentCloseRequest{
		Notes:          "Reviewed",
		ApprovalStatus: "approved",
	})

	require.NoError(t, err)
	assert.Equal(t, models.IncidentStatusClosed, result.Status)
}

func TestIncidentHandle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)

		switch r.URL.Path {
		case "/api/v1/incidents/inc_123/assign":
			json.NewEncoder(w).Encode(models.Incident{
				ID:       "inc_123",
				Assignee: strPtr("analyst@company.com"),
				Status:   models.IncidentStatusInvestigating,
			})
		case "/api/v1/incidents/inc_123/resolve":
			json.NewEncoder(w).Encode(models.Incident{
				ID:                "inc_123",
				Status:            models.IncidentStatusResolved,
				ResolutionSummary: strPtr("Fixed"),
			})
		case "/api/v1/incidents/inc_123/close":
			json.NewEncoder(w).Encode(models.Incident{
				ID:     "inc_123",
				Status: models.IncidentStatusClosed,
			})
		case "/api/v1/incidents/inc_123/annotations":
			json.NewEncoder(w).Encode(models.IncidentAnnotation{
				ID:         "ann_123",
				IncidentID: "inc_123",
				Text:       "Found root cause",
				Author:     strPtr("analyst"),
			})
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", 10*time.Second, false)
	base := NewBaseClient(tr, transport.NewRetryPolicy(3))
	client := NewIncidentsClient(base)

	ctx := context.Background()
	incident := &models.Incident{ID: "inc_123"}
	handle := NewIncidentHandle(incident, client)

	// Test Assign
	handle, err := handle.Assign(ctx, "analyst@company.com", strPtr("Investigating"))
	require.NoError(t, err)
	assert.Equal(t, "analyst@company.com", *handle.Assignee)

	// Test Resolve
	handle, err = handle.Resolve(ctx, "Fixed", nil, nil)
	require.NoError(t, err)
	assert.Equal(t, models.IncidentStatusResolved, handle.Status)

	// Test Close
	handle, err = handle.Close(ctx, "Reviewed", "approved")
	require.NoError(t, err)
	assert.Equal(t, models.IncidentStatusClosed, handle.Status)

	// Test Annotate
	annotation, err := handle.Annotate(ctx, "Found root cause", strPtr("analyst"))
	require.NoError(t, err)
	assert.Equal(t, "ann_123", annotation.ID)
}
