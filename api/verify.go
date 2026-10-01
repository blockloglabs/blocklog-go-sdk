package api

import (
	"context"

	"github.com/blockloglabs/blocklog-go-sdk/models"
)

type VerifyClient struct {
	*BaseClient
}

func NewVerifyClient(base *BaseClient) *VerifyClient {
	return &VerifyClient{BaseClient: base}
}

func (c *VerifyClient) Log(ctx context.Context, logID string) (*models.VerifyLogResponse, error) {
	var result models.VerifyLogResponse
	err := c.get(ctx, "/verify/log/"+logID, nil, &result)
	return &result, err
}

func (c *VerifyClient) Batch(ctx context.Context, batchID string) (*models.VerifyBatchResponse, error) {
	var result models.VerifyBatchResponse
	err := c.get(ctx, "/verify/batch/"+batchID, nil, &result)
	return &result, err
}

func (c *VerifyClient) Decision(ctx context.Context, decisionID string) (*models.VerifyDecisionResponse, error) {
	var result models.VerifyDecisionResponse
	err := c.get(ctx, "/decisions/"+decisionID+"/verify", nil, &result)
	return &result, err
}
