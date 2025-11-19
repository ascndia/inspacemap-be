package entity

import (
	"github.com/google/uuid"
)

type PermissionKey string

const (
	PermVenueCreate  PermissionKey = "venue:create"
	PermVenueUpdate  PermissionKey = "venue:update"
	PermVenueDelete  PermissionKey = "venue:delete"
	PermGraphEdit    PermissionKey = "graph:edit"    // Create/Move Node
	PermGraphPublish PermissionKey = "graph:publish" // Publish Draft
	PermOrgSettings  PermissionKey = "org:settings"  // Edit SSO, Logo
	PermOrgBilling   PermissionKey = "org:billing"
	PermUserInvite   PermissionKey = "user:invite"
	PermMediaUpload  PermissionKey = "media:upload"
	PermMediaDelete  PermissionKey = "media:delete"
)

type Permission struct {
	BaseEntity
	Key         PermissionKey `gorm:"type:varchar(50);uniqueIndex;not null"` // e.g. "venue:create"
	Description string        `gorm:"type:varchar(255)"`                     // e.g. "Allows creating new venues"
	Group       string        `gorm:"type:varchar(50)"`                      // e.g. "CMS", "Graph", "Billing"
}

type Role struct {
	BaseEntity
	OrganizationID *uuid.UUID `gorm:"index"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID"`
	Name        string `gorm:"type:varchar(50);not null"`
	Description string `gorm:"type:varchar(255)"`
	IsSystem    bool   `gorm:"default:false"` 
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type RolePermission struct {
	RoleID			uuid.UUID `gorm:"primaryKey"`
	PermissionID 	uuid.UUID `gorm:"primaryKey"`
}