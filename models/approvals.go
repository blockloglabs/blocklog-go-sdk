package models

import (
	"time"
)

type ApprovalRequest struct {
	DecisionID *string              `json:"decision_id,omitempty"`
	LogID      *string              `json:"log_id,omitempty"`
	Reason     string               `json:"reason"`
	Reviewer   *string              `json:"reviewer,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type ApprovalResponse struct {
	ID        string                 `json:"id"`
	Status    string                 `json:"status"`
	Reason    string                 `json:"reason"`
	Reviewer  *string                `json:"reviewer,omitempty"`
	DecisionID *string               `json:"decision_id,omitempty"`
	LogID     *string                `json:"log_id,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type ApprovalApproveRequest struct {
	ApprovalID string                 `json:"approval_id"`
	AuthResponse map[string]interface{} `json:"auth_response"`
}

type ApprovalRejectRequest struct {
	ApprovalID       string  `json:"approval_id"`
	Reviewer         string  `json:"reviewer"`
	RejectionReason  string  `json:"rejection_reason"`
	DecisionID       *string `json:"decision_id,omitempty"`
}

type ApprovalEscalateRequest struct {
	CurrentReviewer   string `json:"current_reviewer"`
	EscalationTarget  string `json:"escalation_target"`
	EscalationReason  string `json:"escalation_reason"`
	ApprovalID        *string `json:"approval_id,omitempty"`
}

type ApprovalOverride struct {
	ID            int                    `json:"id"`
	ApprovalID    string                 `json:"approval_id"`
	OriginalDecision string              `json:"original_decision"`
	OverrideDecision  string             `json:"override_decision"`
	Reviewer      string                 `json:"reviewer"`
	Reason        string                 `json:"reason"`
	CreatedAt     time.Time              `json:"created_at"`
}

type ApprovalAuditTrailEntry struct {
	ID        int                    `json:"id"`
	Action    string                 `json:"action"`
	ApprovalID string                `json:"approval_id"`
	Reviewer  string                 `json:"reviewer"`
	Reason    string                 `json:"reason"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}