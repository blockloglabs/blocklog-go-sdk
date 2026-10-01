package api

import (
	"context"

	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/blockloglabs/blocklog-go-sdk/transport"
	"github.com/google/uuid"
)

type DecisionsClient struct {
	*BaseClient
}

func NewDecisionsClient(base *BaseClient) *DecisionsClient {
	return &DecisionsClient{BaseClient: base}
}

func (c *DecisionsClient) Create(ctx context.Context, req *models.DecisionCreateRequest) (*models.Decision, error) {
	var result models.Decision
	err := c.post(ctx, "/decisions", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *DecisionsClient) Get(ctx context.Context, id string) (*models.Decision, error) {
	var result models.Decision
	err := c.get(ctx, "/decisions/"+id, nil, &result)
	return &result, err
}

func (c *DecisionsClient) List(ctx context.Context) ([]models.Decision, error) {
	var result []models.Decision
	err := c.get(ctx, "/decisions", nil, &result)
	return result, err
}

func (c *DecisionsClient) Update(ctx context.Context, id string, req map[string]interface{}) (*models.Decision, error) {
	var result models.Decision
	err := c.put(ctx, "/decisions/"+id, &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *DecisionsClient) Verify(ctx context.Context, id string) (*models.VerifyResponse, error) {
	var result models.VerifyResponse
	err := c.get(ctx, "/decisions/"+id+"/verify", nil, &result)
	return &result, err
}

func (c *DecisionsClient) Timeline(ctx context.Context, id string) ([]models.DecisionTimelineItem, error) {
	var result []models.DecisionTimelineItem
	err := c.get(ctx, "/decisions/"+id+"/timeline", nil, &result)
	return result, err
}

func (c *DecisionsClient) Evidence(ctx context.Context, id string) (*models.EvidenceBundle, error) {
	var result models.EvidenceBundle
	err := c.get(ctx, "/decisions/"+id+"/evidence", nil, &result)
	return &result, err
}

func (c *DecisionsClient) Replay(ctx context.Context, id string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.get(ctx, "/decisions/"+id+"/replay", nil, &result)
	return result, err
}

func (c *DecisionsClient) Search(ctx context.Context, query map[string]interface{}) ([]models.Decision, error) {
	params := make(map[string]string)
	for k, v := range query {
		params[k] = toString(v)
	}
	var result []models.Decision
	err := c.get(ctx, "/decisions", &transport.RequestOptions{Params: params}, &result)
	return result, err
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case uuid.UUID:
		return val.String()
	default:
		return ""
	}
}
