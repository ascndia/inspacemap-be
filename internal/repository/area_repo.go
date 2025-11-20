package repository

import (
	"context"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type areaRepo struct {
	BaseRepository[entity.Area, uuid.UUID]
	db *gorm.DB
}

func NewAreaRepository(db *gorm.DB) AreaRepository {
	return &areaRepo{
		BaseRepository: NewBaseRepository[entity.Area, uuid.UUID](db),
		db:             db,
	}
}

func (r *areaRepo) GetByVenueID(ctx context.Context, venueID uuid.UUID) ([]entity.Area, error) {
	var areas []entity.Area
	err := r.db.WithContext(ctx).
		Preload("CoverImage").
		Where("venue_id = ?", venueID).
		Find(&areas).Error
	return areas, err
}

func (r *areaRepo) GetByFloorID(ctx context.Context, floorID uuid.UUID) ([]entity.Area, error) {
	var areas []entity.Area
	err := r.db.WithContext(ctx).
		Preload("CoverImage").
		Where("floor_id = ?", floorID).
		Find(&areas).Error
	return areas, err
}

func (r *areaRepo) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.Area, error) {
	var areas []entity.Area
	err := r.db.WithContext(ctx).
		Joins("JOIN venues ON venues.id = areas.venue_id").
		Where("venues.organization_id = ?", orgID).
		Preload("CoverImage").
		Find(&areas).Error
	return areas, err
}

func (r *areaRepo) FilterAreas(ctx context.Context, filter models.AreaFilter) ([]entity.Area, error) {
	var areas []entity.Area
	query := r.buildFilterQuery(ctx, filter)
	err := query.Find(&areas).Error
	return areas, err
}

func (r *areaRepo) PagedAreas(ctx context.Context, q models.AreaQuery) ([]entity.Area, int64, error) {
	var areas []entity.Area
	var total int64

	db := r.buildFilterQuery(ctx, q.AreaFilter)

	if err := db.Model(&entity.Area{}).Count(&total).Error; err != nil {
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

	err := db.Limit(limit).Offset(offset).Find(&areas).Error
	return areas, total, err
}

func (r *areaRepo) CursorAreas(ctx context.Context, q models.AreaQueryCursor) ([]entity.Area, string, error) {
	var areas []entity.Area
	db := r.buildFilterQuery(ctx, q.AreaFilter)

	if q.Cursor != nil && *q.Cursor != "" {
		if cursorID, err := uuid.Parse(*q.Cursor); err == nil {
			var cursorData entity.Area
			if err := r.db.Select("created_at").First(&cursorData, "id = ?", cursorID).Error; err == nil {
				db = db.Where("(created_at, id) < (?, ?)", cursorData.CreatedAt, cursorID)
			}
		}
	}

	limit := 10
	if q.Limit != nil && *q.Limit > 0 {
		limit = *q.Limit
	}

	err := db.Order("created_at desc, id desc").Limit(limit + 1).Find(&areas).Error
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(areas) > limit {
		areas = areas[:limit]
		nextCursor = areas[len(areas)-1].ID.String()
	}

	return areas, nextCursor, nil
}

func (r *areaRepo) buildFilterQuery(ctx context.Context, f models.AreaFilter) *gorm.DB {
	db := r.db.WithContext(ctx).Preload("CoverImage")

	if f.OrganizationID != nil {
		db = db.Joins("JOIN venues ON venues.id = areas.venue_id").
			Where("venues.organization_id = ?", *f.OrganizationID)
	}

	if f.VenueID != nil {
		db = db.Where("venue_id = ?", *f.VenueID)
	}
	
	if f.FloorID != nil {
		db = db.Where("floor_id = ?", *f.FloorID)
	}

	if f.Name != nil {
		db = db.Where("name ILIKE ?", "%"+*f.Name+"%")
	}
	
	if f.Slug != nil {
		db = db.Where("slug = ?", *f.Slug)
	}
	
	if f.Label != nil {
		db = db.Where("label ILIKE ?", "%"+*f.Label+"%")
	}
	
	if f.Category != nil {
		db = db.Where("category = ?", *f.Category)
	}
	
	if f.Description != nil {
		db = db.Where("description ILIKE ?", "%"+*f.Description+"%")
	}

	return db
}