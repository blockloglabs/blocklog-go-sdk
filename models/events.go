package models

import (
	"time"

	"github.com/google/uuid"
)

type EventEnvelope struct {
	EventType       string                 `json:"event_type"`
	Payload         map[string]interface{} `json:"payload"`
	Source          string                 `json:"source,omitempty"`
	Timestamp       time.Time              `json:"timestamp,omitempty"`
	IdempotencyKey  string                 `json:"idempotency_key,omitempty"`
	TraceID         *uuid.UUID             `json:"trace_id,omitempty"`
	SessionID       *uuid.UUID             `json:"session_id,omitempty"`
	WorkflowID      *uuid.UUID             `json:"workflow_id,omitempty"`
	ParentEventID   *uuid.UUID             `json:"parent_event_id,omitempty"`
	RootEventID     *uuid.UUID             `json:"root_event_id,omitempty"`
	SpanID          string                 `json:"span_id,omitempty"`
	AttemptNo       int                    `json:"attempt_no,omitempty"`
	CausalityType   string                 `json:"causality_type,omitempty"`
	SchemaVersion   string                 `json:"schema_version,omitempty"`
	EventVersion    string                 `json:"event_version,omitempty"`
	AgentType       string                 `json:"agent_type,omitempty"`
	AgentID         string                 `json:"agent_id,omitempty"`
	AgentMetadata   map[string]interface{} `json:"agent_metadata,omitempty"`
}

type SessionContext struct {
	TraceID   uuid.UUID `json:"trace_id"`
	SessionID uuid.UUID `json:"session_id"`
	WorkflowID *uuid.UUID `json:"workflow_id,omitempty"`
	AgentID   string    `json:"agent_id,omitempty"`
	Source    string    `json:"source,omitempty"`
}

type IngestResponse struct {
	Ingested int      `json:"ingested"`
	LogIDs   []string `json:"log_ids"`
}