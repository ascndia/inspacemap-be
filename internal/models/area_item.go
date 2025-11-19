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
