package api

import (
	"context"
	"encoding/json"

	"github.com/blockloglabs/blocklog-go-sdk/transport"
)

type BaseClient struct {
	transport *transport.Transport
	retry     *transport.RetryPolicy
}

func NewBaseClient(transport *transport.Transport, retry *transport.RetryPolicy) *BaseClient {
	return &BaseClient{
		transport: transport,
		retry:     retry,
	}
}

func (c *BaseClient) request(ctx context.Context, method, path string, opts *transport.RequestOptions) ([]byte, error) {
	if method == "GET" {
		return c.retry.Run(func() ([]byte, error) {
			return c.transport.Request(ctx, method, path, opts)
		})
	}
	return c.transport.Request(ctx, method, path, opts)
}

func (c *BaseClient) get(ctx context.Context, path string, opts *transport.RequestOptions, result interface{}) error {
	data, err := c.request(ctx, "GET", path, opts)
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, result)
}

func (c *BaseClient) post(ctx context.Context, path string, opts *transport.RequestOptions, result interface{}) error {
	data, err := c.request(ctx, "POST", path, opts)
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, result)
}

func (c *BaseClient) put(ctx context.Context, path string, opts *transport.RequestOptions, result interface{}) error {
	data, err := c.request(ctx, "PUT", path, opts)
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, result)
}

func (c *BaseClient) patch(ctx context.Context, path string, opts *transport.RequestOptions, result interface{}) error {
	data, err := c.request(ctx, "PATCH", path, opts)
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, result)
}

func (c *BaseClient) delete(ctx context.Context, path string, opts *transport.RequestOptions) error {
	_, err := c.request(ctx, "DELETE", path, opts)
	return err
}
