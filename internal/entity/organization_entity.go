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
	SSOConfigs  []OrganizationSSO    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Venues      []Venue
	MediaAssets []MediaAsset
	ApiKeys     []ApiKey
}

type OrganizationMember struct {
	BaseEntity
	OrganizationID uuid.UUID `gorm:"type:uuid;index;not null;uniqueIndex:idx_org_user"`
	UserID         uuid.UUID `gorm:"type:uuid;index;not null;uniqueIndex:idx_org_user"`
	RoleID         uint      `gorm:"index;not null"`
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

type OrganizationSSO struct {
	BaseEntity
	OrganizationID uuid.UUID      `gorm:"type:uuid;index;not null"`
	Provider       AuthProvider   `gorm:"type:varchar(20);not null"` 	
	DisplayName    string         `gorm:"type:varchar(100)"` 
	ClientID       string         `gorm:"type:varchar(255)"`
	ClientSecret   string         `gorm:"type:varchar(255)"`
	IssuerURL      string         `gorm:"type:varchar(255)"` // URL server login kampus/kantor
	AuthURL        string         `gorm:"type:varchar(255)"`
	TokenURL       string         `gorm:"type:varchar(255)"`
	UserInfoURL    string         `gorm:"type:varchar(255)"`
	EmailDomains   JSONMap        `gorm:"type:jsonb"` 
	IsActive       bool           `gorm:"default:true"`
}