package models

import (
	"time"

	"github.com/google/uuid"
)

type OrgShortInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type OrganizationFilter struct {
	Name     *string `json:"name,omitempty"`
	Domain   *string `json:"domain,omitempty"`
	Slug     *string `json:"slug,omitempty"`
	Website  *string `json:"website,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type OrganizationQuery struct {
	OrganizationFilter
	Limit  *int    `json:"limit,omitempty"`
	Offset *int    `json:"offset,omitempty"`
	Sort   *string `json:"sort,omitempty"`
}

type OrganizationQueryCursor struct {
	OrganizationFilter
	Limit  *int    `json:"limit,omitempty"`
	Cursor *string `json:"cursor,omitempty"`
}

type OrganizationDetail struct {
	ID         uuid.UUID              `json:"id"`
	Name       string                 `json:"name"`
	Slug       string                 `json:"slug"`
	LogoURL    string                 `json:"logo_url"`
	Website    string                 `json:"website"`
	IsActive   bool                   `json:"is_active"`
	Settings   map[string]interface{} `json:"settings"`
	CreatedAt  time.Time              `json:"created_at"`
	VenueCount int                    `json:"venue_count,omitempty"`
	UserCount  int                    `json:"user_count,omitempty"`
}

type CreateOrganizationRequest struct {
	Name     string                 `json:"name" validate:"required,min=3"`
	Slug     string                 `json:"slug" validate:"required,alphanum,min=3"`
	LogoURL  string                 `json:"logo_url" validate:"omitempty,url"`
	Website  string                 `json:"website" validate:"omitempty,url"`
	Settings map[string]interface{} `json:"settings"`
}

type UpdateOrganizationRequest struct {
	Name     *string                `json:"name" validate:"omitempty,min=3"`
	Slug     *string                `json:"slug" validate:"omitempty,alphanum,min=3"`
	LogoURL  *string                `json:"logo_url" validate:"omitempty,url"`
	Website  *string                `json:"website" validate:"omitempty,url"`
	Settings map[string]interface{} `json:"settings"`
}
