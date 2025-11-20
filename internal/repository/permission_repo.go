package repository

import (
	"context"
	"inspacemap/backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)
type permissionRepo struct {
	BaseRepository[entity.Permission, uuid.UUID]
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepo{
		BaseRepository: NewBaseRepository[entity.Permission, uuid.UUID](db),
		db:             db,
	}
}

func (r *permissionRepo) GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error) {
	var perms []entity.Permission
	query := r.db.WithContext(ctx).Where("role_id = ?", roleID)

	if err := query.Find(&perms).Error; err != nil {
		return nil, err
	}	
	return perms, nil
}

func (r *permissionRepo) GetByKey(ctx context.Context, key string) (*entity.Permission, error) {
	var perm entity.Permission
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&perm).Error; err != nil {
		return nil, err
	}
	return &perm, nil
}

func (r *permissionRepo) GetAll(ctx context.Context) ([]entity.Permission, error) {
	var perms []entity.Permission
	if err := r.db.WithContext(ctx).Order("group asc, key asc").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}