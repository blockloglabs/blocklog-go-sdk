package api

import (
	"context"

	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/blockloglabs/blocklog-go-sdk/transport"
)

type IncidentsClient struct {
	*BaseClient
}

func NewIncidentsClient(base *BaseClient) *IncidentsClient {
	return &IncidentsClient{BaseClient: base}
}

func (c *IncidentsClient) Create(ctx context.Context, req *models.IncidentCreateRequest) (*models.Incident, error) {
	var result models.Incident
	err := c.post(ctx, "/incidents", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *IncidentsClient) Get(ctx context.Context, id string) (*models.Incident, error) {
	var result models.Incident
	err := c.get(ctx, "/incidents/"+id, nil, &result)
	return &result, err
}

func (c *IncidentsClient) List(ctx context.Context) ([]models.Incident, error) {
	var result []models.Incident
	err := c.get(ctx, "/incidents", nil, &result)
	return result, err
}

func (c *IncidentsClient) Update(ctx context.Context, id string, req *models.IncidentUpdateRequest) (*models.Incident, error) {
	var result models.Incident
	err := c.patch(ctx, "/incidents/"+id, &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *IncidentsClient) Assign(ctx context.Context, id string, req *models.IncidentAssignRequest) (*models.Incident, error) {
	var result models.Incident
	err := c.post(ctx, "/incidents/"+id+"/assign", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *IncidentsClient) Resolve(ctx context.Context, id string, req *models.IncidentResolveRequest) (*models.Incident, error) {
	var result models.Incident
	err := c.post(ctx, "/incidents/"+id+"/resolve", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *IncidentsClient) Close(ctx context.Context, id string, req *models.IncidentCloseRequest) (*models.Incident, error) {
	var result models.Incident
	err := c.post(ctx, "/incidents/"+id+"/close", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *IncidentsClient) Report(ctx context.Context, id string) (*models.IncidentReport, error) {
	var result models.IncidentReport
	err := c.post(ctx, "/incidents/"+id+"/report", &transport.RequestOptions{JSON: map[string]interface{}{}}, &result)
	return &result, err
}

func (c *IncidentsClient) GetReport(ctx context.Context, id string) (*models.IncidentReport, error) {
	var result models.IncidentReport
	err := c.get(ctx, "/incidents/"+id+"/report", nil, &result)
	return &result, err
}

func (c *IncidentsClient) Annotate(ctx context.Context, id string, req *models.IncidentAnnotationRequest) (*models.IncidentAnnotation, error) {
	var result models.IncidentAnnotation
	err := c.post(ctx, "/incidents/"+id+"/annotations", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *IncidentsClient) Annotations(ctx context.Context, id string) ([]models.IncidentAnnotation, error) {
	var result []models.IncidentAnnotation
	err := c.get(ctx, "/incidents/"+id+"/annotations", nil, &result)
	return result, err
}

func (c *IncidentsClient) AddWorkspaceItem(ctx context.Context, id string, req *models.IncidentWorkspaceRequest) (*models.IncidentWorkspaceItem, error) {
	var result models.IncidentWorkspaceItem
	err := c.post(ctx, "/incidents/"+id+"/workspace", &transport.RequestOptions{JSON: req}, &result)
	return &result, err
}

func (c *IncidentsClient) WorkspaceItems(ctx context.Context, id string) ([]models.IncidentWorkspaceItem, error) {
	var result []models.IncidentWorkspaceItem
	err := c.get(ctx, "/incidents/"+id+"/workspace", nil, &result)
	return result, err
}

type IncidentHandle struct {
	*models.Incident
	client *IncidentsClient
}

func NewIncidentHandle(incident *models.Incident, client *IncidentsClient) *IncidentHandle {
	return &IncidentHandle{Incident: incident, client: client}
}

func (h *IncidentHandle) Assign(ctx context.Context, assignee string, notes *string) (*IncidentHandle, error) {
	incident, err := h.client.Assign(ctx, h.ID, &models.IncidentAssignRequest{
		Assignee: assignee,
		Notes:    notes,
	})
	if err != nil {
		return nil, err
	}
	h.Incident = incident
	return h, nil
}

func (h *IncidentHandle) Resolve(ctx context.Context, summary string, rootCause map[string]interface{}, remediationActions []map[string]interface{}) (*IncidentHandle, error) {
	incident, err := h.client.Resolve(ctx, h.ID, &models.IncidentResolveRequest{
		Summary:            summary,
		RootCause:          rootCause,
		RemediationActions: remediationActions,
	})
	if err != nil {
		return nil, err
	}
	h.Incident = incident
	return h, nil
}

func (h *IncidentHandle) Close(ctx context.Context, notes, approvalStatus string) (*IncidentHandle, error) {
	incident, err := h.client.Close(ctx, h.ID, &models.IncidentCloseRequest{
		Notes:          notes,
		ApprovalStatus: approvalStatus,
	})
	if err != nil {
		return nil, err
	}
	h.Incident = incident
	return h, nil
}

func (h *IncidentHandle) Annotate(ctx context.Context, text string, author *string) (*models.IncidentAnnotation, error) {
	return h.client.Annotate(ctx, h.ID, &models.IncidentAnnotationRequest{
		Text:   text,
		Author: author,
	})
}

func (h *IncidentHandle) AddWorkspaceItem(ctx context.Context, itemType, referenceID string, label *string) (*models.IncidentWorkspaceItem, error) {
	return h.client.AddWorkspaceItem(ctx, h.ID, &models.IncidentWorkspaceRequest{
		ItemType:    itemType,
		ReferenceID: referenceID,
		Label:       label,
	})
}

func (h *IncidentHandle) Report(ctx context.Context) (*models.IncidentReport, error) {
	return h.client.Report(ctx, h.ID)
}

func (h *IncidentHandle) Annotations(ctx context.Context) ([]models.IncidentAnnotation, error) {
	return h.client.Annotations(ctx, h.ID)
}

func (h *IncidentHandle) Workspace(ctx context.Context) ([]models.IncidentWorkspaceItem, error) {
	return h.client.WorkspaceItems(ctx, h.ID)
}

func (h *IncidentHandle) Refresh(ctx context.Context) (*IncidentHandle, error) {
	incident, err := h.client.Get(ctx, h.ID)
	if err != nil {
		return nil, err
	}
	h.Incident = incident
	return h, nil
}
