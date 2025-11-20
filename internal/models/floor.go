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

type FloorFilter struct {
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	GraphRevisionID *uuid.UUID `json:"graph_revision_id,omitempty"`
	VenueID        *uuid.UUID `json:"venue_id,omitempty"`
	Name           *string    `json:"name,omitempty"`
	LevelIndex     *int       `json:"level_index,omitempty"`
	MinLevelIndex  *int       `json:"min_level_index,omitempty"`
	MaxLevelIndex  *int       `json:"max_level_index,omitempty"`
	MinMapWidth    *int       `json:"min_map_width,omitempty"`
	MaxMapWidth    *int       `json:"max_map_width,omitempty"`
	MinMapHeight   *int       `json:"min_map_height,omitempty"`
	MaxMapHeight   *int       `json:"max_map_height,omitempty"`
	PixelsPerMeter *float64   `json:"pixels_per_meter,omitempty"`
	IsActive       *bool      `json:"is_active,omitempty"`
}

type FloorQuery struct {
	FloorFilter
	Limit  *int    `json:"limit,omitempty"`
	Offset *int    `json:"offset,omitempty"`
	Sort   *string `json:"sort,omitempty"`
}

type FloorQueryCursor struct {
	FloorFilter
	Limit  *int    `json:"limit,omitempty"`
	Cursor *string `json:"cursor,omitempty"`
}