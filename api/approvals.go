package api

import (
	"context"

	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/blockloglabs/blocklog-go-sdk/transport"
)

type ApprovalsClient struct {
	*BaseClient
}

func NewApprovalsClient(base *BaseClient) *ApprovalsClient {
	return &ApprovalsClient{BaseClient: base}
}

func (c *ApprovalsClient) Request(ctx context.Context, req *models.ApprovalRequest) (*models.ApprovalResponse, error) {
	var result models.ApprovalResponse
	err := c.post(ctx, "/hitl/request", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *ApprovalsClient) Approve(ctx context.Context, approvalID string, authResponse map[string]interface{}) (*models.ApprovalResponse, error) {
	var result models.ApprovalResponse
	req := models.ApprovalApproveRequest{
		ApprovalID:   approvalID,
		AuthResponse: authResponse,
	}
	err := c.post(ctx, "/hitl/approve", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *ApprovalsClient) Reject(ctx context.Context, approvalID, reviewer, reason string, decisionID *string) (*models.ApprovalResponse, error) {
	var result models.ApprovalResponse
	req := models.ApprovalRejectRequest{
		ApprovalID:      approvalID,
		Reviewer:        reviewer,
		RejectionReason: reason,
		DecisionID:      decisionID,
	}
	err := c.post(ctx, "/hitl/reject", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *ApprovalsClient) Escalate(ctx context.Context, req *models.ApprovalEscalateRequest) (*models.ApprovalResponse, error) {
	var result models.ApprovalResponse
	err := c.post(ctx, "/hitl/escalate", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *ApprovalsClient) ListOverrides(ctx context.Context) ([]models.ApprovalOverride, error) {
	var result []models.ApprovalOverride
	err := c.get(ctx, "/hitl/overrides", nil, &result)
	return result, err
}

func (c *ApprovalsClient) GetOverride(ctx context.Context, id int) (*models.ApprovalOverride, error) {
	var result models.ApprovalOverride
	err := c.get(ctx, "/hitl/overrides/"+string(rune(id)), nil, &result)
	return &result, err
}

func (c *ApprovalsClient) AuditTrail(ctx context.Context) ([]models.ApprovalAuditTrailEntry, error) {
	var result []models.ApprovalAuditTrailEntry
	err := c.get(ctx, "/hitl/audit-trail", nil, &result)
	return result, err
}

func (c *ApprovalsClient) Status(ctx context.Context, id string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.get(ctx, "/hitl/"+id+"/status", nil, &result)
	return result, err
}

func (c *ApprovalsClient) List(ctx context.Context, params map[string]string) ([]models.ApprovalAuditTrailEntry, error) {
	var result []models.ApprovalAuditTrailEntry
	err := c.get(ctx, "/hitl/audit-trail", &transport.RequestOptions{Params: params}, &result)
	return result, err
}
