package repository

import (
	"context"
	"inspacemap/backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =================================================================
// ROLE REPOSITORY (READ ONLY FOR APP USAGE)
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
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetAll: Dipakai untuk dropdown "Pilih Jabatan" saat Invite Member
func (r *roleRepo) GetAll(ctx context.Context) ([]entity.Role, error) {
	var roles []entity.Role
	// Preload permissions tidak wajib untuk list dropdown, tapi berguna untuk debug
	err := r.db.WithContext(ctx).Order("name asc").Find(&roles).Error
	return roles, err
}

func (r *roleRepo) GetPermissions(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error) {
	var perms []entity.Permission
	err := r.db.WithContext(ctx).
		Model(&entity.Role{BaseModel: entity.BaseModel{ID: roleID}}).
		Association("Permissions").
		Find(&perms)
	return perms, err
}

// =================================================================
// PERMISSION REPOSITORY (HELPER)
// =================================================================

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

// GetByUserAndOrg: Method KUNCI untuk Middleware RBAC
func (r *permissionRepo) GetByUserAndOrg(ctx context.Context, userID, orgID uuid.UUID) ([]entity.Permission, error) {
	var perms []entity.Permission
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Joins("JOIN organization_members om ON om.role_id = rp.role_id").
		Where("om.user_id = ? AND om.organization_id = ?", userID, orgID).
		Distinct().
		Find(&perms).Error
	return perms, err
}
