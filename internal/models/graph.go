package models

import (
	"time"

	"github.com/google/uuid"
)

type PublishDraftRequest struct {
	Note string `json:"note" validate:"max=255"`
}

type RevisionHistoryItem struct {
	ID        uuid.UUID `json:"id"` 
	Status    string    `json:"status"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
}

type GraphRevisionDetail struct {
	ID        uuid.UUID `json:"id"`
	VenueID   uuid.UUID `json:"venue_id"`
	Status    string    `json:"status"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Floors    []FloorDetail `json:"floors"`
}

type FloorDetail struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	LevelIndex     int       `json:"level_index"`
	MapImageURL    string    `json:"map_image_url"`
	MapWidth       int       `json:"map_width"`
	MapHeight      int       `json:"map_height"`
	PixelsPerMeter float64   `json:"pixels_per_meter"`
	IsActive       bool      `json:"is_active"`
	NodesCount     int       `json:"nodes_count"`
	AreasCount     int       `json:"areas_count"`
}
type CreateNodeRequest struct {
	FloorID         uuid.UUID `json:"floor_id" validate:"required"`
	X               float64   `json:"x" validate:"required"` 
	Y               float64   `json:"y" validate:"required"`
	PanoramaAssetID uuid.UUID `json:"panorama_asset_id" validate:"required"`
	Label           string    `json:"label"`
}

type UpdateNodePositionRequest struct {
	ID 		 uuid.UUID    `json:"id" validate:"required"`
	X float64 `json:"x" validate:"required"`
	Y float64 `json:"y" validate:"required"`
}

type UpdateNodeCalibrationRequest struct {
	ID              uuid.UUID `json:"id" validate:"required"`
	RotationOffset float64 `json:"rotation_offset" validate:"required"` 
}

type UpdateNodeRequest struct {
	X                *float64   `json:"x,omitempty"`
	Y                *float64   `json:"y,omitempty"`
	PanoramaAssetID  *uuid.UUID `json:"panorama_asset_id,omitempty"`
	Label            *string    `json:"label,omitempty"`
	RotationOffset   *float64   `json:"rotation_offset,omitempty"`
}

type NodeAdminItem struct {
	ID             uuid.UUID      `json:"id"`
	Label          string    `json:"label"`
	X              float64   `json:"x"`
	Y              float64   `json:"y"`
	RotationOffset float64   `json:"rotation_offset"`
	PanoramaURL    string    `json:"panorama_thumbnail"` 
}

type ConnectNodesRequest struct {
	FromNodeID uuid.UUID `json:"from_node_id" validate:"required"`
	ToNodeID   uuid.UUID `json:"to_node_id" validate:"required"`
}

type CreateNodeResponse struct {
	ID        uuid.UUID `json:"id"`
	FloorID   uuid.UUID `json:"floor_id"`
	X         float64   `json:"x"`
	Y         float64   `json:"y"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateEdgeResponse struct {
	ID        uuid.UUID `json:"id"`
	FromNodeID uuid.UUID `json:"from_node_id"`
	ToNodeID   uuid.UUID `json:"to_node_id"`
	Heading   float64   `json:"heading"`
	Distance  float64   `json:"distance"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateRevisionRequest struct {
	Note string `json:"note" validate:"max=255"`
}
