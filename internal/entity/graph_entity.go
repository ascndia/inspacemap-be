package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &j)
}

type RevisionStatus string

const (
	StatusDraft     RevisionStatus = "draft"
	StatusPublished RevisionStatus = "published"
	StatusArchived  RevisionStatus = "archived"
)

type GraphNode struct {
	BaseEntity
	FloorID uuid.UUID `gorm:"index;not null"`
	X float64
	Y float64
	AreaID  *uuid.UUID `gorm:"index" json:"area_id"`
	Area    *Area `gorm:"foreignKey:AreaID"`
	PanoramaAssetID uuid.UUID `gorm:"type:uuid;index;not null"`
	Panorama        *MediaAsset `gorm:"foreignKey:PanoramaAssetID"` 
	RotationOffset float64 `gorm:"default:0.0"`
	Label          string
	Properties     JSONMap `gorm:"type:jsonb"`
	IsActive       bool    `gorm:"default:true"`
	OutgoingEdges  []GraphEdge `gorm:"foreignKey:FromNodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type GraphEdge struct {
	BaseEntity
	FromNodeID uuid.UUID `gorm:"index;not null"`
	FromNode   *GraphNode `gorm:"foreignKey:FromNodeID"`
	ToNodeID   uuid.UUID `gorm:"index;not null"`
	ToNode     *GraphNode `gorm:"foreignKey:ToNodeID"`
	Heading  float64
	Distance float64
	Type     string `gorm:"default:'walk'"`
	IsActive   bool   `gorm:"default:true"`
}

type GraphRevision struct {
	BaseEntity
	OrganizationID uuid.UUID          `gorm:"index;not null"`
	CreatedByID    uuid.UUID          `gorm:"index;not null"`
	CreatedBy      User               `gorm:"foreignKey:CreatedByID"`
	VenueID   uuid.UUID           `gorm:"index;not null"`
	Venue     Venue               `gorm:"foreignKey:VenueID"`
	Status    RevisionStatus `gorm:"type:varchar(20);default:'draft'"`
	Note      string
	StartNodeID *uuid.UUID      `gorm:"index" json:"start_node_id"`
	StartNode   *GraphNode `gorm:"foreignKey:StartNodeID"`
	Floors    []Floor `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
