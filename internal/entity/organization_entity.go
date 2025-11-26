package entity

import (
	"github.com/google/uuid"
)

type Organization struct {
	BaseEntity
	Name        string      `gorm:"type:varchar(100);not null"`
	Slug        string      `gorm:"type:varchar(100);uniqueIndex;not null"`
	LogoID      *uuid.UUID  `gorm:"index"`
	Logo        *MediaAsset `gorm:"foreignKey:LogoID"`
	Website     string
	IsActive    bool    `gorm:"default:true"`
	Settings    JSONMap `gorm:"type:jsonb"`
	Users       []User  `gorm:"foreignKey:OrganizationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Venues      []Venue
	MediaAssets []MediaAsset
	ApiKeys     []ApiKey
}

type ApiKey struct {
	BaseEntity
	OrganizationID uuid.UUID `gorm:"type:uuid;index;not null"`
	Name           string
	IsActive       bool `gorm:"default:true"`
}
