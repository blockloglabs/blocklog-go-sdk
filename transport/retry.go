package transport

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var (
	ErrMissingAuth    = errors.New("missing apiKey/accessToken for authenticated request")
	ErrAuthentication = errors.New("authentication failed")
	ErrAuthorization  = errors.New("authorization failed")
	ErrNotFound       = errors.New("resource not found")
	ErrConflict       = errors.New("resource conflict")
	ErrRateLimit      = errors.New("rate limit exceeded")
	ErrServerError    = errors.New("server error")
	ErrValidation     = errors.New("validation error")
	ErrTransport      = errors.New("transport error")
)

func mapHTTPError(statusCode int, message string, body []byte) error {
	switch statusCode {
	case 401:
		return ErrAuthentication
	case 403:
		return ErrAuthorization
	case 404:
		return ErrNotFound
	case 409:
		return ErrConflict
	case 422:
		return ErrValidation
	case 429:
		return ErrRateLimit
	case 500, 502, 503, 504:
		return ErrServerError
	default:
		return fmt.Errorf("%w: http status %d", ErrTransport, statusCode)
	}
}

type RetryPolicy struct {
	maxRetries int
	baseDelay  time.Duration
}

func NewRetryPolicy(maxRetries int) *RetryPolicy {
	return &RetryPolicy{
		maxRetries: maxRetries,
		baseDelay:  250 * time.Millisecond,
	}
}

func (r *RetryPolicy) Run(fn func() ([]byte, error)) ([]byte, error) {
	return r.RunContext(context.Background(), fn)
}

// RunContext retries only retryable errors and stops promptly when ctx is done.
func (r *RetryPolicy) RunContext(ctx context.Context, fn func() ([]byte, error)) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		result, err := fn()
		if err == nil {
			return result, nil
		}

		if isRetryableError(err) {
			lastErr = err
			if attempt == r.maxRetries {
				break
			}
			delay := r.backoff(attempt)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
			continue
		}

		return nil, err
	}
	return nil, lastErr
}

func (r *RetryPolicy) RunAsync(fn func() ([]byte, error)) ([]byte, error) {
	return r.Run(fn)
}

func (r *RetryPolicy) backoff(attempt int) time.Duration {
	delay := r.baseDelay * time.Duration(1<<attempt)
	jitter := time.Duration(rand.Int63n(int64(delay / 10)))
	return delay + jitter
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrRateLimit) || errors.Is(err, ErrServerError) || errors.Is(err, ErrTransport)
}
