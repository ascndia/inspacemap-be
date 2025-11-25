package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type BoundaryPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Boundary []BoundaryPoint

func (b Boundary) Value() (driver.Value, error) {
	return json.Marshal(b)
}

func (b *Boundary) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, b)
}

type Area struct {
	BaseEntity
	VenueID         uuid.UUID     `gorm:"index;not null"`
	Venue           Venue         `gorm:"foreignKey:VenueID"`
	FloorID         uuid.UUID     `gorm:"index;not null"`
	Floor           Floor         `gorm:"foreignKey:FloorID"`
	GraphRevisionID uuid.UUID     `gorm:"index;not null"`
	GraphRevision   GraphRevision `gorm:"foreignKey:GraphRevisionID"`
	Name            string        `gorm:"type:varchar(100);not null"`
	Slug            string        `gorm:"type:varchar(100);index"`
	Label           string        `gorm:"type:varchar(100)"`
	Description     string        `gorm:"type:text"`
	Latitude        float64       `gorm:"type:decimal(10,8)"`
	Longitude       float64       `gorm:"type:decimal(11,8)"`
	Boundary        Boundary      `gorm:"type:jsonb"`
	StartNodeID     *uuid.UUID    `gorm:"index"`
	StartNode       *GraphNode    `gorm:"foreignKey:StartNodeID"`
	LabelX          float64
	LabelY          float64
	Category        string `gorm:"type:varchar(50);index"`
	CoverImageID    *uuid.UUID
	CoverImage      *MediaAsset       `gorm:"foreignKey:CoverImageID"`
	Gallery         []AreaGalleryItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type AreaGalleryItem struct {
	BaseEntity
	AreaID       uuid.UUID  `gorm:"type:uuid;index:idx_area_media,unique;not null"`
	MediaAssetID uuid.UUID  `gorm:"type:uuid;index:idx_area_media,unique;not null"`
	SortOrder    int        `gorm:"default:0"`
	Caption      string     `gorm:"type:varchar(255)"`
	IsVisible    bool       `gorm:"default:true"`
	Area         Area       `gorm:"foreignKey:AreaID"`
	MediaAsset   MediaAsset `gorm:"foreignKey:MediaAssetID"`
}
