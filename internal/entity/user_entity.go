package entity

import (
	"github.com/google/uuid"
)

type UserRole string

const (
	RoleOwner  UserRole = "owner"
	RoleAdmin  UserRole = "admin"
	RoleEditor UserRole = "editor"
	RoleViewer UserRole = "viewer"
)

type User struct {
	BaseEntity
	Email           string       `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash    string       `gorm:"type:varchar(255)"` // Nullable jika login via Google
	FullName        string       `gorm:"type:varchar(100)"`
	AvatarURL       string       `gorm:"type:text"`
	IsEmailVerified bool         `gorm:"default:false"`
	OrganizationID  uuid.UUID    `gorm:"index;not null"`
	Organization    Organization `gorm:"foreignKey:OrganizationID"`
	RoleID          uuid.UUID    `gorm:"index;not null"`
	Role            Role         `gorm:"foreignKey:RoleID"`
}
