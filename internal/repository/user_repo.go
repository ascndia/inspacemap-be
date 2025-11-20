package repository

import (
	"context"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)
type userRepo struct {
	BaseRepository[entity.User, uuid.UUID] 
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{
		BaseRepository: NewBaseRepository[entity.User, uuid.UUID](db),
		db:             db,
	}
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Preload("Memberships").
		Preload("Memberships.Organization").
		Preload("Memberships.Role").
		Where("email = ?", email).
		First(&user).Error
	
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.User, error) {
	var users []entity.User
	err := r.db.WithContext(ctx).
		Joins("JOIN organization_members om ON om.user_id = users.id").
		Where("om.organization_id = ?", orgID).
		Preload("Memberships", "organization_id = ?", orgID).
		Preload("Memberships.Role").
		Find(&users).Error
	return users, err
}

func (r *userRepo) FilterUsers(ctx context.Context, filter models.UserFilter) ([]entity.User, error) {
	var users []entity.User
	query := r.buildFilterQuery(ctx, filter)
	err := query.Find(&users).Error
	return users, err
}

func (r *userRepo) PagedUsers(ctx context.Context, q models.UserQuery) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	db := r.buildFilterQuery(ctx, q.UserFilter)
	if err := db.Model(&entity.User{}).Count(&total).Error; err != nil {
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
	
	err := db.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *userRepo) CursorUsers(ctx context.Context, q models.UserQueryCursor) ([]entity.User, string, error) {
	var users []entity.User
	
	db := r.buildFilterQuery(ctx, q.UserFilter)
	
	if q.Cursor != nil && *q.Cursor != "" {
		if cursorID, err := uuid.Parse(*q.Cursor); err == nil {
			var cursorUser entity.User
			if err := r.db.Select("created_at").First(&cursorUser, "id = ?", cursorID).Error; err == nil {
				db = db.Where("(created_at, id) < (?, ?)", cursorUser.CreatedAt, cursorID)
			}
		}
	}

	limit := 10
	if q.Limit != nil && *q.Limit > 0 {
		limit = *q.Limit
	}

	err := db.Order("created_at desc, id desc").
		Limit(limit + 1). 
		Find(&users).Error

	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(users) > limit {
		users = users[:limit]
		nextCursor = users[len(users)-1].ID.String()
	}

	return users, nextCursor, nil
}

func (r *userRepo) buildFilterQuery(ctx context.Context, f models.UserFilter) *gorm.DB {
	db := r.db.WithContext(ctx)

	if f.OrganizationID != nil {
		db = db.Joins("JOIN organization_members om ON om.user_id = users.id").
			Where("om.organization_id = ?", *f.OrganizationID)
		
		if f.RoleID != nil {
			db = db.Where("om.role_id = ?", *f.RoleID)
		}
	}

	if f.Email != nil {
		db = db.Where("email ILIKE ?", "%"+*f.Email+"%")
	}

	if f.FullName != nil {
		db = db.Where("full_name ILIKE ?", "%"+*f.FullName+"%")
	}

	db = db.Preload("Memberships").Preload("Memberships.Role")

	return db
}