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
	err := r.db.WithContext(ctx).
		Model(&entity.Role{BaseEntity: entity.BaseEntity{ID: roleID}}).
		Association("Permissions").
		Find(&perms)
	return perms, err
}

func (r *permissionRepo) GetByUserAndOrg(ctx context.Context, userID, orgID uuid.UUID) ([]entity.Permission, error) {
	var perms []entity.Permission

	// Query: User -> Member -> Role -> RolePermission -> Permission
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Joins("JOIN organization_members om ON om.role_id = rp.role_id").
		Where("om.user_id = ? AND om.organization_id = ?", userID, orgID).
		Find(&perms).Error

	return perms, err
}

func (r *permissionRepo) GetByOrganizationID(ctx context.Context, orgID *uuid.UUID) ([]entity.Permission, error) {
	var perms []entity.Permission
	query := r.db.WithContext(ctx)

	if orgID != nil {
		query = query.Where("organization_id = ?", *orgID)
	} else {
		query = query.Where("organization_id IS NULL")
	}

	if err := query.Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

func (r *permissionRepo) GetByOrganizationIDAndKey(ctx context.Context, orgID *uuid.UUID, key string) (*entity.Permission, error) {
	var perm entity.Permission
	query := r.db.WithContext(ctx).Where("key = ?", key)

	if orgID != nil {
		query = query.Where("organization_id = ?", *orgID)
	} else {
		query = query.Where("organization_id IS NULL")
	}

	if err := query.First(&perm).Error; err != nil {
		return nil, err
	}
	return &perm, nil
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
