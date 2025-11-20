package models

type CreateRoleRequest struct {
	Name          string `json:"name" validate:"required,min=3"`
	Description   string `json:"description"`
	PermissionIDs []uint `json:"permission_ids" validate:"required,min=1"`
}

type RoleDetail struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	IsSystem    bool     `json:"is_system"`
	Permissions []string `json:"permissions"`
}

type PermissionNode struct {
	Group string           `json:"group"`
	Items []PermissionItem `json:"items"`
}

type PermissionItem struct {
	ID          uint   `json:"id"`
	Key         string `json:"key"`
	Description string `json:"description"`
}