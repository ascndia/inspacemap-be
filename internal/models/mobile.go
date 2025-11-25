package models

import (
	"time"

	"github.com/google/uuid"
)

type ManifestResponse struct {
	VenueID     uuid.UUID   `json:"venue_id"`
	VenueName   string      `json:"venue_name"`
	LastUpdated time.Time   `json:"last_updated"`
	Floors      []FloorData `json:"floors"`
	StartNodeID uuid.UUID   `json:"start_node_id"`
}

type FloorData struct {
	ID          uuid.UUID  `json:"id"`
	LevelName   string     `json:"name"`
	LevelIndex  int        `json:"level_index"`
	MapImageURL string     `json:"map_image_url"`
	MapWidth    int        `json:"width"`
	MapHeight   int        `json:"height"`
	Nodes       []NodeData `json:"nodes"`
	Areas       []AreaData `json:"areas"`
}

type NodeData struct {
	ID             uuid.UUID      `json:"id"`
	X              int            `json:"x"`
	Y              int            `json:"y"`
	PanoramaURL    string         `json:"panorama_url"` // Full Resolution URL
	RotationOffset float64        `json:"rotation_offset"`
	AreaID         *uuid.UUID     `json:"area_id,omitempty"`
	AreaName       string         `json:"area_name,omitempty"`
	Label          string         `json:"label,omitempty"`
	Neighbors      []NeighborData `json:"neighbors"`
}

type NeighborData struct {
	TargetNodeID uuid.UUID `json:"target_node_id"`
	Heading      float64   `json:"heading"`   // Arah kompas absolut
	Distance     float64   `json:"distance"`  // Jarak dalam pixel/meter
	Type         string    `json:"type"`      // 'walk', 'stairs' (Untuk icon panah beda)
	IsActive     bool      `json:"is_active"` // Jika false, jangan gambar panah
}

type AreaData struct {
	ID            uuid.UUID           `json:"id"`
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	Category      string              `json:"category"`
	Latitude      float64             `json:"latitude"`
	Longitude     float64             `json:"longitude"`
	Boundary      []BoundaryPoint     `json:"boundary"`
	StartNodeID   *uuid.UUID          `json:"start_node_id,omitempty"`
	CoverImageURL string              `json:"cover_image_url,omitempty"`
	Gallery       []AreaGalleryDetail `json:"gallery"`
}
