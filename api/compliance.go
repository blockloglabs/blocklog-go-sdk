package api

import (
	"context"

	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/blockloglabs/blocklog-go-sdk/transport"
)

type ComplianceClient struct {
	*BaseClient
}

func NewComplianceClient(base *BaseClient) *ComplianceClient {
	return &ComplianceClient{BaseClient: base}
}

func (c *ComplianceClient) Generate(ctx context.Context, req *models.ComplianceGenerateRequest) (*models.ComplianceReport, error) {
	var result models.ComplianceReport
	err := c.post(ctx, "/compliance/reports", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *ComplianceClient) Get(ctx context.Context, id string) (*models.ComplianceReport, error) {
	var result models.ComplianceReport
	err := c.get(ctx, "/compliance/reports/"+id, nil, &result)
	return &result, err
}

func (c *ComplianceClient) List(ctx context.Context) ([]models.ComplianceReport, error) {
	var result models.ComplianceListResponse
	err := c.get(ctx, "/compliance/reports", nil, &result)
	return result.Items, err
}

func (c *ComplianceClient) Dashboard(ctx context.Context) (*models.ComplianceDashboard, error) {
	var result models.ComplianceDashboard
	err := c.get(ctx, "/compliance/dashboard", nil, &result)
	return &result, err
}

func (c *ComplianceClient) Share(ctx context.Context, id string, req *models.ComplianceShareRequest) (*models.ComplianceShareResponse, error) {
	var result models.ComplianceShareResponse
	err := c.post(ctx, "/compliance/reports/"+id+"/share", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *ComplianceClient) Export(ctx context.Context, id string, download bool) (*models.ComplianceExportResponse, error) {
	var result models.ComplianceExportResponse
	params := map[string]string{"download": "false"}
	if download {
		params["download"] = "true"
	}
	err := c.get(ctx, "/compliance/reports/"+id+"/export", &transport.RequestOptions{Params: params}, &result)
	return &result, err
}
