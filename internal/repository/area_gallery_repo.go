package repository

import (
	"context"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type areaGalleryRepo struct {
	BaseRepository[entity.AreaGalleryItem, uuid.UUID]
	db *gorm.DB
}

func NewAreaGalleryRepository(db *gorm.DB) AreaGalleryRepository {
	return &areaGalleryRepo{
		BaseRepository: NewBaseRepository[entity.AreaGalleryItem, uuid.UUID](db),
		db:             db,
	}
}

func (r *areaGalleryRepo) AddAreaItems(ctx context.Context, items []entity.AreaGalleryItem) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&items).Error
}

func (r *areaGalleryRepo) UpdateAreaItem(ctx context.Context, item *entity.AreaGalleryItem) error {
	return r.db.WithContext(ctx).Model(&entity.AreaGalleryItem{}).
		Where("area_id = ? AND media_asset_id = ?", item.AreaID, item.MediaAssetID).
		Updates(map[string]interface{}{
			"caption":    item.Caption,
			"is_visible": item.IsVisible,
			"sort_order": item.SortOrder,
		}).Error
}

func (r *areaGalleryRepo) RemoveAreaItem(ctx context.Context, areaID, mediaID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("area_id = ? AND media_asset_id = ?", areaID, mediaID).
		Delete(&entity.AreaGalleryItem{}).Error
}

func (r *areaGalleryRepo) ReorderAreaItems(ctx context.Context, areaID uuid.UUID, mediaIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, mediaID := range mediaIDs {
			if err := tx.Model(&entity.AreaGalleryItem{}).
				Where("area_id = ? AND media_asset_id = ?", areaID, mediaID).
				Update("sort_order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
func (r *areaGalleryRepo) GetByAreaID(ctx context.Context, areaID uuid.UUID) ([]entity.AreaGalleryItem, error) {
	var items []entity.AreaGalleryItem
	err := r.db.WithContext(ctx).
		Preload("MediaAsset"). // Load data medianya
		Where("area_id = ?", areaID).
		Order("sort_order asc"). // Urutkan sesuai urutan display
		Find(&items).Error
	return items, err
}

func (r *areaGalleryRepo) GetByVenueID(ctx context.Context, venueID uuid.UUID) ([]entity.AreaGalleryItem, error) {
	var items []entity.AreaGalleryItem
	// Join ke table Area untuk filter by VenueID
	err := r.db.WithContext(ctx).
		Joins("JOIN areas ON areas.id = area_gallery_items.area_id").
		Where("areas.venue_id = ?", venueID).
		Preload("MediaAsset").
		Preload("Area"). // Load info areanya juga
		Order("area_gallery_items.sort_order asc").
		Find(&items).Error
	return items, err
}

// 2. Filter Logic

func (r *areaGalleryRepo) FilterAreaGalleries(ctx context.Context, filter models.AreaGalleryFilter) ([]entity.AreaGalleryItem, error) {
	var items []entity.AreaGalleryItem
	query := r.buildFilterQuery(ctx, filter)
	err := query.Find(&items).Error
	return items, err
}

func (r *areaGalleryRepo) PagedAreaGalleries(ctx context.Context, q models.AreaGalleryQuery) ([]entity.AreaGalleryItem, int64, error) {
	var items []entity.AreaGalleryItem
	var total int64

	db := r.buildFilterQuery(ctx, q.AreaGalleryFilter)

	if err := db.Model(&entity.AreaGalleryItem{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if q.Sort != nil {
		db = db.Order(*q.Sort)
	} else {
		// Default sort by SortOrder, lalu CreatedAt
		db = db.Order("sort_order asc, created_at desc")
	}

	limit := 10
	if q.Limit != nil && *q.Limit > 0 {
		limit = *q.Limit
	}
	offset := 0
	if q.Offset != nil && *q.Offset >= 0 {
		offset = *q.Offset
	}

	err := db.Limit(limit).Offset(offset).Find(&items).Error
	return items, total, err
}

func (r *areaGalleryRepo) CursorAreaGalleries(ctx context.Context, q models.AreaGalleryCursor) ([]entity.AreaGalleryItem, string, error) {
	var items []entity.AreaGalleryItem
	db := r.buildFilterQuery(ctx, q.AreaGalleryFilter)

	if q.Cursor != nil && *q.Cursor != "" {
		if cursorID, err := uuid.Parse(*q.Cursor); err == nil {
			var cursorItem entity.AreaGalleryItem
			if err := r.db.Select("created_at").First(&cursorItem, "id = ?", cursorID).Error; err == nil {
				db = db.Where("(created_at, id) < (?, ?)", cursorItem.CreatedAt, cursorID)
			}
		}
	}

	limit := 10
	if q.Limit != nil && *q.Limit > 0 {
		limit = *q.Limit
	}

	err := db.Order("created_at desc, id desc").Limit(limit + 1).Find(&items).Error
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(items) > limit {
		items = items[:limit]
		nextCursor = items[len(items)-1].ID.String()
	}

	return items, nextCursor, nil
}

// --- HELPER: Query Builder ---
func (r *areaGalleryRepo) buildFilterQuery(ctx context.Context, f models.AreaGalleryFilter) *gorm.DB {
	db := r.db.WithContext(ctx).Preload("MediaAsset")

	if f.AreaID != nil {
		db = db.Where("area_id = ?", *f.AreaID)
	}

	// Filter by Venue (Requires Join)
	if f.VenueID != nil {
		db = db.Joins("JOIN areas ON areas.id = area_gallery_items.area_id").
			Where("areas.venue_id = ?", *f.VenueID)
	}

	if f.MediaAssetID != nil {
		db = db.Where("media_asset_id = ?", *f.MediaAssetID)
	}

	if f.Caption != nil {
		db = db.Where("caption ILIKE ?", "%"+*f.Caption+"%")
	}

	if f.IsVisible != nil {
		db = db.Where("is_visible = ?", *f.IsVisible)
	}

	if f.IsFeatured != nil {
		db = db.Where("is_featured = ?", *f.IsFeatured)
	}

	if f.SortOrder != nil {
		db = db.Where("sort_order = ?", *f.SortOrder)
	}

	return db
}
