package repository

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type venueRepo struct {
	BaseRepository[entity.Venue, uuid.UUID]
	db *gorm.DB
}

func NewVenueRepository(db *gorm.DB) VenueRepository {
	return &venueRepo{
		BaseRepository: NewBaseRepository[entity.Venue, uuid.UUID](db),
		db:             db,
	}
}

func (r *venueRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Venue, error) {
	var venue entity.Venue
	err := r.db.WithContext(ctx).
		Preload("CoverImage").
		Preload("Gallery.MediaAsset").
		Preload("PointsOfInterest"). // Load Area/POI
		Where("id = ?", id).
		First(&venue).Error
	return &venue, err
}

func (r *venueRepo) GetBySlug(ctx context.Context, slug string) (*entity.Venue, error) {
	var venue entity.Venue
	err := r.db.WithContext(ctx).
		Preload("CoverImage").
		Preload("Gallery.MediaAsset").
		Preload("PointsOfInterest").
		Where("slug = ?", slug).
		First(&venue).Error
	return &venue, err
}

func (r *venueRepo) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.Venue, error) {
	var venues []entity.Venue
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Preload("CoverImage").
		Order("created_at desc").
		Find(&venues).Error
	return venues, err
}

func (r *venueRepo) FilterVenues(ctx context.Context, filter models.VenueFilter) ([]entity.Venue, error) {
	var venues []entity.Venue
	query := r.buildFilterQuery(ctx, filter)
	err := query.Find(&venues).Error
	return venues, err
}

func (r *venueRepo) PagedVenues(ctx context.Context, q models.VenueQuery) ([]entity.Venue, int64, error) {
	var venues []entity.Venue
	var total int64

	db := r.buildFilterQuery(ctx, q.VenueFilter)

	if err := db.Model(&entity.Venue{}).Count(&total).Error; err != nil {
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

	err := db.Limit(limit).Offset(offset).Find(&venues).Error
	return venues, total, err
}

func (r *venueRepo) CursorVenues(ctx context.Context, q models.VenueQueryCursor) ([]entity.Venue, string, error) {
	var venues []entity.Venue
	db := r.buildFilterQuery(ctx, q.VenueFilter)

	if q.Cursor != nil && *q.Cursor != "" {
		if cursorID, err := uuid.Parse(*q.Cursor); err == nil {
			var cursorVenue entity.Venue
			if err := r.db.Select("created_at").First(&cursorVenue, "id = ?", cursorID).Error; err == nil {
				db = db.Where("(created_at, id) < (?, ?)", cursorVenue.CreatedAt, cursorID)
			}
		}
	}

	limit := 10
	if q.Limit != nil && *q.Limit > 0 {
		limit = *q.Limit
	}

	err := db.Order("created_at desc, id desc").
		Limit(limit + 1).
		Find(&venues).Error

	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(venues) > limit {
		venues = venues[:limit]
		nextCursor = venues[len(venues)-1].ID.String()
	}

	return venues, nextCursor, nil
}

// --- HELPER: Query Builder ---
func (r *venueRepo) buildFilterQuery(ctx context.Context, f models.VenueFilter) *gorm.DB {
	db := r.db.WithContext(ctx)

	if f.OrganizationID != nil {
		db = db.Where("organization_id = ?", *f.OrganizationID)
	}

	if f.Name != nil {
		db = db.Where("name ILIKE ?", "%"+*f.Name+"%")
	}

	if f.Slug != nil {
		db = db.Where("slug = ?", *f.Slug)
	}

	if f.Description != nil {
		db = db.Where("description ILIKE ?", "%"+*f.Description+"%")
	}

	if f.Address != nil {
		db = db.Where("address ILIKE ?", "%"+*f.Address+"%")
	}
	if f.City != nil {
		db = db.Where("city ILIKE ?", "%"+*f.City+"%")
	}
	if f.Province != nil {
		db = db.Where("province ILIKE ?", "%"+*f.Province+"%")
	}
	if f.PostalCode != nil {
		db = db.Where("postal_code = ?", *f.PostalCode)
	}

	if f.Visibility != nil {
		db = db.Where("visibility = ?", *f.Visibility)
	}

	if f.IsLive != nil {
		if *f.IsLive {
			db = db.Where("live_revision_id IS NOT NULL")
		} else {
			db = db.Where("live_revision_id IS NULL")
		}
	}

	db = db.Preload("CoverImage")

	return db
}

func (r *venueRepo) GetLiveManifestData(venueSlug string) (*entity.Venue, error) {
	var venue entity.Venue

	err := r.db.
		Preload("LiveRevision").
		Preload("LiveRevision.Floors").
		Preload("LiveRevision.Floors.MapImage").
		Preload("LiveRevision.Floors.Nodes", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Panorama").Where("is_active = ?", true)
		}).
		Preload("LiveRevision.Floors.Nodes.Area").
		Preload("LiveRevision.Floors.Nodes.OutgoingEdges", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true)
		}).
		Where("slug = ?", venueSlug).
		First(&venue).Error

	if err != nil {
		return nil, err
	}

	if venue.LiveRevisionID == uuid.Nil {
		return nil, errors.New("venue has no published version yet")
	}

	return &venue, nil
}

func (r *venueRepo) GetDraftDataUUID(venueID uuid.UUID) (*entity.GraphRevision, error) {
	var venue entity.Venue
	if err := r.db.First(&venue, "id = ?", venueID).Error; err != nil {
		return nil, err
	}

	if venue.DraftRevisionID == nil {
		return nil, errors.New("no active draft found")
	}

	var draft entity.GraphRevision
	err := r.db.
		Preload("Floors.Nodes.Panorama").
		Preload("Floors.Nodes.OutgoingEdges").
		First(&draft, "id = ?", *venue.DraftRevisionID).Error

	return &draft, err
}
