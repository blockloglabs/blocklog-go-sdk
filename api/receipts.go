package api

import (
	"context"

	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/blockloglabs/blocklog-go-sdk/transport"
)

type ReceiptVerificationClient struct {
	*BaseClient
}

func NewReceiptVerificationClient(base *BaseClient) *ReceiptVerificationClient {
	return &ReceiptVerificationClient{BaseClient: base}
}

func (c *ReceiptVerificationClient) Verify(ctx context.Context, receipt string) (*models.ReceiptVerificationResult, error) {
	var result models.ReceiptVerificationResult
	err := c.post(ctx, "/receipts/verify", &transport.RequestOptions{JSON: map[string]string{"receipt": receipt}}, &result)
	return &result, err
}
