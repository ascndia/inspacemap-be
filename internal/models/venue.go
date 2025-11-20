package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateVenueRequest struct {
	Name         string      `json:"name" validate:"required"`
	Slug         string      `json:"slug" validate:"required,alphanum"`
	Description  string      `json:"description"`
	Address      string      `json:"address"`
	City         string      `json:"city"`
	Province     string      `json:"province"`
	Latitude     float64     `json:"latitude" validate:"required"`
	Longitude    float64     `json:"longitude" validate:"required"`
	Visibility   string      `json:"visibility" validate:"oneof=public private unlisted"`
	CoverImageID *uuid.UUID  `json:"cover_image_id"`
	Gallery      []VenueGalleryItemRequest `json:"gallery"` // Langsung set gallery saat create
}

type VenueGalleryItemRequest struct {
	MediaAssetID uuid.UUID `json:"media_asset_id" validate:"required"`
	SortOrder    int       `json:"sort_order"`
	Caption      string    `json:"caption"`
}

type VenueDetail struct {
	ID               uint            `json:"id"`
	Name             string          `json:"name"`
	Slug             string          `json:"slug"`
	Address          string          `json:"address"`
	Description    string    `json:"description"`
	City             string          `json:"city"`
	FullAddress      string          `json:"full_address"` // Gabungan address + city
	IsPublic         bool            `json:"is_public"`
	Coordinates      GeoPoint        `json:"coordinates"`
	CoverImageURL    string          `json:"cover_image_url"`
	Gallery          []VenueGalleryDetail `json:"gallery"`
	PointsOfInterest []AreaPinDetail `json:"pois"` 
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}


type VenueGalleryDetail struct {
	MediaID      uuid.UUID `json:"media_id"`
	URL          string    `json:"url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Caption      string    `json:"caption"`
	SortOrder    int       `json:"sort_order"`
}

type VenueFilter struct {
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	Name       *string `json:"name,omitempty"`
	Slug       *string `json:"slug,omitempty"`
	Description *string `json:"description,omitempty"`
	Address    *string `json:"address,omitempty"`
	City       *string `json:"city,omitempty"`
	Province   *string `json:"province,omitempty"`
	PostalCode *string `json:"postal_code,omitempty"`
	Visibility *string `json:"visibility,omitempty"`
	IsLive     *bool   `json:"is_live,omitempty"`
}

type VenueQuery struct {
	VenueFilter
	Limit  *int `json:"limit,omitempty"`
	Offset *int `json:"offset,omitempty"`
	Sort   *string `json:"sort,omitempty"`
}

type VenueQueryCursor struct {
	VenueFilter
	Limit  *int    `json:"limit,omitempty"`
	Cursor *string `json:"cursor,omitempty"`
}