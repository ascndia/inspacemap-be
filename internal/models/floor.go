package models

import "github.com/google/uuid"

type CreateFloorRequest struct {
	Name           string     `json:"name" validate:"required"`
	LevelIndex     int        `json:"level_index" validate:"required"` // 1, 2, -1
	MapImageID     *uuid.UUID `json:"map_image_id" validate:"required"`
	PixelsPerMeter float64    `json:"pixels_per_meter" validate:"gt=0"`
}

type FloorAdminDetail struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	LevelIndex  int             `json:"level_index"`
	MapImageURL string          `json:"map_image_url"`
	NodesCount  int             `json:"nodes_count"`
	Nodes       []NodeAdminItem `json:"nodes,omitempty"` // Optional list
}