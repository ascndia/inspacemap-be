package models

import (
	"time"

	"github.com/google/uuid"
)

type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type BoundaryPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type CreateAreaRequest struct {
	Name         string            `json:"name" validate:"required"`
	FloorID      *uuid.UUID        `json:"floor_id"`
	Description  string            `json:"description"`
	Category     string            `json:"category"`
	Latitude     float64           `json:"latitude"`
	Longitude    float64           `json:"longitude"`
	Boundary     []BoundaryPoint   `json:"boundary"`
	CoverImageID *uuid.UUID        `json:"cover_image_id"`
	Gallery      []AreaItemRequest `json:"gallery"`
}

type UpdateAreaRequest struct {
	Name         *string         `json:"name,omitempty"`
	Description  *string         `json:"description,omitempty"`
	Category     *string         `json:"category,omitempty"`
	Latitude     *float64        `json:"latitude,omitempty"`
	Longitude    *float64        `json:"longitude,omitempty"`
	Boundary     []BoundaryPoint `json:"boundary,omitempty"`
	FloorID      *uuid.UUID      `json:"floor_id,omitempty"`
	CoverImageID *uuid.UUID      `json:"cover_image_id,omitempty"`
}

type AreaPinDetail struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	Coordinates  GeoPoint  `json:"coordinates"`
	ThumbnailURL string    `json:"thumbnail_url"`
	FloorName    string    `json:"floor_name,omitempty"`
}

type AreaListItem struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	FloorID      uuid.UUID `json:"floor_id"`
	FloorName    string    `json:"floor_name"`
	RevisionID   uuid.UUID `json:"revision_id"`
	CoverURL     string    `json:"cover_url"`
	GalleryCount int       `json:"gallery_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AreaDetail struct {
	ID          uuid.UUID           `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Gallery     []AreaGalleryDetail `json:"gallery"`
	StartNodeID *uuid.UUID          `json:"start_node_id"`
}

type AreaEditorDetail struct {
	ID           uuid.UUID           `json:"id"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Category     string              `json:"category"`
	Latitude     float64             `json:"latitude"`
	Longitude    float64             `json:"longitude"`
	Boundary     []BoundaryPoint     `json:"boundary"`
	StartNodeID  *uuid.UUID          `json:"start_node_id"`
	FloorID      uuid.UUID           `json:"floor_id"`
	FloorName    string              `json:"floor_name"`
	RevisionID   uuid.UUID           `json:"revision_id"`
	CoverImageID *uuid.UUID          `json:"cover_image_id"`
	CoverURL     string              `json:"cover_url"`
	Gallery      []AreaGalleryDetail `json:"gallery"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type AreaFilter struct {
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	VenueID        *uuid.UUID `json:"venue_id,omitempty"`
	FloorID        *uuid.UUID `json:"floor_id,omitempty"`
	RevisionID     *uuid.UUID `json:"revision_id,omitempty"`
	Name           *string    `json:"name,omitempty"`
	Slug           *string    `json:"slug,omitempty"`
	Label          *string    `json:"label,omitempty"`
	Description    *string    `json:"description,omitempty"`
	Category       *string    `json:"category,omitempty"`
	Status         *string    `json:"status,omitempty"` // "published", "draft", "all"
}

type AreaQuery struct {
	AreaFilter
	Limit  *int    `json:"limit,omitempty"`
	Offset *int    `json:"offset,omitempty"`
	Sort   *string `json:"sort,omitempty"`
	Status *string `json:"status,omitempty"` // "published", "draft", "all"
}

type AreaQueryCursor struct {
	AreaFilter
	Limit  *int    `json:"limit,omitempty"`
	Cursor *string `json:"cursor,omitempty"`
}

type SetStartNodeRequest struct {
	NodeID uuid.UUID `json:"node_id" validate:"required"`
}

type SetStartNodeResponse struct {
	Warning string `json:"warning,omitempty"`
}
