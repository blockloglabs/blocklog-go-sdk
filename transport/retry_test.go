package transport

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetryPolicy_Run_Success(t *testing.T) {
	policy := NewRetryPolicy(3)

	callCount := 0
	result, err := policy.Run(func() ([]byte, error) {
		callCount++
		return []byte("success"), nil
	})

	require.NoError(t, err)
	assert.Equal(t, []byte("success"), result)
	assert.Equal(t, 1, callCount)
}

func TestRetryPolicy_Run_RetryOnRetryableError(t *testing.T) {
	policy := NewRetryPolicy(3)

	callCount := 0
	result, err := policy.Run(func() ([]byte, error) {
		callCount++
		if callCount < 3 {
			return nil, ErrRateLimit
		}
		return []byte("success"), nil
	})

	require.NoError(t, err)
	assert.Equal(t, []byte("success"), result)
	assert.Equal(t, 3, callCount)
}

func TestRetryPolicy_Run_MaxRetriesExceeded(t *testing.T) {
	policy := NewRetryPolicy(2)

	callCount := 0
	_, err := policy.Run(func() ([]byte, error) {
		callCount++
		return nil, ErrRateLimit
	})

	require.Error(t, err)
	assert.Equal(t, ErrRateLimit, err)
	assert.Equal(t, 3, callCount) // initial + 2 retries
}

func TestRetryPolicy_Run_NonRetryableError(t *testing.T) {
	policy := NewRetryPolicy(3)

	callCount := 0
	_, err := policy.Run(func() ([]byte, error) {
		callCount++
		return nil, ErrAuthentication
	})

	require.Error(t, err)
	assert.Equal(t, ErrAuthentication, err)
	assert.Equal(t, 1, callCount)
}

func TestRetryPolicy_Run_ServerErrorRetryable(t *testing.T) {
	policy := NewRetryPolicy(2)

	callCount := 0
	result, err := policy.Run(func() ([]byte, error) {
		callCount++
		if callCount < 2 {
			return nil, ErrServerError
		}
		return []byte("success"), nil
	})

	require.NoError(t, err)
	assert.Equal(t, []byte("success"), result)
	assert.Equal(t, 2, callCount)
}

func TestRetryPolicy_Run_TransportErrorRetryable(t *testing.T) {
	policy := NewRetryPolicy(2)

	callCount := 0
	result, err := policy.Run(func() ([]byte, error) {
		callCount++
		if callCount < 2 {
			return nil, ErrTransport
		}
		return []byte("success"), nil
	})

	require.NoError(t, err)
	assert.Equal(t, []byte("success"), result)
	assert.Equal(t, 2, callCount)
}

func TestRetryPolicy_Run_ContextCancellation(t *testing.T) {
	policy := NewRetryPolicy(3)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	callCount := 0
	_, err := policy.Run(func() ([]byte, error) {
		callCount++
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			time.Sleep(10 * time.Millisecond)
			return []byte("success"), nil
		}
	})

	require.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestRetryPolicy_Backoff(t *testing.T) {
	policy := NewRetryPolicy(3)

	delay0 := policy.backoff(0)
	delay1 := policy.backoff(1)
	delay2 := policy.backoff(2)

	// With jitter, delays should be roughly: 250ms, 500ms, 1000ms
	assert.GreaterOrEqual(t, delay0, 250*time.Millisecond)
	assert.LessOrEqual(t, delay0, 275*time.Millisecond) // 250 + jitter (max 25ms)

	assert.GreaterOrEqual(t, delay1, 500*time.Millisecond)
	assert.LessOrEqual(t, delay1, 550*time.Millisecond) // 500 + jitter (max 50ms)

	assert.GreaterOrEqual(t, delay2, 1000*time.Millisecond)
	assert.LessOrEqual(t, delay2, 1100*time.Millisecond) // 1000 + jitter (max 100ms)
}

func TestRetryPolicy_BackoffIncreases(t *testing.T) {
	policy := NewRetryPolicy(5)

	prevDelay := policy.backoff(0)
	for i := 1; i < 5; i++ {
		delay := policy.backoff(i)
		assert.Greater(t, delay, prevDelay, "backoff should increase exponentially")
		prevDelay = delay
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		retryable bool
	}{
		{"rate limit", ErrRateLimit, true},
		{"server error", ErrServerError, true},
		{"transport error", ErrTransport, true},
		{"authentication", ErrAuthentication, false},
		{"authorization", ErrAuthorization, false},
		{"not found", ErrNotFound, false},
		{"validation", ErrValidation, false},
		{"conflict", ErrConflict, false},
		{"nil error", nil, false},
		{"wrapped rate limit", errors.New("wrapped: " + ErrRateLimit.Error()), false}, // Not wrapped with errors.Is
		{"wrapped server error", errors.New("wrapped: " + ErrServerError.Error()), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRetryableError(tt.err)
			assert.Equal(t, tt.retryable, result)
		})
	}
}

func TestRetryPolicy_RunAsync(t *testing.T) {
	policy := NewRetryPolicy(2)

	callCount := 0
	result, err := policy.RunAsync(func() ([]byte, error) {
		callCount++
		if callCount < 2 {
			return nil, ErrRateLimit
		}
		return []byte("async success"), nil
	})

	require.NoError(t, err)
	assert.Equal(t, []byte("async success"), result)
	assert.Equal(t, 2, callCount)
}

func TestRetryPolicy_MaxRetriesZero(t *testing.T) {
	policy := NewRetryPolicy(0)

	callCount := 0
	_, err := policy.Run(func() ([]byte, error) {
		callCount++
		return nil, ErrRateLimit
	})

	require.Error(t, err)
	assert.Equal(t, 1, callCount) // Only initial attempt, no retries
}

func TestRetryPolicy_CustomBaseDelay(t *testing.T) {
	policy := &RetryPolicy{
		maxRetries: 3,
		baseDelay:  100 * time.Millisecond,
	}

	delay0 := policy.backoff(0)
	delay1 := policy.backoff(1)

	assert.GreaterOrEqual(t, delay0, 100*time.Millisecond)
	assert.LessOrEqual(t, delay0, 110*time.Millisecond)

	assert.GreaterOrEqual(t, delay1, 200*time.Millisecond)
	assert.LessOrEqual(t, delay1, 220*time.Millisecond)
}

func TestMapHTTPError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expected   error
	}{
		{"401", 401, ErrAuthentication},
		{"403", 403, ErrAuthorization},
		{"404", 404, ErrNotFound},
		{"409", 409, ErrConflict},
		{"422", 422, ErrValidation},
		{"429", 429, ErrRateLimit},
		{"500", 500, ErrServerError},
		{"502", 502, ErrServerError},
		{"503", 503, ErrServerError},
		{"504", 504, ErrServerError},
		{"400", 400, ErrTransport},
		{"418", 418, ErrTransport},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mapHTTPError(tt.statusCode, "error message", []byte("body"))
			assert.ErrorIs(t, err, tt.expected)
		})
	}
}