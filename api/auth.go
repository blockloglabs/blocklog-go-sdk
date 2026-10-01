package api

import (
	"context"

	"github.com/blockloglabs/blocklog-go-sdk/models"
	"github.com/blockloglabs/blocklog-go-sdk/transport"
)

type AuthClient struct {
	*BaseClient
}

func NewAuthClient(base *BaseClient) *AuthClient {
	return &AuthClient{BaseClient: base}
}

func (c *AuthClient) Signup(ctx context.Context, req *models.SignupRequest) (*models.SignupResponse, error) {
	var result models.SignupResponse
	err := c.post(ctx, "/auth/signup", &transport.RequestOptions{JSON: req, SkipAuth: true}, &result)
	return &result, err
}

func (c *AuthClient) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	var result models.LoginResponse
	err := c.post(ctx, "/auth/login", &transport.RequestOptions{JSON: req, SkipAuth: true}, &result)
	return &result, err
}

func (c *AuthClient) Me(ctx context.Context) (*models.User, error) {
	var result models.User
	err := c.get(ctx, "/auth/me", nil, &result)
	return &result, err
}

func (c *AuthClient) Refresh(ctx context.Context) (*models.TokenResponse, error) {
	var result models.TokenResponse
	err := c.post(ctx, "/auth/refresh", &transport.RequestOptions{SkipAuth: true}, &result)
	return &result, err
}

func (c *AuthClient) Logout(ctx context.Context) error {
	return c.post(ctx, "/auth/logout", &transport.RequestOptions{}, nil)
}
