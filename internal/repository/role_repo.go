package repository

import (
	"context"
	"inspacemap/backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// =================================================================
// ROLE REPOSITORY
// =================================================================

type roleRepo struct {
	BaseRepository[entity.Role, uuid.UUID]
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepo{
		BaseRepository: NewBaseRepository[entity.Role, uuid.UUID](db),
		db:             db,
	}
}

func (r *roleRepo) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	var role entity.Role
	query := r.db.WithContext(ctx).Where("name = ?", name)

	if err := query.First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepo) GetByOrganizationID(ctx context.Context, orgID *uuid.UUID) ([]entity.Role, error) {
	var roles []entity.Role
	query := r.db.WithContext(ctx)

	if orgID != nil {
		query = query.Where("organization_id = ?", *orgID)
	} else {
		query = query.Where("organization_id IS NULL")
	}

	if err := query.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepo) GetByOrganizationIDAndName(ctx context.Context, orgID *uuid.UUID, name string) (*entity.Role, error) {
	var role entity.Role
	query := r.db.WithContext(ctx).Where("name = ?", name)

	if orgID != nil {
		query = query.Where("organization_id = ?", *orgID)
	} else {
		query = query.Where("organization_id IS NULL")
	}

	if err := query.First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepo) AttachPermission(ctx context.Context, roleID uuid.UUID, permID uuid.UUID) error {
	pivot := entity.RolePermission{
		RoleID:       roleID,
		PermissionID: permID,
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&pivot).Error
}

func (r *roleRepo) DetachPermission(ctx context.Context, roleID uuid.UUID, permID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permID).
		Delete(&entity.RolePermission{}).Error
}

func (r *roleRepo) GetPermissions(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error) {
	var permissions []entity.Permission

	// Join otomatis via GORM Association
	// SELECT * FROM permissions JOIN role_permissions ON ...
	err := r.db.WithContext(ctx).
		Model(&entity.Role{BaseEntity: entity.BaseEntity{ID: roleID}}).
		Association("Permissions").
		Find(&permissions)

	return permissions, err
}
