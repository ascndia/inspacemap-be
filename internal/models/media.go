package models

import "github.com/google/uuid"

type MediaAssetRequest struct {
	OrganizationID uuid.UUID `json:"organization_id" validate:"required"`
	FileType       string    `json:"file_type" validate:"required"`
}

type MediaAssetResponse struct {
	ID           uuid.UUID `json:"id"`
	PublicURL    string    `json:"url"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	Type         string    `json:"type"` // 'panorama', 'image'
	Width        int       `json:"width,omitempty"`
	Height       int       `json:"height,omitempty"`
}

type PresignedUploadRequest struct {
	FileName string `json:"file_name" validate:"required"`
	FileType string `json:"file_type" validate:"required"` 
	Category string `json:"category" validate:"oneof=panorama icon floorplan"`
}

type PresignedUploadResponse struct {
	UploadURL string    `json:"upload_url"` 
	AssetID   uuid.UUID `json:"asset_id"`   
}

type SizeFilter struct {
	SizeInBytes  *int64     `json:"size_in_bytes,omitempty"`
	MinSizeInBytes *int64     `json:"min_size_in_bytes,omitempty"`
	MaxSizeInBytes *int64     `json:"max_size_in_bytes,omitempty"`
	Width	  *int       `json:"width,omitempty"`
	Height	  *int       `json:"height,omitempty"`
	MinWidth  *int `json:"min_width,omitempty"`
	MaxWidth  *int `json:"max_width,omitempty"`
	MinHeight *int `json:"min_height,omitempty"`
	MaxHeight *int `json:"max_height,omitempty"`
}

type StorageProviderFilter struct {
	StorageProvider string `json:"storage_provider,omitempty"`
	Bucket		  string `json:"bucket,omitempty"`
	Key		  string `json:"key,omitempty"`
	Region		  string `json:"region,omitempty"`
}

type MediaAssetFilter struct {
	SizeFilter
	StorageProviderFilter
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	FileType      *string    `json:"file_type,omitempty"`
	MimeType      *string    `json:"mime_type,omitempty"`
	Type          *string    `json:"type,omitempty"`
	Height        *int       `json:"height,omitempty"`
	UploadedAt    *string    `json:"uploaded_at,omitempty"`
	Visibility    *string    `json:"visibility,omitempty"`
	UploadedBy    *uuid.UUID `json:"uploaded_by,omitempty"`
	IsEncrypted  *bool      `json:"is_encrypted,omitempty"`
	BlurHash	 *string    `json:"blur_hash,omitempty"`
	AltText     *string    `json:"alt_text,omitempty"`
	UploadedBefore *string   `json:"uploaded_before,omitempty"`
	UploadedAfter  *string   `json:"uploaded_after,omitempty"`
}

type MediaAssetQuery struct {
	MediaAssetFilter
	Limit  *int    `json:"limit,omitempty"`
	Offset *int    `json:"offset,omitempty"`
}

type MediaAssetQueryCursor struct {
	MediaAssetFilter
	Limit  *int    `json:"limit,omitempty"`
	Cursor *string `json:"cursor,omitempty"`
}