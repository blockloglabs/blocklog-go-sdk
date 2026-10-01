package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTransport(t *testing.T) {
	t.Run("with http prefix", func(t *testing.T) {
		tr := NewTransport("http://localhost:8080/api/v1", "api-key", "", 10*time.Second, false)
		assert.NotNil(t, tr)
	})

	t.Run("with https prefix", func(t *testing.T) {
		tr := NewTransport("https://api.example.com", "api-key", "", 10*time.Second, false)
		assert.NotNil(t, tr)
	})

	t.Run("without protocol prefix", func(t *testing.T) {
		tr := NewTransport("api.example.com", "api-key", "", 10*time.Second, false)
		assert.NotNil(t, tr)
	})
}

func TestTransport_Request_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/api/v1/test", r.URL.Path)
		assert.Equal(t, "Bearer api-key", r.Header.Get("Authorization"))
		assert.Equal(t, "blocklog-go/1.0.0", r.Header.Get("User-Agent"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "", 10*time.Second, false)

	ctx := context.Background()
	resp, err := tr.Request(ctx, "GET", "/test", nil)
	require.NoError(t, err)

	var result map[string]string
	err = json.Unmarshal(resp, &result)
	require.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
}

func TestTransport_Request_WithJSONBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "test-value", body["key"])

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": "123"})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "", 10*time.Second, false)

	ctx := context.Background()
	resp, err := tr.Request(ctx, "POST", "/create", &RequestOptions{
		JSON: map[string]interface{}{"key": "test-value"},
	})
	require.NoError(t, err)

	var result map[string]string
	err = json.Unmarshal(resp, &result)
	require.NoError(t, err)
	assert.Equal(t, "123", result["id"])
}

func TestTransport_Request_WithQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "5", r.URL.Query().Get("offset"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]string{})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "", 10*time.Second, false)

	ctx := context.Background()
	_, err := tr.Request(ctx, "GET", "/list", &RequestOptions{
		Params: map[string]string{
			"limit":  "10",
			"offset": "5",
		},
	})
	require.NoError(t, err)
}

func TestTransport_Request_WithCustomHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "custom-value", r.Header.Get("X-Custom-Header"))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "", 10*time.Second, false)

	ctx := context.Background()
	_, err := tr.Request(ctx, "GET", "/test", &RequestOptions{
		Headers: map[string]string{
			"X-Custom-Header": "custom-value",
		},
	})
	require.NoError(t, err)
}

func TestTransport_Request_WithAccessToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer access-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "access-token", 10*time.Second, false)

	ctx := context.Background()
	_, err := tr.Request(ctx, "GET", "/test", nil)
	require.NoError(t, err)
}

func TestTransport_Request_TokenOverride(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer override-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "access-token", 10*time.Second, false)

	ctx := context.Background()
	_, err := tr.Request(ctx, "GET", "/test", &RequestOptions{
		TokenOverride: "override-token",
	})
	require.NoError(t, err)
}

func TestTransport_Request_SkipAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "", 10*time.Second, false)

	ctx := context.Background()
	_, err := tr.Request(ctx, "GET", "/public", &RequestOptions{
		SkipAuth: true,
	})
	require.NoError(t, err)
}

func TestTransport_Request_MissingAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "", "", 10*time.Second, false)

	ctx := context.Background()
	_, err := tr.Request(ctx, "GET", "/test", nil)
	require.Error(t, err)
	assert.Equal(t, ErrMissingAuth, err)
}

func TestTransport_Request_AuthenticationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "invalid-key", "", 10*time.Second, false)

	ctx := context.Background()
	_, err := tr.Request(ctx, "GET", "/test", nil)
	require.Error(t, err)
	assert.Equal(t, ErrAuthentication, err)
}

func TestTransport_Request_AuthorizationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "forbidden"})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "", 10*time.Second, false)

	ctx := context.Background()
	_, err := tr.Request(ctx, "GET", "/test", nil)
	require.Error(t, err)
	assert.Equal(t, ErrAuthorization, err)
}

func TestTransport_Request_NotFoundError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "", 10*time.Second, false)

	ctx := context.Background()
	_, err := tr.Request(ctx, "GET", "/missing", nil)
	require.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

