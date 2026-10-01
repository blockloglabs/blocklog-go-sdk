package api

import (
	"context"

	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/blockloglabs/blocklog-go-sdk/transport"
)

type ExecutionsClient struct {
	*BaseClient
}

func NewExecutionsClient(base *BaseClient) *ExecutionsClient {
	return &ExecutionsClient{BaseClient: base}
}

func (c *ExecutionsClient) List(ctx context.Context, params *models.ExecutionListParams) ([]models.Execution, error) {
	var result []models.Execution
	queryParams := make(map[string]string)
	if params != nil {
		if params.TraceID != nil {
			queryParams["trace_id"] = params.TraceID.String()
		}
		if params.SessionID != nil {
			queryParams["session_id"] = params.SessionID.String()
		}
		if params.Status != nil {
			queryParams["status"] = *params.Status
		}
		if params.Limit > 0 {
			queryParams["limit"] = string(rune(params.Limit))
		}
	}
	err := c.get(ctx, "/executions", &transport.RequestOptions{Params: queryParams}, &result)
	return result, err
}

func (c *ExecutionsClient) Get(ctx context.Context, id string) (*models.Execution, error) {
	var result models.Execution
	err := c.get(ctx, "/executions/"+id, nil, &result)
	return &result, err
}

func (c *ExecutionsClient) Create(ctx context.Context, req *models.ExecutionCreateRequest) (*models.Execution, error) {
	var result models.Execution
	err := c.post(ctx, "/executions", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}
