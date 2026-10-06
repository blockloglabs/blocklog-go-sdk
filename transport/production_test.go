package transport

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRetryPolicyRunContextStopsDuringBackoff(t *testing.T) {
	policy := NewRetryPolicy(3)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { time.Sleep(10 * time.Millisecond); cancel() }()
	_, err := policy.RunContext(ctx, func() ([]byte, error) { return nil, ErrRateLimit })
	require.ErrorIs(t, err, context.Canceled)
}

func TestHTTPErrorDoesNotExposeResponseBody(t *testing.T) {
	err := mapHTTPError(418, "api_key=secret", []byte("api_key=secret"))
	require.ErrorIs(t, err, ErrTransport)
	require.NotContains(t, err.Error(), "secret")
	require.False(t, errors.Is(err, ErrAuthentication))
}