func TestTransport_Request_RateLimitError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{"error": "rate limited"})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "", 10*time.Second, false)

	ctx := context.Background()
	_, err := tr.Request(ctx, "GET", "/test", nil)
	require.Error(t, err)
	assert.Equal(t, ErrRateLimit, err)
}

func TestTransport_Request_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "server error"})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "", 10*time.Second, false)

	ctx := context.Background()
	_, err := tr.Request(ctx, "GET", "/test", nil)
	require.Error(t, err)
	assert.Equal(t, ErrServerError, err)
}

func TestTransport_Request_ValidationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": "validation failed"})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "", 10*time.Second, false)

	ctx := context.Background()
	_, err := tr.Request(ctx, "POST", "/create", &RequestOptions{
		JSON: map[string]interface{}{"invalid": "data"},
	})
	require.Error(t, err)
	assert.Equal(t, ErrValidation, err)
}

func TestTransport_Request_ConflictError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "conflict"})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "", 10*time.Second, false)

	ctx := context.Background()
	_, err := tr.Request(ctx, "POST", "/create", &RequestOptions{
		JSON: map[string]interface{}{"id": "existing"},
	})
	require.Error(t, err)
	assert.Equal(t, ErrConflict, err)
}

func TestTransport_Request_NetworkError(t *testing.T) {
	tr := NewTransport("http://invalid-host-that-does-not-exist.local", "api-key", "", 10*time.Millisecond, false)

	ctx := context.Background()
	_, err := tr.Request(ctx, "GET", "/test", nil)
	require.Error(t, err)
}

func TestTransport_Request_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "", 50*time.Millisecond, false)

	ctx := context.Background()
	_, err := tr.Request(ctx, "GET", "/test", nil)
	require.Error(t, err)
}

func TestTransport_Request_204NoContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "", 10*time.Second, false)

	ctx := context.Background()
	resp, err := tr.Request(ctx, "DELETE", "/test", nil)
	require.NoError(t, err)
	assert.Nil(t, resp)
}

func TestTransport_Request_DebugLogging(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "", 10*time.Second, true)

	ctx := context.Background()
	_, err := tr.Request(ctx, "GET", "/test", nil)
	require.NoError(t, err)
}

func TestTransport_HTTPMethods(t *testing.T) {
	methods := []struct {
		method string
		fn     func(*Transport, context.Context, string, *RequestOptions) ([]byte, error)
	}{
		{"GET", (*Transport).Get},
		{"POST", (*Transport).Post},
		{"PUT", (*Transport).Put},
		{"PATCH", (*Transport).Patch},
		{"DELETE", (*Transport).Delete},
	}

	for _, m := range methods {
		t.Run(m.method, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, m.method, r.Method)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{})
			}))
			defer server.Close()

			tr := NewTransport(server.URL, "api-key", "", 10*time.Second, false)
			ctx := context.Background()

			_, err := m.fn(tr, ctx, "/test", nil)
			require.NoError(t, err)
		})
	}
}

func TestTransport_BuildURL(t *testing.T) {
	tests := []struct {
		name       string
		baseURL    string
		path       string
		params     map[string]string
		expected   string
	}{
		{
			name:     "simple path",
			baseURL:  "https://api.example.com",
			path:     "/test",
			params:   nil,
			expected: "https://api.example.com/api/v1/test",
		},
		{
			name:     "path with query params",
			baseURL:  "https://api.example.com",
			path:     "/test",
			params:   map[string]string{"key": "value"},
			expected: "https://api.example.com/api/v1/test?key=value",
		},
		{
			name:     "base URL with api/v1",
			baseURL:  "https://api.example.com/api/v1",
			path:     "/test",
			params:   nil,
			expected: "https://api.example.com/api/v1/test",
		},
		{
			name:     "path already has api/v1",
			baseURL:  "https://api.example.com",
			path:     "/api/v1/test",
			params:   nil,
			expected: "https://api.example.com/api/v1/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildURL(tt.baseURL, tt.path, tt.params)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTransport_SetAccessToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer new-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{})
	}))
	defer server.Close()

	tr := NewTransport(server.URL, "api-key", "old-token", 10*time.Second, false)
	tr.SetAccessToken("new-token")

	ctx := context.Background()
	_, err := tr.Request(ctx, "GET", "/test", nil)
	require.NoError(t, err)
}