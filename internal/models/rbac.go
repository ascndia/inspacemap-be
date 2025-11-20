package models

import (
	"github.com/google/uuid"
)
type CreateRoleRequest struct {
	Name          string `json:"name" validate:"required,min=3"`
	Description   string `json:"description"`
	PermissionIDs []uuid.UUID `json:"permission_ids" validate:"required,min=1"`
}

type UpdateUserRoleRequest struct {
	TargetUserID uuid.UUID `json:"target_user_id" validate:"required"`
	NewRoleID    uuid.UUID `json:"new_role_id" validate:"required"`
}

type RoleDetail struct {
	ID          uuid.UUID `json:"id"`
	Name        string   `json:"name"`
	IsSystem    bool     `json:"is_system"`
	Permissions []string `json:"permissions"`
}

type PermissionNode struct {
	Group string           `json:"group"`
	Items []PermissionItem `json:"items"`
}

type PermissionItem struct {
	ID          uuid.UUID `json:"id"`
	Key         string `json:"key"`
	Description string `json:"description"`
}