package models

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	CompanyID string `json:"company_id"`
	IsActive bool   `json:"is_active"`
	IsAdmin  bool   `json:"is_admin"`
}

type SignupRequest struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	CompanyID    string `json:"company_id,omitempty"`
	WorkspaceName string `json:"workspace_name,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	UserID       string `json:"user_id"`
	CompanyID    string `json:"company_id"`
	TeamID       *string `json:"team_id,omitempty"`
}

type SignupResponse struct {
	User      User   `json:"user"`
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	Team      Team   `json:"team"`
}

type LoginResponse struct {
	User      User    `json:"user"`
	Token     string  `json:"token"`
	ExpiresIn int     `json:"expires_in"`
	Team      *Team   `json:"team,omitempty"`
}