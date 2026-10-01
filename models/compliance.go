package models

import (
	"time"
)

type ComplianceReport struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	ReportKind  string                 `json:"report_kind"`
	ScopeType   string                 `json:"scope_type"`
	Framework   *string                `json:"framework,omitempty"`
	TraceID     *string                `json:"trace_id,omitempty"`
	DateFrom    *time.Time             `json:"date_from,omitempty"`
	DateTo      *time.Time             `json:"date_to,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Status      string                 `json:"status"`
	Content     map[string]interface{} `json:"content,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type ComplianceGenerateRequest struct {
	Title       string                 `json:"title"`
	ReportKind  string                 `json:"report_kind"`
	ScopeType   string                 `json:"scope_type"`
	TraceID     *string                `json:"trace_id,omitempty"`
	Framework   *string                `json:"framework,omitempty"`
	DateFrom    *time.Time             `json:"date_from,omitempty"`
	DateTo      *time.Time             `json:"date_to,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type ComplianceListResponse struct {
	Items []ComplianceReport `json:"items"`
}

type ComplianceDashboard struct {
	TotalReports    int     `json:"total_reports"`
	ByFramework     map[string]int `json:"by_framework"`
	ByStatus        map[string]int `json:"by_status"`
	RecentReports   []ComplianceReport `json:"recent_reports"`
}

type ComplianceShareRequest struct {
	Recipients            []string `json:"recipients,omitempty"`
	RecipientEmail        *string  `json:"recipient_email,omitempty"`
	ExpiresIn             *int     `json:"expires_in,omitempty"`
	ExpiresInDays         *int     `json:"expires_in_days,omitempty"`
	Note                  *string  `json:"note,omitempty"`
	CreateAuditorAPIKey   bool     `json:"create_auditor_api_key,omitempty"`
}

type ComplianceShareResponse struct {
	ShareURL      string    `json:"share_url"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	AuditorAPIKey *string   `json:"auditor_api_key,omitempty"`
}

type ComplianceExportResponse struct {
	Content     map[string]interface{} `json:"content"`
	DownloadURL *string                `json:"download_url,omitempty"`
}