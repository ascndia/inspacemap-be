package entity

import (
	"github.com/google/uuid"
)


type Floor struct {
	BaseEntity
	GraphRevisionID uuid.UUID `gorm:"index;not null"`
	VenueID         uuid.UUID `gorm:"index;not null"` 
	Name       string `gorm:"type:varchar(100);not null"` 
	LevelIndex int    `gorm:"not null"` 
	MapImageID *uuid.UUID  `gorm:"type:uuid"`
	MapImage   *MediaAsset `gorm:"foreignKey:MapImageID"` 
	MapWidth       int     `gorm:"default:0"` 
	MapHeight      int     `gorm:"default:0"`
	PixelsPerMeter float64 `gorm:"default:1.0"` 
	GeoReference JSONMap `gorm:"type:jsonb"` 
	IsActive bool `gorm:"default:true"` 
	Nodes []GraphNode `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Areas []Area      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
