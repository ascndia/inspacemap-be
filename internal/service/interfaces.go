package service

import (
	"context"
	"inspacemap/backend/internal/models"
)

type AuthService interface {
	Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error)
	Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error)
	AcceptInvitation(ctx context.Context, req models.AcceptInviteRequest) (*models.AuthResponse, error)
}