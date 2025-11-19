package models

import (
	"time"
)

type ManifestResponse struct {
	VenueID      uint          `json:"venue_id"`
	VenueName    string        `json:"venue_name"`
	LastUpdated  time.Time     `json:"last_updated"`
	Floors       []FloorData   `json:"floors"`
	StartNodeID  uint        `json:"start_node_id"` 
}

type FloorData struct {
	ID         uint        `json:"id"`
	LevelName  string      `json:"name"`
	LevelIndex int         `json:"level_index"`
	Nodes      []NodeData  `json:"nodes"`
}

type NodeData struct {
	ID             uint           `json:"id"`
	// Koordinat X/Y dikirim agar fitur "Hotspot Scaling by Distance" di Flutter jalan
	X              int            `json:"x"` 
	Y              int            `json:"y"`
	PanoramaURL    string         `json:"panorama"` // Full Resolution URL
	RotationOffset float64        `json:"rotation_offset"`
	Label          string         `json:"label,omitempty"`
	Neighbors      []NeighborData `json:"neighbors"`
}

type NeighborData struct {
	TargetNodeID uint    `json:"target"`
	Heading      float64 `json:"heading"`  // Arah kompas absolut
	Distance     float64 `json:"distance"` // Jarak dalam pixel/meter
	Type         string  `json:"type"`     // 'walk', 'stairs' (Untuk icon panah beda)
	IsActive     bool    `json:"is_active"` // Jika false, jangan gambar panah
}