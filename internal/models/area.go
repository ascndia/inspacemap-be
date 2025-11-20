package models

import "github.com/google/uuid"

type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CreateAreaRequest struct {
	Name         string      `json:"name" validate:"required"`
	FloorID      *uint       `json:"floor_id"` 
	Description  string      `json:"description"`
	Category     string      `json:"category"` 
	Latitude     float64     `json:"latitude"` 
	Longitude    float64     `json:"longitude"`
	MapX         float64     `json:"map_x"`    
	MapY         float64     `json:"map_y"`
	CoverImageID *uuid.UUID  `json:"cover_image_id"`
	Gallery      []AreaItemRequest `json:"gallery"`
}

type AreaPinDetail struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	Coordinates  GeoPoint  `json:"coordinates"`
	ThumbnailURL string    `json:"thumbnail_url"`
	FloorName    string    `json:"floor_name,omitempty"` 
}

type AreaDetail struct {
	ID            uint            `json:"id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Gallery       []AreaGalleryDetail `json:"gallery"`	
	NearestNodeID *uint           `json:"nearest_node_id"` 
}


type AreaFilter struct {
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	VenueID      *uuid.UUID    `json:"venue_id,omitempty"`
	FloorID     *uuid.UUID   `json:"floor_id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Slug        *string `json:"slug,omitempty"`
	Label 	 *string `json:"label,omitempty"`
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty"`
}

type AreaQuery struct {
	AreaFilter
	Limit          *int       `json:"limit,omitempty"`
	Offset         *int       `json:"offset,omitempty"`
	Sort 		*string    `json:"sort,omitempty"`
}

type AreaQueryCursor struct {
	AreaFilter
	Limit  *int    `json:"limit,omitempty"`
	Cursor *string `json:"cursor,omitempty"`
}