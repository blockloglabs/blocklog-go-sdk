package models

import (
	"time"

	"github.com/google/uuid"
)

type Execution struct {
	ID          string                 `json:"id"`
	TraceID     uuid.UUID              `json:"trace_id"`
	SessionID   uuid.UUID              `json:"session_id"`
	Status      string                 `json:"status"`
	Input       map[string]interface{} `json:"input,omitempty"`
	Output      map[string]interface{} `json:"output,omitempty"`
	Error       *string                `json:"error,omitempty"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type ExecutionListParams struct {
	TraceID   *uuid.UUID `json:"trace_id,omitempty"`
	SessionID *uuid.UUID `json:"session_id,omitempty"`
	Status    *string    `json:"status,omitempty"`
	Limit     int        `json:"limit,omitempty"`
}

type ExecutionCreateRequest struct {
	TraceID   uuid.UUID              `json:"trace_id"`
	SessionID uuid.UUID              `json:"session_id"`
	Input     map[string]interface{} `json:"input,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type ExecutionGatewayRiskRequest struct {
	ExecutionID string                 `json:"execution_id"`
	RiskFactors map[string]interface{} `json:"risk_factors"`
}

type ExecutionGatewayAuthRequest struct {
	ExecutionID string                 `json:"execution_id"`
	Permissions []string               `json:"permissions"`
}

type ExecutionGatewayConsumptionRequest struct {
	ExecutionID string  `json:"execution_id"`
	Tokens      int     `json:"tokens"`
	Model       string  `json:"model"`
}