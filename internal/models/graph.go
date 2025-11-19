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
