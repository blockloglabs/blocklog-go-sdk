package models

import (
	"time"

	"github.com/google/uuid"
)

type Trace struct {
	ID          string                 `json:"id"`
	TraceID     uuid.UUID              `json:"trace_id"`
	SessionID   uuid.UUID              `json:"session_id"`
	WorkflowID  *uuid.UUID             `json:"workflow_id,omitempty"`
	Source      string                 `json:"source"`
	EventType   string                 `json:"event_type"`
	Timestamp   time.Time              `json:"timestamp"`
	Payload     map[string]interface{} `json:"payload"`
}

type TraceListParams struct {
	TraceID    *uuid.UUID `json:"trace_id,omitempty"`
	SessionID  *uuid.UUID `json:"session_id,omitempty"`
	WorkflowID *uuid.UUID `json:"workflow_id,omitempty"`
	Source     *string    `json:"source,omitempty"`
	EventType  *string    `json:"event_type,omitempty"`
	FromTS     *time.Time `json:"from,omitempty"`
	ToTS       *time.Time `json:"to,omitempty"`
	Limit      int        `json:"limit,omitempty"`
}

type TraceListResponse struct {
	Items       []Trace `json:"items"`
	NextCursor  *string `json:"next_cursor,omitempty"`
}

type SessionTimelineParams struct {
	Cursor *string `json:"cursor,omitempty"`
	Limit  int     `json:"limit,omitempty"`
}

type SessionTimelineResponse struct {
	Items      []Trace `json:"items"`
	NextCursor *string `json:"next_cursor,omitempty"`
}