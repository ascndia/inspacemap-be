package models

import "github.com/google/uuid"

type CreateAreaItemRequest struct {
	AreaID       uuid.UUID `json:"area_id" validate:"required"`
	MediaAssetID uuid.UUID `json:"media_asset_id" validate:"required"`
	SortOrder    int       `json:"sort_order"`
	Caption      string    `json:"caption"`
	IsVisible    bool      `json:"is_visible"`
}

type AreaGalleryDetail struct {
	MediaID      uuid.UUID `json:"media_id"`
	URL          string    `json:"url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Caption      string    `json:"caption"`
	SortOrder    int       `json:"sort_order"`
}

type AreaItemRequest struct {
	MediaAssetID uuid.UUID `json:"media_asset_id" validate:"required"`
	SortOrder    int       `json:"sort_order"`
	Caption      string    `json:"caption"`
}

type AreaGalleryFilter struct {
	AreaID       *uuid.UUID `json:"area_id,omitempty"`
	VenueID     *uuid.UUID `json:"venue_id,omitempty"`
	MediaAssetID *uuid.UUID `json:"media_asset_id,omitempty"`
	SortOrder    *int       `json:"sort_order,omitempty"`
	Caption      *string    `json:"caption,omitempty"`
	IsVisible    *bool      `json:"is_visible,omitempty"`
	IsFeatured   *bool      `json:"is_featured,omitempty"`
}

type AreaGalleryQuery struct {
	AreaGalleryFilter
	Limit  *int    `json:"limit,omitempty"`
	Offset *int    `json:"offset,omitempty"`
	Sort   *string `json:"sort,omitempty"`
}

type AreaGalleryCursor struct {
	AreaGalleryFilter
	Limit  *int    `json:"limit,omitempty"`
	Cursor *string `json:"cursor,omitempty"`
}