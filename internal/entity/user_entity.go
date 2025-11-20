package entity

import (
	"time"

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
	Email        string `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255)"` // Nullable jika login via Google
	FullName     string `gorm:"type:varchar(100)"`
	AvatarURL    string `gorm:"type:text"`
	IsEmailVerified bool   `gorm:"default:false"`
	Memberships []OrganizationMember `gorm:"foreignKey:UserID"`
}

type UserInvitation struct {
	BaseEntity
	OrganizationID uuid.UUID `gorm:"type:uuid;index;not null"`
	Organization    Organization `gorm:"foreignKey:OrganizationID"`
	Email          string    `gorm:"type:varchar(100);not null"`
	RoleID         uuid.UUID `gorm:"type:uuid;not null"`
	Role           Role      `gorm:"foreignKey:RoleID"`	
	Token          string    `gorm:"type:varchar(255);not null;index"`
	ExpiresAt      time.Time
	InvitedByUserID uuid.UUID `gorm:"type:uuid"`
	Status          string    `gorm:"default:'pending'"`
	AcceptedAt      *time.Time
	InvitedByUser   *User `gorm:"foreignKey:InvitedByUserID"`
}