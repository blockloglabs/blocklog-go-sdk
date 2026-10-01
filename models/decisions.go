package models

import (
	"time"

	"github.com/google/uuid"
)

type Decision struct {
	ID                   string                 `json:"id"`
	DecisionType         string                 `json:"decision_type"`
	Agent                *string                `json:"agent,omitempty"`
	AgentID              *string                `json:"agent_id,omitempty"`
	Model                *string                `json:"model,omitempty"`
	Prompt               *string                `json:"prompt,omitempty"`
	Inputs               map[string]interface{} `json:"inputs,omitempty"`
	Outputs              map[string]interface{} `json:"outputs,omitempty"`
	Tools                []map[string]interface{} `json:"tools,omitempty"`
	Policies             []map[string]interface{} `json:"policies,omitempty"`
	EvidenceLinks        []map[string]interface{} `json:"evidence_links,omitempty"`
	Status               *string                `json:"status,omitempty"`
	Asset                *string                `json:"asset,omitempty"`
	Confidence           *float64               `json:"confidence,omitempty"`
	ConfidenceScore      *float64               `json:"confidence_score,omitempty"`
	Metadata             map[string]interface{} `json:"metadata,omitempty"`
	TraceID              *uuid.UUID             `json:"trace_id,omitempty"`
	SessionID            *uuid.UUID             `json:"session_id,omitempty"`
	WorkflowID           *uuid.UUID             `json:"workflow_id,omitempty"`
	ApprovalReferences   []map[string]interface{} `json:"approval_references,omitempty"`
	Signatures           []map[string]interface{} `json:"signatures,omitempty"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
}

type DecisionCreateRequest struct {
	DecisionType         string                 `json:"decision_type"`
	Agent                *string                `json:"agent,omitempty"`
	AgentID              *string                `json:"agent_id,omitempty"`
	Model                *string                `json:"model,omitempty"`
	Prompt               *string                `json:"prompt,omitempty"`
	Inputs               map[string]interface{} `json:"inputs,omitempty"`
	Outputs              map[string]interface{} `json:"outputs,omitempty"`
	Tools                []map[string]interface{} `json:"tools,omitempty"`
	Policies             []map[string]interface{} `json:"policies,omitempty"`
	EvidenceLinks        []map[string]interface{} `json:"evidence_links,omitempty"`
	Status               *string                `json:"status,omitempty"`
	Asset                *string                `json:"asset,omitempty"`
	Confidence           *float64               `json:"confidence,omitempty"`
	ConfidenceScore      *float64               `json:"confidence_score,omitempty"`
	Metadata             map[string]interface{} `json:"metadata,omitempty"`
	TraceID              *uuid.UUID             `json:"trace_id,omitempty"`
	SessionID            *uuid.UUID             `json:"session_id,omitempty"`
	WorkflowID           *uuid.UUID             `json:"workflow_id,omitempty"`
	ApprovalReferences   []map[string]interface{} `json:"approval_references,omitempty"`
	Signatures           []map[string]interface{} `json:"signatures,omitempty"`
}

type VerifyResponse struct {
	Verified    bool   `json:"verified"`
	Signature   string `json:"signature,omitempty"`
	Hash        string `json:"hash,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

type DecisionTimelineItem struct {
	ID        string                 `json:"id"`
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

type EvidenceBundle struct {
	DecisionID string                 `json:"decision_id"`
	Evidence   map[string]interface{} `json:"evidence"`
}