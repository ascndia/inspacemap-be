package repository

import (
	"context"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type floorRepo struct {
	BaseRepository[entity.Floor, uuid.UUID]
	db *gorm.DB
}

func NewFloorRepository(db *gorm.DB) FloorRepository {
	return &floorRepo{
		BaseRepository: NewBaseRepository[entity.Floor, uuid.UUID](db),
		db:             db,
	}
}

func (r *floorRepo) GetByVenueID(ctx context.Context, venueID uuid.UUID) ([]entity.Floor, error) {
	var floors []entity.Floor
	err := r.db.WithContext(ctx).
		Where("venue_id = ?", venueID).
		Preload("MapImage").
		Order("level_index asc"). 
		Find(&floors).Error
	return floors, err
}

func (r *floorRepo) GetByGraphRevisionID(ctx context.Context, revisionID uuid.UUID) ([]entity.Floor, error) {
	var floors []entity.Floor
	err := r.db.WithContext(ctx).
		Where("graph_revision_id = ?", revisionID).
		Preload("MapImage").
		Order("level_index asc").
		Find(&floors).Error
	return floors, err
}


func (r *floorRepo) UpdateFloorMap(ctx context.Context, id uuid.UUID, mapImageID *uuid.UUID, pixelsPerMeter float64) error {
	updates := map[string]interface{}{
		"pixels_per_meter": pixelsPerMeter,
	}
	if mapImageID != nil {
		updates["map_image_id"] = *mapImageID
	} else {
		updates["map_image_id"] = nil
	}

	return r.db.WithContext(ctx).Model(&entity.Floor{}).
		Where("id = ?", id).
		Updates(updates).Error
}


func (r *floorRepo) FilterFloors(ctx context.Context, filter models.FloorFilter) ([]entity.Floor, error) {
	var floors []entity.Floor
	query := r.buildFilterQuery(ctx, filter)
	err := query.Find(&floors).Error
	return floors, err
}

func (r *floorRepo) PagedFloors(ctx context.Context, q models.FloorQuery) ([]entity.Floor, int64, error) {
	var floors []entity.Floor
	var total int64

	db := r.buildFilterQuery(ctx, q.FloorFilter)

	if err := db.Model(&entity.Floor{}).Count(&total).Error; err != nil {
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

	err := db.Limit(limit).Offset(offset).Find(&floors).Error
	return floors, total, err
}

func (r *floorRepo) CursorFloors(ctx context.Context, q models.FloorQueryCursor) ([]entity.Floor, string, error) {
	var floors []entity.Floor
	db := r.buildFilterQuery(ctx, q.FloorFilter)

	if q.Cursor != nil && *q.Cursor != "" {
		if cursorID, err := uuid.Parse(*q.Cursor); err == nil {
			var cursorData entity.Floor
			if err := r.db.Select("created_at").First(&cursorData, "id = ?", cursorID).Error; err == nil {
				db = db.Where("(created_at, id) < (?, ?)", cursorData.CreatedAt, cursorID)
			}
		}
	}

	limit := 10
	if q.Limit != nil && *q.Limit > 0 {
		limit = *q.Limit
	}

	err := db.Order("created_at desc, id desc").Limit(limit + 1).Find(&floors).Error
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(floors) > limit {
		floors = floors[:limit]
		nextCursor = floors[len(floors)-1].ID.String()
	}

	return floors, nextCursor, nil
}

func (r *floorRepo) buildFilterQuery(ctx context.Context, f models.FloorFilter) *gorm.DB {
	db := r.db.WithContext(ctx).Preload("MapImage")

	if f.OrganizationID != nil {
		db = db.Joins("JOIN venues ON venues.id = floors.venue_id").
			Where("venues.organization_id = ?", *f.OrganizationID)
	}

	if f.VenueID != nil {
		db = db.Where("venue_id = ?", *f.VenueID)
	}
	if f.GraphRevisionID != nil {
		db = db.Where("graph_revision_id = ?", *f.GraphRevisionID)
	}

	if f.Name != nil {
		db = db.Where("name ILIKE ?", "%"+*f.Name+"%")
	}

	if f.LevelIndex != nil {
		db = db.Where("level_index = ?", *f.LevelIndex)
	}
	if f.MinLevelIndex != nil {
		db = db.Where("level_index >= ?", *f.MinLevelIndex)
	}
	if f.MaxLevelIndex != nil {
		db = db.Where("level_index <= ?", *f.MaxLevelIndex)
	}

	if f.MinMapWidth != nil {
		db = db.Where("map_width >= ?", *f.MinMapWidth)
	}
	if f.MaxMapWidth != nil {
		db = db.Where("map_width <= ?", *f.MaxMapWidth)
	}
	if f.MinMapHeight != nil {
		db = db.Where("map_height >= ?", *f.MinMapHeight)
	}
	if f.MaxMapHeight != nil {
		db = db.Where("map_height <= ?", *f.MaxMapHeight)
	}

	if f.PixelsPerMeter != nil {
		db = db.Where("pixels_per_meter = ?", *f.PixelsPerMeter)
	}
	
	if f.IsActive != nil {
		db = db.Where("is_active = ?", *f.IsActive)
	}

	return db
}