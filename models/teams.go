package models

import (
	"time"
)

type TeamRole string

const (
	TeamRoleOwner  TeamRole = "owner"
	TeamRoleAdmin  TeamRole = "admin"
	TeamRoleMember TeamRole = "member"
)

type BackendTeamRole string

const (
	BackendTeamRoleOwner     BackendTeamRole = "owner"
	BackendTeamRoleLead      BackendTeamRole = "lead"
	BackendTeamRoleReviewer  BackendTeamRole = "reviewer"
	BackendTeamRoleObserver  BackendTeamRole = "observer"
)

type TeamNotificationChannel string

const (
	TeamNotificationChannelSlack       TeamNotificationChannel = "slack"
	TeamNotificationChannelMSTeams     TeamNotificationChannel = "ms_teams"
	TeamNotificationChannelPagerDuty   TeamNotificationChannel = "pagerduty"
	TeamNotificationChannelCustomWebhook TeamNotificationChannel = "custom_webhook"
	TeamNotificationChannelEmail       TeamNotificationChannel = "email"
)

type Team struct {
	ID                    string                 `json:"id"`
	CompanyID             string                 `json:"company_id"`
	Name                  string                 `json:"name"`
	Slug                  string                 `json:"slug"`
	Description           *string                `json:"description,omitempty"`
	OwnerUserID           string                 `json:"owner_user_id"`
	CreatedBy             *string                `json:"created_by,omitempty"`
	IsActive              bool                   `json:"is_active"`
	DefaultSLAMinutes     int                    `json:"default_sla_minutes"`
	CreatedAt             time.Time              `json:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at"`
	SlackWebhookURL       *string                `json:"slack_webhook_url,omitempty"`
	SlackChannel          *string                `json:"slack_channel,omitempty"`
	EmailAddresses        []string               `json:"email_addresses,omitempty"`
	PagerDutyIntegrationKey *string              `json:"pagerduty_integration_key,omitempty"`
	MSTeamsWebhookURL     *string                `json:"ms_teams_webhook_url,omitempty"`
	CustomWebhookURL      *string                `json:"custom_webhook_url,omitempty"`
	CustomWebhookHeaders  map[string]string      `json:"custom_webhook_headers,omitempty"`
	TeamMetadata          map[string]interface{} `json:"team_metadata,omitempty"`
	CurrentUserIsOwner    bool                   `json:"current_user_is_owner,omitempty"`
}

type TeamMember struct {
	ID                   string                  `json:"id"`
	TeamID               string                  `json:"team_id"`
	UserID               string                  `json:"user_id"`
	CompanyID            *string                 `json:"company_id,omitempty"`
	UserEmail            *string                 `json:"user_email,omitempty"`
	UserName             *string                 `json:"user_name,omitempty"`
	Role                 TeamRole                `json:"role"`
	BackendRole          BackendTeamRole         `json:"backend_role"`
	IsOnCall             bool                    `json:"is_on_call"`
	NotificationChannels []TeamNotificationChannel `json:"notification_channels,omitempty"`
	AddedBy              *string                 `json:"added_by,omitempty"`
	CreatedAt            *time.Time              `json:"created_at,omitempty"`
}

type TeamCreateRequest struct {
	Name                   string                 `json:"name"`
	Slug                   *string                `json:"slug,omitempty"`
	Description            *string                `json:"description,omitempty"`
	SlackWebhookURL        *string                `json:"slack_webhook_url,omitempty"`
	SlackChannel           *string                `json:"slack_channel,omitempty"`
	EmailAddresses         []string               `json:"email_addresses,omitempty"`
	PagerDutyIntegrationKey *string               `json:"pagerduty_integration_key,omitempty"`
	MSTeamsWebhookURL      *string                `json:"ms_teams_webhook_url,omitempty"`
	CustomWebhookURL       *string                `json:"custom_webhook_url,omitempty"`
	CustomWebhookHeaders   map[string]string      `json:"custom_webhook_headers,omitempty"`
	DefaultSLAMinutes      *int                   `json:"default_sla_minutes,omitempty"`
	TeamMetadata           map[string]interface{} `json:"team_metadata,omitempty"`
}

type TeamUpdateRequest struct {
	Name                   *string                `json:"name,omitempty"`
	Slug                   *string                `json:"slug,omitempty"`
	Description            *string                `json:"description,omitempty"`
	SlackWebhookURL        *string                `json:"slack_webhook_url,omitempty"`
	SlackChannel           *string                `json:"slack_channel,omitempty"`
	EmailAddresses         []string               `json:"email_addresses,omitempty"`
	PagerDutyIntegrationKey *string               `json:"pagerduty_integration_key,omitempty"`
	MSTeamsWebhookURL      *string                `json:"ms_teams_webhook_url,omitempty"`
	CustomWebhookURL       *string                `json:"custom_webhook_url,omitempty"`
	CustomWebhookHeaders   map[string]string      `json:"custom_webhook_headers,omitempty"`
	DefaultSLAMinutes      *int                   `json:"default_sla_minutes,omitempty"`
	TeamMetadata           map[string]interface{} `json:"team_metadata,omitempty"`
	IsActive               *bool                  `json:"is_active,omitempty"`
}

type TeamMemberAddRequest struct {
	UserID              *string                  `json:"user_id,omitempty"`
	Email               *string                  `json:"email,omitempty"`
	Role                *TeamRole                `json:"role,omitempty"`
	IsOnCall            *bool                    `json:"is_on_call,omitempty"`
	NotificationChannels []TeamNotificationChannel `json:"notification_channels,omitempty"`
}

type TeamMemberUpdateRequest struct {
	Role                *TeamRole                 `json:"role,omitempty"`
	IsOnCall            *bool                     `json:"is_on_call,omitempty"`
	NotificationChannels []TeamNotificationChannel `json:"notification_channels,omitempty"`
}

type NotifyTestResponse struct {
	Results map[string]bool `json:"results"`
}

func MapRole(role BackendTeamRole) TeamRole {
	switch role {
	case BackendTeamRoleOwner:
		return TeamRoleOwner
	case BackendTeamRoleLead:
		return TeamRoleAdmin
	default:
		return TeamRoleMember
	}
}

func EncodeRole(role TeamRole) BackendTeamRole {
	switch role {
	case TeamRoleOwner:
		return BackendTeamRoleOwner
	case TeamRoleAdmin:
		return BackendTeamRoleLead
	case TeamRoleMember:
		return BackendTeamRoleReviewer
	default:
		return BackendTeamRoleReviewer
	}
}