package repository

import (
	"context"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type organizationRepo struct {
	BaseRepository[entity.Organization, uuid.UUID]
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepo{
		BaseRepository: NewBaseRepository[entity.Organization, uuid.UUID](db),
		db:             db,
	}
}

// GetByDomain: Mencari single organization berdasarkan identifier domain (Slug / Website)
func (r *organizationRepo) GetByDomain(ctx context.Context, domain string) (*entity.Organization, error) {
	var org entity.Organization
	// Cek apakah domain cocok dengan Slug ATAU Website
	err := r.db.WithContext(ctx).
		Where("slug = ? OR website = ?", domain, domain).
		First(&org).Error
	
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepo) FilterOrganizations(ctx context.Context, filter models.OrganizationFilter) ([]entity.Organization, error) {
	var orgs []entity.Organization
	query := r.buildFilterQuery(ctx, filter)
	err := query.Find(&orgs).Error
	return orgs, err
}

func (r *organizationRepo) PagedOrganizations(ctx context.Context, q models.OrganizationQuery) ([]entity.Organization, int64, error) {
	var orgs []entity.Organization
	var total int64

	db := r.buildFilterQuery(ctx, q.OrganizationFilter)

	if err := db.Model(&entity.Organization{}).Count(&total).Error; err != nil {
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

	err := db.Limit(limit).Offset(offset).Find(&orgs).Error
	return orgs, total, err
}

func (r *organizationRepo) CursorOrganizations(ctx context.Context, q models.OrganizationQueryCursor) ([]entity.Organization, string, error) {
	var orgs []entity.Organization
	db := r.buildFilterQuery(ctx, q.OrganizationFilter)

	// Cursor Logic: (CreatedAt, ID)
	if q.Cursor != nil && *q.Cursor != "" {
		if cursorID, err := uuid.Parse(*q.Cursor); err == nil {
			var cursorOrg entity.Organization
			if err := r.db.Select("created_at").First(&cursorOrg, "id = ?", cursorID).Error; err == nil {
				db = db.Where("(created_at, id) < (?, ?)", cursorOrg.CreatedAt, cursorID)
			}
		}
	}

	limit := 10
	if q.Limit != nil && *q.Limit > 0 {
		limit = *q.Limit
	}

	// Deterministik Sort
	err := db.Order("created_at desc, id desc").
		Limit(limit + 1).
		Find(&orgs).Error

	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(orgs) > limit {
		orgs = orgs[:limit]
		nextCursor = orgs[len(orgs)-1].ID.String()
	}

	return orgs, nextCursor, nil
}

// --- HELPER ---
func (r *organizationRepo) buildFilterQuery(ctx context.Context, f models.OrganizationFilter) *gorm.DB {
	db := r.db.WithContext(ctx)

	if f.Name != nil {
		db = db.Where("name ILIKE ?", "%"+*f.Name+"%")
	}
	
	// Jika filter Domain diisi, cari di slug atau website
	if f.Domain != nil {
		db = db.Where("slug ILIKE ? OR website ILIKE ?", "%"+*f.Domain+"%", "%"+*f.Domain+"%")
	}

	if f.Slug != nil {
		db = db.Where("slug = ?", *f.Slug)
	}

	if f.Website != nil {
		db = db.Where("website ILIKE ?", "%"+*f.Website+"%")
	}

	if f.IsActive != nil {
		db = db.Where("is_active = ?", *f.IsActive)
	}

	return db
}