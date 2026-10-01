package blocklog

import (
	"context"

	"github.com/blockloglabs/blocklog-go-sdk/config"
	"github.com/blockloglabs/blocklog-go-sdk/models"
)

var (
	globalClient *Client
)

func Init(cfg *config.Config) (*Client, error) {
	client, err := NewClient(&ClientOptions{Config: cfg})
	if err != nil {
		return nil, err
	}
	globalClient = client
	return client, nil
}

func InitFromEnv() (*Client, error) {
	cfg := config.NewConfig().FromEnv()
	return Init(cfg)
}

func GetGlobalClient() *Client {
	return globalClient
}

func Health(ctx context.Context) (*HealthStatus, error) {
	if globalClient == nil {
		return nil, ErrNotInitialized
	}
	return globalClient.Health(ctx)
}

func Event(ctx context.Context, eventType string, payload map[string]interface{}, options map[string]interface{}) (*models.IngestResponse, error) {
	if globalClient == nil {
		return nil, ErrNotInitialized
	}
	return globalClient.Event(ctx, eventType, payload, options)
}

func Enqueue(ctx context.Context, eventType string, payload map[string]interface{}, options map[string]interface{}) (*models.IngestResponse, error) {
	if globalClient == nil {
		return nil, ErrNotInitialized
	}
	return globalClient.Enqueue(ctx, eventType, payload, options)
}

func Flush(ctx context.Context) (*models.IngestResponse, error) {
	if globalClient == nil {
		return nil, ErrNotInitialized
	}
	resp, err := globalClient.Flush(ctx)
	return resp, err
}

func Shutdown(ctx context.Context) error {
	if globalClient == nil {
		return ErrNotInitialized
	}
	return globalClient.Shutdown(ctx)
}

var ErrNotInitialized = &InitError{"blocklog not initialized: call Init() first"}

type InitError struct {
	msg string
}

func (e *InitError) Error() string {
	return e.msg
}
