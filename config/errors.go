package config

import "errors"

var (
	ErrMissingCredentials = errors.New("missing API key or access token")
	ErrInvalidEndpoint    = errors.New("endpoint must not be empty")
	ErrInvalidTimeout     = errors.New("timeout must be positive")
	ErrInvalidBatchSize   = errors.New("batch size must be positive")
	ErrInvalidRetryCount  = errors.New("retry count must not be negative")
)
