package blocklog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewClientAcceptsRootConfigAndDefaults(t *testing.T) {
	client, err := NewClient(Config{APIKey: "test-key"})
	require.NoError(t, err)
	require.Equal(t, 10*time.Second, client.Config().Timeout)
}

func TestClientUsesConfiguredHTTPClientAndUserAgent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "example-agent/1.0", r.UserAgent())
		_, _ = w.Write([]byte("[]"))
	}))
	defer server.Close()

	client, err := NewClient(Config{APIKey: "test-key", Endpoint: server.URL, UserAgent: "example-agent/1.0", HTTPClient: server.Client()})
	require.NoError(t, err)
	_, err = client.Executions.List(context.Background(), nil)
	require.NoError(t, err)
}
