package entity

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
    BaseEntity
	Name        string  `gorm:"type:varchar(100);not null"`
	Slug        string  `gorm:"type:varchar(100);uniqueIndex;not null"`
	LogoURL     string
	Website     string
	IsActive    bool    `gorm:"default:true"`
	Settings    JSONMap `gorm:"type:jsonb"`
	Members     []OrganizationMember `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Invitations []UserInvitation     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Venues      []Venue
	MediaAssets []MediaAsset
	ApiKeys     []ApiKey
}

type OrganizationMember struct {
	BaseEntity
	OrganizationID uuid.UUID `gorm:"type:uuid;index;not null;uniqueIndex:idx_org_user"`
	UserID         uuid.UUID `gorm:"type:uuid;index;not null;uniqueIndex:idx_org_user"`
	RoleID         uuid.UUID `gorm:"index;not null"`
	Role           Role      `gorm:"foreignKey:RoleID"`
	JoinedAt       time.Time `gorm:"autoCreateTime"`
	User         User         `gorm:"foreignKey:UserID"`
	Organization Organization `gorm:"foreignKey:OrganizationID"`
}

type ApiKey struct {
	BaseEntity
	OrganizationID uuid.UUID `gorm:"type:uuid;index;not null"`
	Name           string
	IsActive       bool      `gorm:"default:true"`
}