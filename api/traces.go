package api

import (
	"context"

	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/blockloglabs/blocklog-go-sdk/transport"
)

type TracesClient struct {
	*BaseClient
}

func NewTracesClient(base *BaseClient) *TracesClient {
	return &TracesClient{BaseClient: base}
}

func (c *TracesClient) List(ctx context.Context, params *models.TraceListParams) (*models.TraceListResponse, error) {
	var result models.TraceListResponse
	queryParams := make(map[string]string)
	if params != nil {
		if params.TraceID != nil {
			queryParams["trace_id"] = params.TraceID.String()
		}
		if params.SessionID != nil {
			queryParams["session_id"] = params.SessionID.String()
		}
		if params.WorkflowID != nil {
			queryParams["workflow_id"] = params.WorkflowID.String()
		}
		if params.Source != nil {
			queryParams["source"] = *params.Source
		}
		if params.EventType != nil {
			queryParams["event_type"] = *params.EventType
		}
		if params.FromTS != nil {
			queryParams["from"] = params.FromTS.Format("2006-01-02T15:04:05Z")
		}
		if params.ToTS != nil {
			queryParams["to"] = params.ToTS.Format("2006-01-02T15:04:05Z")
		}
		if params.Limit > 0 {
			queryParams["limit"] = string(rune(params.Limit))
		}
	}
	err := c.get(ctx, "/traces", &transport.RequestOptions{Params: queryParams}, &result)
	return &result, err
}

func (c *TracesClient) Get(ctx context.Context, traceID string) (*models.Trace, error) {
	var result models.Trace
	err := c.get(ctx, "/traces/"+traceID, nil, &result)
	return &result, err
}

func (c *TracesClient) SessionTimeline(ctx context.Context, sessionID string, params *models.SessionTimelineParams) (*models.SessionTimelineResponse, error) {
	var result models.SessionTimelineResponse
	queryParams := make(map[string]string)
	if params != nil {
		if params.Cursor != nil {
			queryParams["cursor"] = *params.Cursor
		}
		if params.Limit > 0 {
			queryParams["limit"] = string(rune(params.Limit))
		}
	}
	err := c.get(ctx, "/sessions/"+sessionID+"/timeline", &transport.RequestOptions{Params: queryParams}, &result)
	return &result, err
}
