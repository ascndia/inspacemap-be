package entity

import (
	"github.com/google/uuid"
)

type Area struct {
	BaseEntity
	VenueID uuid.UUID `gorm:"index;not null"`
	FloorID uuid.UUID `gorm:"index;not null"` 
	Name        string `gorm:"type:varchar(100);not null"`
	Slug        string `gorm:"type:varchar(100);index"`
	Label	    string `gorm:"type:varchar(100)"`
	Description string `gorm:"type:text"`
	Latitude    float64 `gorm:"type:decimal(10,8)"` 
	Longitude   float64 `gorm:"type:decimal(11,8)"`
	MapX        float64 
	MapY        float64
	Category    string `gorm:"type:varchar(50);index"`
	CoverImageID *uuid.UUID
	CoverImage   *MediaAsset `gorm:"foreignKey:CoverImageID"`
	Gallery      []AreaGalleryItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type AreaGalleryItem struct {
	BaseEntity
	AreaID       uuid.UUID  `gorm:"primaryKey;autoIncrement:false"`
	MediaAssetID uuid.UUID  `gorm:"type:uuid;primaryKey;autoIncrement:false"`
	SortOrder    int        `gorm:"default:0"`
	Caption      string     `gorm:"type:varchar(255)"`
	IsVisible    bool       `gorm:"default:true"`
	Area         Area       `gorm:"foreignKey:AreaID"`
	MediaAsset   MediaAsset `gorm:"foreignKey:MediaAssetID"`
}


