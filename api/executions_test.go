package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/blockloglabs/blocklog-go-sdk/transport"
	"github.com/stretchr/testify/require"
)

func TestExecutionsListEncodesNumericLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "25", r.URL.Query().Get("limit"))
		_, _ = w.Write([]byte("[]"))
	}))
	defer server.Close()

	tr := transport.NewTransport(server.URL, "test-key", "", time.Second, false)
	client := NewExecutionsClient(NewBaseClient(tr, transport.NewRetryPolicy(0)))
	_, err := client.List(context.Background(), &models.ExecutionListParams{Limit: 25})
	require.NoError(t, err)
}
