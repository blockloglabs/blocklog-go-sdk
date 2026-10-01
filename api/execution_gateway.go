package api

import (
	"context"

	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/blockloglabs/blocklog-go-sdk/transport"
)

type ExecutionGatewayClient struct {
	*BaseClient
}

func NewExecutionGatewayClient(base *BaseClient) *ExecutionGatewayClient {
	return &ExecutionGatewayClient{BaseClient: base}
}

func (c *ExecutionGatewayClient) RiskAssessment(ctx context.Context, req *models.ExecutionGatewayRiskRequest) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.post(ctx, "/execution-gateway/risk", &transport.RequestOptions{JSON: req}, &result)
	return result, err
}

func (c *ExecutionGatewayClient) Authorization(ctx context.Context, req *models.ExecutionGatewayAuthRequest) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.post(ctx, "/execution-gateway/auth", &transport.RequestOptions{JSON: req}, &result)
	return result, err
}

func (c *ExecutionGatewayClient) Consumption(ctx context.Context, req *models.ExecutionGatewayConsumptionRequest) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.post(ctx, "/execution-gateway/consumption", &transport.RequestOptions{JSON: req}, &result)
	return result, err
}
