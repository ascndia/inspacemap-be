package models

import (
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SSOLoginRequest struct {
	Email string `json:"email" validate:"required,email"` 
}

type SSOCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int         `json:"expires_in"`
	User         UserDetail  `json:"user"`
}


type SSOSetupRequest struct {
	DisplayName  string   `json:"display_name"` 
	Provider     string   `json:"provider" validate:"oneof=oidc saml"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	IssuerURL    string   `json:"issuer_url"`
	EmailDomains []string `json:"email_domains"`
}

type InviteUserRequest struct {
	Email  string `json:"email" validate:"required,email"`
	RoleID uint   `json:"role_id" validate:"required"`
}

type AcceptInviteRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"` 
	FullName string `json:"full_name" validate:"required"`
}


type CreateApiKeyRequest struct {
	Name string `json:"name" validate:"required,min=3"`
}

type ApiKeyResponse struct {
	Key       uuid.UUID `json:"api_key"` // The secret key
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}