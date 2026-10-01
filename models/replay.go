package models

import (
	"time"

	"github.com/google/uuid"
)

type ReplaySession struct {
	ID             string                 `json:"id"`
	ReplaySessionID string                `json:"replay_session_id,omitempty"`
	TraceID        uuid.UUID              `json:"trace_id"`
	TokenID        *string                `json:"token_id,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type ReplayCreateRequest struct {
	TraceID  uuid.UUID              `json:"trace_id"`
	TokenID  *string                `json:"token_id,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type ReplayTimelineItem struct {
	ID        string                 `json:"id"`
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

type ReplayRootCause struct {
	Detected        bool   `json:"detected"`
	RootCauseType   string `json:"root_cause_type,omitempty"`
	Description     string `json:"description,omitempty"`
	Confidence      float64 `json:"confidence,omitempty"`
	Remediation     string `json:"remediation,omitempty"`
}

type ReplayCausalGraph struct {
	Nodes []CausalGraphNode `json:"nodes"`
	Edges []CausalGraphEdge `json:"edges"`
}

type CausalGraphNode struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Label    string                 `json:"label"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type CausalGraphEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
	Label  string `json:"label,omitempty"`
}

type ReplayStaleness struct {
	OverallStalenessRating string                 `json:"overall_staleness_rating"`
	Findings               []StalenessFinding     `json:"findings"`
}

type StalenessFinding struct {
	SourceID   string  `json:"source_id"`
	SourceType string  `json:"source_type"`
	Staleness  float64 `json:"staleness"`
	Details    string  `json:"details,omitempty"`
}

type ReplayDivergence struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Severity    string                 `json:"severity"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

type ReplayCounterfactualRequest struct {
	TokenID         string                 `json:"token_id"`
	ModifiedInputs  map[string]interface{} `json:"modified_inputs"`
}

type ReplayCounterfactualResponse struct {
	Result map[string]interface{} `json:"result"`
}

type ReplayCompareRequest struct {
	BaselineSessionID  string `json:"baseline_session_id"`
	CandidateSessionID string `json:"candidate_session_id"`
}

type ReplayComparison struct {
	ID           string                 `json:"id"`
	Differences  []ComparisonDifference `json:"differences"`
	CreatedAt    time.Time              `json:"created_at"`
}

type ComparisonDifference struct {
	Path        string                 `json:"path"`
	Type        string                 `json:"type"`
	Baseline    interface{}            `json:"baseline"`
	Candidate   interface{}            `json:"candidate"`
	Description string                 `json:"description,omitempty"`
}