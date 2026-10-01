package api

import (
	"context"

	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/blockloglabs/blocklog-go-sdk/transport"
	"github.com/google/uuid"
)

type ReplayClient struct {
	*BaseClient
}

func NewReplayClient(base *BaseClient) *ReplayClient {
	return &ReplayClient{BaseClient: base}
}

func (c *ReplayClient) Create(ctx context.Context, req *models.ReplayCreateRequest) (*models.ReplaySession, error) {
	var result models.ReplaySession
	err := c.post(ctx, "/forensics/replays", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *ReplayClient) Get(ctx context.Context, id string) (*models.ReplaySession, error) {
	var result models.ReplaySession
	err := c.get(ctx, "/forensics/replays/"+id, nil, &result)
	return &result, err
}

func (c *ReplayClient) Timeline(ctx context.Context, id string) ([]models.ReplayTimelineItem, error) {
	var result []models.ReplayTimelineItem
	err := c.get(ctx, "/forensics/replays/"+id+"/timeline", nil, &result)
	return result, err
}

func (c *ReplayClient) RootCause(ctx context.Context, id string) (*models.ReplayRootCause, error) {
	var result models.ReplayRootCause
	err := c.get(ctx, "/forensics/replays/"+id+"/root-cause", nil, &result)
	return &result, err
}

func (c *ReplayClient) CausalGraph(ctx context.Context, id string) (*models.ReplayCausalGraph, error) {
	var result models.ReplayCausalGraph
	err := c.get(ctx, "/forensics/replays/"+id+"/causal-graph", nil, &result)
	return &result, err
}

func (c *ReplayClient) Staleness(ctx context.Context, id string) (*models.ReplayStaleness, error) {
	var result models.ReplayStaleness
	err := c.get(ctx, "/forensics/replays/"+id+"/staleness", nil, &result)
	return &result, err
}

func (c *ReplayClient) Divergence(ctx context.Context, id string) ([]models.ReplayDivergence, error) {
	var result []models.ReplayDivergence
	err := c.get(ctx, "/forensics/replays/"+id+"/divergence", nil, &result)
	return result, err
}

func (c *ReplayClient) Counterfactual(ctx context.Context, id string, req *models.ReplayCounterfactualRequest) (*models.ReplayCounterfactualResponse, error) {
	var result models.ReplayCounterfactualResponse
	err := c.post(ctx, "/forensics/replays/"+id+"/counterfactuals", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *ReplayClient) Compare(ctx context.Context, baselineID, candidateID string) (*models.ReplayComparison, error) {
	var result models.ReplayComparison
	req := models.ReplayCompareRequest{
		BaselineSessionID:  baselineID,
		CandidateSessionID: candidateID,
	}
	err := c.post(ctx, "/forensics/compare", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *ReplayClient) GetComparison(ctx context.Context, id string) (*models.ReplayComparison, error) {
	var result models.ReplayComparison
	err := c.get(ctx, "/forensics/compare/"+id, nil, &result)
	return &result, err
}

type ReplaySession struct {
	*models.ReplaySession
	client *ReplayClient
}

func NewReplaySession(session *models.ReplaySession, client *ReplayClient) *ReplaySession {
	return &ReplaySession{ReplaySession: session, client: client}
}

func (s *ReplaySession) Timeline(ctx context.Context) ([]models.ReplayTimelineItem, error) {
	return s.client.Timeline(ctx, s.ID)
}

func (s *ReplaySession) RootCause(ctx context.Context) (*models.ReplayRootCause, error) {
	return s.client.RootCause(ctx, s.ID)
}

func (s *ReplaySession) CausalGraph(ctx context.Context) (*models.ReplayCausalGraph, error) {
	return s.client.CausalGraph(ctx, s.ID)
}

func (s *ReplaySession) Staleness(ctx context.Context) (*models.ReplayStaleness, error) {
	return s.client.Staleness(ctx, s.ID)
}

func (s *ReplaySession) Divergence(ctx context.Context) ([]models.ReplayDivergence, error) {
	return s.client.Divergence(ctx, s.ID)
}

func (s *ReplaySession) Counterfactual(ctx context.Context, tokenID string, modifiedInputs map[string]interface{}) (*models.ReplayCounterfactualResponse, error) {
	return s.client.Counterfactual(ctx, s.ID, &models.ReplayCounterfactualRequest{
		TokenID:        tokenID,
		ModifiedInputs: modifiedInputs,
	})
}

func (s *ReplaySession) Compare(ctx context.Context, otherTraceID uuid.UUID) (*models.ReplayComparison, error) {
	other, err := s.client.Create(ctx, &models.ReplayCreateRequest{TraceID: otherTraceID})
	if err != nil {
		return nil, err
	}
	return s.client.Compare(ctx, s.ID, other.ID)
}
