package models

import (
	"github.com/google/uuid"
)

type OrgMemberDetail struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	RoleName       string    `json:"role_name"` 
}