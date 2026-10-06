package blocklog

import (
	"context"
	"sync"

	"github.com/blockloglabs/blocklog-go-sdk/config"
	"github.com/blockloglabs/blocklog-go-sdk/models"
)

var (
	globalClient *Client
	globalMu     sync.RWMutex
)

func Init(cfg *config.Config) (*Client, error) {
	client, err := NewClient(&ClientOptions{Config: cfg})
	if err != nil {
		return nil, err
	}
	globalMu.Lock()
	globalClient = client
	globalMu.Unlock()
	return client, nil
}

func InitFromEnv() (*Client, error) {
	cfg := config.NewConfig().FromEnv()
	return Init(cfg)
}

func GetGlobalClient() *Client {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalClient
}

func Health(ctx context.Context) (*HealthStatus, error) {
	client := GetGlobalClient()
	if client == nil {
		return nil, ErrNotInitialized
	}
	return client.Health(ctx)
}

func Event(ctx context.Context, eventType string, payload map[string]interface{}, options map[string]interface{}) (*models.IngestResponse, error) {
	client := GetGlobalClient()
	if client == nil {
		return nil, ErrNotInitialized
	}
	return client.Event(ctx, eventType, payload, options)
}

func Enqueue(ctx context.Context, eventType string, payload map[string]interface{}, options map[string]interface{}) (*models.IngestResponse, error) {
	client := GetGlobalClient()
	if client == nil {
		return nil, ErrNotInitialized
	}
	return client.Enqueue(ctx, eventType, payload, options)
}

func Flush(ctx context.Context) (*models.IngestResponse, error) {
	client := GetGlobalClient()
	if client == nil {
		return nil, ErrNotInitialized
	}
	resp, err := client.Flush(ctx)
	return resp, err
}

func Shutdown(ctx context.Context) error {
	client := GetGlobalClient()
	if client == nil {
		return ErrNotInitialized
	}
	return client.Shutdown(ctx)
}

var ErrNotInitialized = &InitError{"blocklog not initialized: call Init() first"}

type InitError struct {
	msg string
}

func (e *InitError) Error() string {
	return e.msg
}
