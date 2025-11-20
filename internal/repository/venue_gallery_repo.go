package repository

import (
	"context"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type venueGalleryRepo struct {
	BaseRepository[entity.VenueGalleryItem, uuid.UUID]
	db *gorm.DB
}

func NewVenueGalleryRepository(db *gorm.DB) VenueGalleryRepository {
	return &venueGalleryRepo{
		BaseRepository: NewBaseRepository[entity.VenueGalleryItem, uuid.UUID](db),
		db: db,
	}
}

func (r *venueGalleryRepo) GetByVenueID(ctx context.Context, venueID uuid.UUID) ([]entity.VenueGalleryItem, error) {
	var items []entity.VenueGalleryItem
	err := r.db.WithContext(ctx).
		Preload("MediaAsset").
		Where("venue_id = ?", venueID).
		Order("sort_order asc").
		Find(&items).Error
	return items, err
}

func (r *venueGalleryRepo) AddVenueItems(ctx context.Context, items []entity.VenueGalleryItem) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&items).Error
}


func (r *venueGalleryRepo) ReorderVenueItems(ctx context.Context, venueID uuid.UUID, mediaIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, mediaID := range mediaIDs {
			if err := tx.Model(&entity.VenueGalleryItem{}).
				Where("venue_id = ? AND media_asset_id = ?", venueID, mediaID).
				Update("sort_order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *venueGalleryRepo) FilterVenueGalleries(ctx context.Context, filter models.VenueGalleryFilter) ([]entity.VenueGalleryItem, error) {
	var items []entity.VenueGalleryItem
	query := r.buildFilterQuery(ctx, filter)
	err := query.Find(&items).Error
	return items, err
}

func (r *venueGalleryRepo) PagedVenueGalleries(ctx context.Context, q models.VenueGalleryQuery) ([]entity.VenueGalleryItem, int64, error) {
	var items []entity.VenueGalleryItem
	var total int64

	db := r.buildFilterQuery(ctx, q.VenueGalleryFilter)

	if err := db.Model(&entity.VenueGalleryItem{}).Count(&total).Error; err != nil {
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

	err := db.Limit(limit).Offset(offset).Find(&items).Error
	return items, total, err
}

func (r *venueGalleryRepo) CursorVenueGalleries(ctx context.Context, q models.VenueGalleryCursor) ([]entity.VenueGalleryItem, string, error) {
	var items []entity.VenueGalleryItem
	db := r.buildFilterQuery(ctx, q.VenueGalleryFilter)

	if q.Cursor != nil && *q.Cursor != "" {
		if cursorID, err := uuid.Parse(*q.Cursor); err == nil {
			var cursorItem entity.VenueGalleryItem
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

// --- QUERY BUILDER ---
func (r *venueGalleryRepo) buildFilterQuery(ctx context.Context, f models.VenueGalleryFilter) *gorm.DB {
	db := r.db.WithContext(ctx).Model(&entity.VenueGalleryItem{})
	if f.VenueID != nil {
		db = db.Where("venue_id = ?", *f.VenueID)
	} 
	if f.MediaAssetID != nil {
		db = db.Where("media_asset_id = ?", *f.MediaAssetID)
	}
	if f.SortOrder != nil {
		db = db.Where("sort_order = ?", *f.SortOrder)
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
	return db
}