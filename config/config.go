package config

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultEndpoint       = "https://blocklogsecurity.com/api/v1"
	defaultBatchSize      = 100
	defaultFlushInterval  = 2 * time.Second
	defaultTimeout        = 10 * time.Second
	defaultRetryCount     = 3
	defaultSigningAlg     = "ed25519"
	defaultEnableSigning  = true
	defaultEnableCompress = true
	defaultDebug          = false
)

type Config struct {
	APIKey           string
	AccessToken      string
	Endpoint         string
	BatchSize        int
	FlushInterval    time.Duration
	Timeout          time.Duration
	RetryCount       int
	EnableSigning    bool
	EnableCompression bool
	PersistenceEnabled bool
	Debug            bool
	SigningKey       string
	SigningAlg       string
}

func NewConfig() *Config {
	return &Config{
		Endpoint:          defaultEndpoint,
		BatchSize:         defaultBatchSize,
		FlushInterval:     defaultFlushInterval,
		Timeout:           defaultTimeout,
		RetryCount:        defaultRetryCount,
		EnableSigning:     defaultEnableSigning,
		EnableCompression: defaultEnableCompress,
		Debug:             defaultDebug,
		SigningAlg:        defaultSigningAlg,
	}
}

func (c *Config) FromEnv() *Config {
	if v := os.Getenv("BLOCKLOG_API_KEY"); v != "" {
		c.APIKey = v
	}
	if v := os.Getenv("BLOCKLOG_ACCESS_TOKEN"); v != "" {
		c.AccessToken = v
	}
	if v := os.Getenv("BLOCKLOG_ENDPOINT"); v != "" {
		c.Endpoint = v
	}
	if v := os.Getenv("BLOCKLOG_BATCH_SIZE"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			c.BatchSize = i
		}
	}
	if v := os.Getenv("BLOCKLOG_FLUSH_INTERVAL"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			c.FlushInterval = time.Duration(i) * time.Second
		}
	}
	if v := os.Getenv("BLOCKLOG_TIMEOUT"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			c.Timeout = time.Duration(i) * time.Second
		}
	}
	if v := os.Getenv("BLOCKLOG_RETRY_COUNT"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			c.RetryCount = i
		}
	}
	if v := os.Getenv("BLOCKLOG_ENABLE_SIGNING"); v != "" {
		c.EnableSigning = v == "true"
	}
	if v := os.Getenv("BLOCKLOG_ENABLE_COMPRESSION"); v != "" {
		c.EnableCompression = v == "true"
	}
	if v := os.Getenv("BLOCKLOG_DEBUG"); v != "" {
		c.Debug = v == "true"
	}
	if v := os.Getenv("BLOCKLOG_SIGNING_KEY"); v != "" {
		c.SigningKey = v
	}
	if v := os.Getenv("BLOCKLOG_SIGNING_ALG"); v != "" {
		c.SigningAlg = v
	}
	return c
}

func (c *Config) Validate() error {
	if c.APIKey == "" && c.AccessToken == "" {
		return ErrMissingCredentials
	}
	return nil
}