package models

import (
	"time"

	"github.com/google/uuid"
)

type CreateVenueRequest struct {
	Name         string                    `json:"name" validate:"required"`
	Slug         string                    `json:"slug" validate:"required,alphanum"`
	Description  string                    `json:"description"`
	Address      string                    `json:"address"`
	City         string                    `json:"city"`
	PostalCode   string                    `json:"postal_code"`
	Province     string                    `json:"province"`
	Latitude     float64                   `json:"latitude" validate:"required"`
	Longitude    float64                   `json:"longitude" validate:"required"`
	Visibility   string                    `json:"visibility" validate:"oneof=public private unlisted"`
	CoverImageID *uuid.UUID                `json:"cover_image_id"`
	Gallery      []VenueGalleryItemRequest `json:"gallery"` // Langsung set gallery saat create
}

type UpdateVenueRequest struct {
	Name         *string    `json:"name,omitempty" validate:"omitempty,min=3"`
	Slug         *string    `json:"slug,omitempty" validate:"omitempty,alphanum,min=3"`
	Description  *string    `json:"description,omitempty"`
	Address      *string    `json:"address,omitempty"`
	City         *string    `json:"city,omitempty"`
	Province     *string    `json:"province,omitempty"`
	PostalCode   *string    `json:"postal_code,omitempty"`
	Latitude     *float64   `json:"latitude,omitempty"`
	Longitude    *float64   `json:"longitude,omitempty"`
	Visibility   *string    `json:"visibility,omitempty" validate:"omitempty,oneof=public private unlisted"`
	CoverImageID *uuid.UUID `json:"cover_image_id,omitempty"`
}

type VenueGalleryItemRequest struct {
	MediaAssetID uuid.UUID `json:"media_asset_id" validate:"required"`
	SortOrder    int       `json:"sort_order"`
	Caption      string    `json:"caption"`
	IsVisible    bool      `json:"is_visible"`
	IsFeatured   bool      `json:"is_featured"`
}

type VenueDetail struct {
	ID               uuid.UUID            `json:"id"`
	OrganizationID   uuid.UUID            `json:"organization_id"`
	Name             string               `json:"name"`
	Slug             string               `json:"slug"`
	Description      string               `json:"description"`
	Address          string               `json:"address"`
	City             string               `json:"city"`
	Province         string               `json:"province"`
	PostalCode       string               `json:"postal_code"`
	FullAddress      string               `json:"full_address"`
	Coordinates      GeoPoint             `json:"coordinates"`
	Visibility       string               `json:"visibility"`
	CoverImageURL    string               `json:"cover_image_url,omitempty"`
	Gallery          []VenueGalleryDetail `json:"gallery"`
	PointsOfInterest []AreaPinDetail      `json:"pois"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
}

type VenueListItem struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	City          string    `json:"city"`
	CoverImageURL string    `json:"cover_image_url,omitempty"`
	Visibility    string    `json:"visibility"`
	IsLive        bool      `json:"is_live"`
}

type VenueFilter struct {
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	Name           *string    `json:"name,omitempty"`
	Slug           *string    `json:"slug,omitempty"`
	Description    *string    `json:"description,omitempty"`
	Address        *string    `json:"address,omitempty"`
	City           *string    `json:"city,omitempty"`
	Province       *string    `json:"province,omitempty"`
	PostalCode     *string    `json:"postal_code,omitempty"`
	Visibility     *string    `json:"visibility,omitempty"`
	IsLive         *bool      `json:"is_live,omitempty"`
}

type VenueQuery struct {
	VenueFilter
	Limit  *int    `json:"limit,omitempty"`
	Offset *int    `json:"offset,omitempty"`
	Sort   *string `json:"sort,omitempty"`
}

type VenueQueryCursor struct {
	VenueFilter
	Limit  *int    `json:"limit,omitempty"`
	Cursor *string `json:"cursor,omitempty"`
}
