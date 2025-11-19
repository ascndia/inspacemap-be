package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	AuthLocal  AuthProvider = "local"
	AuthGoogle AuthProvider = "google"
	AuthApple  AuthProvider = "apple"
	AuthOIDC   AuthProvider = "oidc"
	AuthSAML   AuthProvider = "saml"
)

const (
	RoleOwner  UserRole = "owner"
	RoleAdmin  UserRole = "admin"
	RoleEditor UserRole = "editor"
	RoleViewer UserRole = "viewer"
)

type AuthProvider string

type User struct {
	BaseEntity
	FullName  string `gorm:"type:varchar(100)"`
	Email     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash   string       `json:"-"` 
	AuthProvider   AuthProvider `gorm:"type:varchar(20);default:'local'"` 
	ProviderUserID string       `gorm:"type:varchar(100);index"`
	SSOConfigID    *uuid.UUID       `gorm:"type:uuid;index"`
	SSOConfig      *OrganizationSSO `gorm:"foreignKey:SSOConfigID"`
	IsEmailVerified bool `gorm:"default:false"`
	AvatarURL       string
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