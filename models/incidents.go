package models

import (
	"time"

	"github.com/google/uuid"
)

type IncidentSeverity string

const (
	IncidentSeverityLow      IncidentSeverity = "low"
	IncidentSeverityMedium   IncidentSeverity = "medium"
	IncidentSeverityHigh     IncidentSeverity = "high"
	IncidentSeverityCritical IncidentSeverity = "critical"
)

type IncidentStatus string

const (
	IncidentStatusOpen       IncidentStatus = "open"
	IncidentStatusAssigned   IncidentStatus = "assigned"
	IncidentStatusInvestigating IncidentStatus = "investigating"
	IncidentStatusResolved   IncidentStatus = "resolved"
	IncidentStatusClosed     IncidentStatus = "closed"
)

type Incident struct {
	ID              string                 `json:"id"`
	Title           string                 `json:"title"`
	Severity        IncidentSeverity       `json:"severity"`
	Status          IncidentStatus         `json:"status"`
	Description     *string                `json:"description,omitempty"`
	TraceID         *uuid.UUID             `json:"trace_id,omitempty"`
	Assignee        *string                `json:"assignee,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	ResolvedAt      *time.Time             `json:"resolved_at,omitempty"`
	ClosedAt        *time.Time             `json:"closed_at,omitempty"`
	ResolutionSummary *string              `json:"resolution_summary,omitempty"`
	RootCause       map[string]interface{} `json:"root_cause,omitempty"`
	RemediationActions []map[string]interface{} `json:"remediation_actions,omitempty"`
}

type IncidentCreateRequest struct {
	Title       string                 `json:"title"`
	Severity    IncidentSeverity       `json:"severity"`
	TraceID     *uuid.UUID             `json:"trace_id,omitempty"`
	Description *string                `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type IncidentUpdateRequest struct {
	Title       *string                `json:"title,omitempty"`
	Severity    *IncidentSeverity      `json:"severity,omitempty"`
	Status      *IncidentStatus        `json:"status,omitempty"`
	Description *string                `json:"description,omitempty"`
	Assignee    *string                `json:"assignee,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type IncidentAssignRequest struct {
	Assignee string  `json:"assignee"`
	Notes    *string `json:"notes,omitempty"`
}

type IncidentResolveRequest struct {
	Summary             string                   `json:"resolution_summary"`
	RootCause           map[string]interface{}   `json:"root_cause,omitempty"`
	RemediationActions  []map[string]interface{} `json:"remediation_actions,omitempty"`
}

type IncidentCloseRequest struct {
	Notes            string  `json:"closure_notes,omitempty"`
	ApprovalStatus string  `json:"approval_status,omitempty"`
}

type IncidentAnnotation struct {
	ID        string    `json:"id"`
	IncidentID string   `json:"incident_id"`
	Text      string    `json:"text"`
	Author    *string   `json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type IncidentAnnotationRequest struct {
	Text   string  `json:"text"`
	Author *string `json:"author,omitempty"`
}

type IncidentWorkspaceItem struct {
	ID            string    `json:"id"`
	IncidentID    string    `json:"incident_id"`
	ItemType      string    `json:"item_type"`
	ReferenceID   string    `json:"reference_id"`
	Label         *string   `json:"label,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type IncidentWorkspaceRequest struct {
	ItemType     string  `json:"item_type"`
	ReferenceID  string  `json:"reference_id"`
	Label        *string `json:"label,omitempty"`
}

type IncidentReport struct {
	IncidentID   string                 `json:"incident_id"`
	Title        string                 `json:"title"`
	GeneratedAt  time.Time              `json:"generated_at"`
	Content      map[string]interface{} `json:"content"`
	DownloadURL  *string                `json:"download_url,omitempty"`
}