package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	assert.Equal(t, defaultEndpoint, cfg.Endpoint)
	assert.Equal(t, defaultBatchSize, cfg.BatchSize)
	assert.Equal(t, defaultFlushInterval, cfg.FlushInterval)
	assert.Equal(t, defaultTimeout, cfg.Timeout)
	assert.Equal(t, defaultRetryCount, cfg.RetryCount)
	assert.Equal(t, defaultEnableSigning, cfg.EnableSigning)
	assert.Equal(t, defaultEnableCompress, cfg.EnableCompression)
	assert.Equal(t, defaultDebug, cfg.Debug)
	assert.Equal(t, defaultSigningAlg, cfg.SigningAlg)
}

func TestConfig_FromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected Config
	}{
		{
			name: "API key from env",
			envVars: map[string]string{
				"BLOCKLOG_API_KEY": "test-api-key",
			},
			expected: Config{
				APIKey:           "test-api-key",
				Endpoint:         defaultEndpoint,
				BatchSize:        defaultBatchSize,
				FlushInterval:    defaultFlushInterval,
				Timeout:          defaultTimeout,
				RetryCount:       defaultRetryCount,
				EnableSigning:    defaultEnableSigning,
				EnableCompression: defaultEnableCompress,
				Debug:            defaultDebug,
				SigningAlg:       defaultSigningAlg,
			},
		},
		{
			name: "Access token from env",
			envVars: map[string]string{
				"BLOCKLOG_ACCESS_TOKEN": "test-access-token",
			},
			expected: Config{
				AccessToken:      "test-access-token",
				Endpoint:         defaultEndpoint,
				BatchSize:        defaultBatchSize,
				FlushInterval:    defaultFlushInterval,
				Timeout:          defaultTimeout,
				RetryCount:       defaultRetryCount,
				EnableSigning:    defaultEnableSigning,
				EnableCompression: defaultEnableCompress,
				Debug:            defaultDebug,
				SigningAlg:       defaultSigningAlg,
			},
		},
		{
			name: "Custom endpoint from env",
			envVars: map[string]string{
				"BLOCKLOG_ENDPOINT": "https://custom.example.com/api/v1",
			},
			expected: Config{
				Endpoint:         "https://custom.example.com/api/v1",
				BatchSize:        defaultBatchSize,
				FlushInterval:    defaultFlushInterval,
				Timeout:          defaultTimeout,
				RetryCount:       defaultRetryCount,
				EnableSigning:    defaultEnableSigning,
				EnableCompression: defaultEnableCompress,
				Debug:            defaultDebug,
				SigningAlg:       defaultSigningAlg,
			},
		},
		{
			name: "Custom batch size from env",
			envVars: map[string]string{
				"BLOCKLOG_BATCH_SIZE": "50",
			},
			expected: Config{
				Endpoint:         defaultEndpoint,
				BatchSize:        50,
				FlushInterval:    defaultFlushInterval,
				Timeout:          defaultTimeout,
				RetryCount:       defaultRetryCount,
				EnableSigning:    defaultEnableSigning,
				EnableCompression: defaultEnableCompress,
				Debug:            defaultDebug,
				SigningAlg:       defaultSigningAlg,
			},
		},
		{
			name: "Custom timeout from env",
			envVars: map[string]string{
				"BLOCKLOG_TIMEOUT": "30",
			},
			expected: Config{
				Endpoint:         defaultEndpoint,
				BatchSize:        defaultBatchSize,
				FlushInterval:    defaultFlushInterval,
				Timeout:          30 * time.Second,
				RetryCount:       defaultRetryCount,
				EnableSigning:    defaultEnableSigning,
				EnableCompression: defaultEnableCompress,
				Debug:            defaultDebug,
				SigningAlg:       defaultSigningAlg,
			},
		},
		{
			name: "Custom retry count from env",
			envVars: map[string]string{
				"BLOCKLOG_RETRY_COUNT": "5",
			},
			expected: Config{
				Endpoint:         defaultEndpoint,
				BatchSize:        defaultBatchSize,
				FlushInterval:    defaultFlushInterval,
				Timeout:          defaultTimeout,
				RetryCount:       5,
				EnableSigning:    defaultEnableSigning,
				EnableCompression: defaultEnableCompress,
				Debug:            defaultDebug,
				SigningAlg:       defaultSigningAlg,
			},
		},
		{
			name: "Enable signing from env",
			envVars: map[string]string{
				"BLOCKLOG_ENABLE_SIGNING": "false",
			},
			expected: Config{
				Endpoint:         defaultEndpoint,
				BatchSize:        defaultBatchSize,
				FlushInterval:    defaultFlushInterval,
				Timeout:          defaultTimeout,
				RetryCount:       defaultRetryCount,
				EnableSigning:    false,
				EnableCompression: defaultEnableCompress,
				Debug:            defaultDebug,
				SigningAlg:       defaultSigningAlg,
			},
		},
		{
			name: "Enable compression from env",
			envVars: map[string]string{
				"BLOCKLOG_ENABLE_COMPRESSION": "false",
			},
			expected: Config{
				Endpoint:         defaultEndpoint,
				BatchSize:        defaultBatchSize,
				FlushInterval:    defaultFlushInterval,
				Timeout:          defaultTimeout,
				RetryCount:       defaultRetryCount,
				EnableSigning:    defaultEnableSigning,
				EnableCompression: false,
				Debug:            defaultDebug,
				SigningAlg:       defaultSigningAlg,
			},
		},
		{
			name: "Debug from env",
			envVars: map[string]string{
				"BLOCKLOG_DEBUG": "true",
			},
			expected: Config{
				Endpoint:         defaultEndpoint,
				BatchSize:        defaultBatchSize,
				FlushInterval:    defaultFlushInterval,
				Timeout:          defaultTimeout,
				RetryCount:       defaultRetryCount,
				EnableSigning:    defaultEnableSigning,
				EnableCompression: defaultEnableCompress,
				Debug:            true,
				SigningAlg:       defaultSigningAlg,
			},
		},
		{
			name: "Signing key from env",
			envVars: map[string]string{
				"BLOCKLOG_SIGNING_KEY": "test-signing-key",
			},
			expected: Config{
				Endpoint:         defaultEndpoint,
				BatchSize:        defaultBatchSize,
				FlushInterval:    defaultFlushInterval,
				Timeout:          defaultTimeout,
				RetryCount:       defaultRetryCount,
				EnableSigning:    defaultEnableSigning,
				EnableCompression: defaultEnableCompress,
				Debug:            defaultDebug,
				SigningAlg:       defaultSigningAlg,
				SigningKey:       "test-signing-key",
			},
		},
		{
			name: "Signing alg from env",
			envVars: map[string]string{
				"BLOCKLOG_SIGNING_ALG": "ed25519",
			},
			expected: Config{
				Endpoint:         defaultEndpoint,
				BatchSize:        defaultBatchSize,
				FlushInterval:    defaultFlushInterval,
				Timeout:          defaultTimeout,
				RetryCount:       defaultRetryCount,
				EnableSigning:    defaultEnableSigning,
				EnableCompression: defaultEnableCompress,
				Debug:            defaultDebug,
				SigningAlg:       "ed25519",
			},
		},
		{
			name: "Multiple env vars",
			envVars: map[string]string{
				"BLOCKLOG_API_KEY":          "test-key",
				"BLOCKLOG_ENDPOINT":         "https://test.example.com",
				"BLOCKLOG_BATCH_SIZE":       "200",
				"BLOCKLOG_TIMEOUT":          "60",
				"BLOCKLOG_RETRY_COUNT":      "10",
				"BLOCKLOG_DEBUG":            "true",
				"BLOCKLOG_ENABLE_SIGNING":   "false",
			},
			expected: Config{
				APIKey:           "test-key",
				Endpoint:         "https://test.example.com",
				BatchSize:        200,
				FlushInterval:    defaultFlushInterval,
				Timeout:          60 * time.Second,
				RetryCount:       10,
				EnableSigning:    false,
				EnableCompression: defaultEnableCompress,
				Debug:            true,
				SigningAlg:       defaultSigningAlg,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env vars
			for key := range tt.envVars {
				os.Unsetenv(key)
			}
			// Set env vars for test
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			cfg := NewConfig().FromEnv()

			assert.Equal(t, tt.expected.APIKey, cfg.APIKey)
			assert.Equal(t, tt.expected.AccessToken, cfg.AccessToken)
			assert.Equal(t, tt.expected.Endpoint, cfg.Endpoint)
			assert.Equal(t, tt.expected.BatchSize, cfg.BatchSize)
			assert.Equal(t, tt.expected.FlushInterval, cfg.FlushInterval)
			assert.Equal(t, tt.expected.Timeout, cfg.Timeout)
			assert.Equal(t, tt.expected.RetryCount, cfg.RetryCount)
			assert.Equal(t, tt.expected.EnableSigning, cfg.EnableSigning)
			assert.Equal(t, tt.expected.EnableCompression, cfg.EnableCompression)
			assert.Equal(t, tt.expected.Debug, cfg.Debug)
			assert.Equal(t, tt.expected.SigningKey, cfg.SigningKey)
			assert.Equal(t, tt.expected.SigningAlg, cfg.SigningAlg)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *Config
		expectError bool
	}{
		{
			name: "Valid with API key",
			cfg: &Config{
				APIKey: "test-api-key",
			},
			expectError: false,
		},
		{
			name: "Valid with access token",
			cfg: &Config{
				AccessToken: "test-access-token",
			},
			expectError: false,
		},
		{
			name: "Valid with both",
			cfg: &Config{
				APIKey:      "test-api-key",
				AccessToken: "test-access-token",
			},
			expectError: false,
		},
		{
			name: "Invalid - no credentials",
			cfg: &Config{},
			expectError: true,
		},
		{
			name: "Invalid - empty API key",
			cfg: &Config{
				APIKey: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.expectError {
				require.Error(t, err)
				assert.Equal(t, ErrMissingCredentials, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_InvalidBatchSize(t *testing.T) {
	os.Setenv("BLOCKLOG_BATCH_SIZE", "invalid")
	defer os.Unsetenv("BLOCKLOG_BATCH_SIZE")

	cfg := NewConfig().FromEnv()
	assert.Equal(t, defaultBatchSize, cfg.BatchSize)
}

func TestConfig_InvalidTimeout(t *testing.T) {
	os.Setenv("BLOCKLOG_TIMEOUT", "invalid")
	defer os.Unsetenv("BLOCKLOG_TIMEOUT")

	cfg := NewConfig().FromEnv()
	assert.Equal(t, defaultTimeout, cfg.Timeout)
}

func TestConfig_InvalidRetryCount(t *testing.T) {
	os.Setenv("BLOCKLOG_RETRY_COUNT", "invalid")
	defer os.Unsetenv("BLOCKLOG_RETRY_COUNT")

	cfg := NewConfig().FromEnv()
	assert.Equal(t, defaultRetryCount, cfg.RetryCount)
}

func TestConfig_InvalidFlushInterval(t *testing.T) {
	os.Setenv("BLOCKLOG_FLUSH_INTERVAL", "invalid")
	defer os.Unsetenv("BLOCKLOG_FLUSH_INTERVAL")

	cfg := NewConfig().FromEnv()
	assert.Equal(t, defaultFlushInterval, cfg.FlushInterval)
}