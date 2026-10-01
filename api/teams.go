package api

import (
	"context"

	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/blockloglabs/blocklog-go-sdk/transport"
)

type TeamsClient struct {
	*BaseClient
	Members *TeamMembersClient
}

func NewTeamsClient(base *BaseClient) *TeamsClient {
	return &TeamsClient{
		BaseClient: base,
		Members:    NewTeamMembersClient(base),
	}
}

func (c *TeamsClient) List(ctx context.Context) ([]models.Team, error) {
	var result []models.Team
	err := c.get(ctx, "/teams", nil, &result)
	return result, err
}

func (c *TeamsClient) Get(ctx context.Context, id string) (*models.Team, error) {
	var result models.Team
	err := c.get(ctx, "/teams/"+id, nil, &result)
	return &result, err
}

func (c *TeamsClient) Create(ctx context.Context, req *models.TeamCreateRequest) (*models.Team, error) {
	var result models.Team
	err := c.post(ctx, "/teams", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *TeamsClient) Update(ctx context.Context, id string, req *models.TeamUpdateRequest) (*models.Team, error) {
	var result models.Team
	err := c.patch(ctx, "/teams/"+id, &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *TeamsClient) Delete(ctx context.Context, id string) error {
	return c.delete(ctx, "/teams/"+id, nil)
}

func (c *TeamsClient) NotifyTest(ctx context.Context, id string) (*models.NotifyTestResponse, error) {
	var result models.NotifyTestResponse
	err := c.post(ctx, "/teams/"+id+"/notify-test", &transport.RequestOptions{JSON: map[string]interface{}{}}, &result)
	return &result, err
}

type TeamMembersClient struct {
	*BaseClient
}

func NewTeamMembersClient(base *BaseClient) *TeamMembersClient {
	return &TeamMembersClient{BaseClient: base}
}

func (c *TeamMembersClient) List(ctx context.Context, teamID string) ([]models.TeamMember, error) {
	var result []models.TeamMember
	err := c.get(ctx, "/teams/"+teamID+"/members", nil, &result)
	return result, err
}

func (c *TeamMembersClient) Add(ctx context.Context, teamID string, req *models.TeamMemberAddRequest) (*models.TeamMember, error) {
	var result models.TeamMember
	payload := map[string]interface{}{
		"is_on_call":            req.IsOnCall,
		"notification_channels": req.NotificationChannels,
	}
	if req.UserID != nil {
		payload["user_id"] = *req.UserID
	}
	if req.Email != nil {
		payload["email"] = *req.Email
	}
	if req.Role != nil {
		payload["role"] = req.Role
	}
	err := c.post(ctx, "/teams/"+teamID+"/members", &transport.RequestOptions{JSON: payload}, &result)
	return &result, err
}

func (c *TeamMembersClient) Update(ctx context.Context, teamID, memberID string, req *models.TeamMemberUpdateRequest) (*models.TeamMember, error) {
	var result models.TeamMember
	payload := make(map[string]interface{})
	if req.Role != nil {
		payload["role"] = *req.Role
	}
	if req.IsOnCall != nil {
		payload["is_on_call"] = *req.IsOnCall
	}
	if req.NotificationChannels != nil {
		payload["notification_channels"] = req.NotificationChannels
	}
	err := c.patch(ctx, "/teams/"+teamID+"/members/"+memberID, &transport.RequestOptions{JSON: payload}, &result)
	return &result, err
}

func (c *TeamMembersClient) Remove(ctx context.Context, teamID, memberID string) error {
	return c.delete(ctx, "/teams/"+teamID+"/members/"+memberID, nil)
}
