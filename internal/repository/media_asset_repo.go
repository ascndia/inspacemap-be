package repository

import (
	"context"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mediaRepo struct {
	BaseRepository[entity.MediaAsset, uuid.UUID]
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) MediaAssetRepository {
	return &mediaRepo{
		BaseRepository: NewBaseRepository[entity.MediaAsset, uuid.UUID](db),
		db:             db,
	}
}

func (r *mediaRepo) GetAssetByID(ctx context.Context, id uuid.UUID) (*entity.MediaAsset, error) {
	var asset entity.MediaAsset
	if err := r.db.WithContext(ctx).First(&asset, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *mediaRepo) FilterMediaAssets(ctx context.Context, filter models.MediaAssetFilter) ([]entity.MediaAsset, error) {
	var assets []entity.MediaAsset
	query := r.buildFilterQuery(ctx, filter)
	
	err := query.Find(&assets).Error
	return assets, err
}

func (r *mediaRepo) PagedMediaAssets(ctx context.Context, q models.MediaAssetQuery) ([]entity.MediaAsset, int64, error) {
	var assets []entity.MediaAsset
	var total int64

	db := r.buildFilterQuery(ctx, q.MediaAssetFilter)

	if err := db.Model(&entity.MediaAsset{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if q.Sort != nil {
		db = db.Order(*q.Sort)
	} else {
		db = db.Order("created_at desc")
	}

	limit := 10
	if q.Limit != nil && *q.Limit > 0 {
		limit = *q.Limit
	}
	offset := 0
	if q.Offset != nil && *q.Offset >= 0 {
		offset = *q.Offset
	}

	err := db.Limit(limit).Offset(offset).Find(&assets).Error
	return assets, total, err
}

func (r *mediaRepo) CursorMediaAssets(ctx context.Context, q models.MediaAssetQueryCursor) ([]entity.MediaAsset, string, error) {
	var assets []entity.MediaAsset
	db := r.buildFilterQuery(ctx, q.MediaAssetFilter)

	if q.Cursor != nil && *q.Cursor != "" {
		if cursorID, err := uuid.Parse(*q.Cursor); err == nil {
			var cursorAsset entity.MediaAsset
			if err := r.db.Select("created_at").First(&cursorAsset, "id = ?", cursorID).Error; err == nil {
				db = db.Where("(created_at, id) < (?, ?)", cursorAsset.CreatedAt, cursorID)
			}
		}
	}

	limit := 10
	if q.Limit != nil && *q.Limit > 0 {
		limit = *q.Limit
	}

	err := db.Order("created_at desc, id desc").
		Limit(limit + 1).
		Find(&assets).Error

	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(assets) > limit {
		assets = assets[:limit]
		nextCursor = assets[len(assets)-1].ID.String()
	}

	return assets, nextCursor, nil
}

func (r *mediaRepo) buildFilterQuery(ctx context.Context, f models.MediaAssetFilter) *gorm.DB {
	db := r.db.WithContext(ctx)

	if f.OrganizationID != nil {
		db = db.Where("organization_id = ?", *f.OrganizationID)
	}

	if f.StorageProvider != "" {
		db = db.Where("storage_provider = ?", f.StorageProvider) 
	}
	if f.Bucket != "" {
		db = db.Where("bucket = ?", f.Bucket)
	}
	if f.Region != "" {
		db = db.Where("region = ?", f.Region) 
	}

	if f.MimeType != nil {
		db = db.Where("mime_type = ?", *f.MimeType)
	}
	if f.Type != nil {
		db = db.Where("type = ?", *f.Type)
	}
	if f.Visibility != nil {
		db = db.Where("visibility = ?", *f.Visibility)
	}
	if f.AltText != nil {
		db = db.Where("file_name ILIKE ? OR alt_text ILIKE ?", "%"+*f.AltText+"%", "%"+*f.AltText+"%")
	}
	if f.BlurHash != nil {
		db = db.Where("blur_hash IS NOT NULL") 
	}

	if f.MinSizeInBytes != nil {
		db = db.Where("size_in_bytes >= ?", *f.MinSizeInBytes)
	}
	if f.MaxSizeInBytes != nil {
		db = db.Where("size_in_bytes <= ?", *f.MaxSizeInBytes)
	}

	if f.UploadedBefore != nil {
		if t, err := time.Parse(time.RFC3339, *f.UploadedBefore); err == nil {
			db = db.Where("created_at <= ?", t)
		}
	}
	if f.UploadedAfter != nil {
		if t, err := time.Parse(time.RFC3339, *f.UploadedAfter); err == nil {
			db = db.Where("created_at >= ?", t)
		}
	}

	return db
}