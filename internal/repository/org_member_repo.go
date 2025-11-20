package repository

import (
	"context"
	"inspacemap/backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orgMemberRepo struct {
	BaseRepository[entity.OrganizationMember, uuid.UUID]
	db *gorm.DB
}

func NewOrganizationMemberRepository(db *gorm.DB) OrganizationMemberRepository {
	return &orgMemberRepo{
		BaseRepository: NewBaseRepository[entity.OrganizationMember, uuid.UUID](db),
		db:             db,
	}
}

func (r *orgMemberRepo) AddMember(ctx context.Context, member *entity.OrganizationMember) error {
	return r.db.WithContext(ctx).
		FirstOrCreate(member, entity.OrganizationMember{
			OrganizationID: member.OrganizationID,
			UserID:         member.UserID,
		}).Error
}

func (r *orgMemberRepo) RemoveMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Delete(&entity.OrganizationMember{}).Error
}

func (r *orgMemberRepo) UpdateRole(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, roleID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.OrganizationMember{}).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Update("role_id", roleID).Error
}

func (r *orgMemberRepo) GetMember(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) (*entity.OrganizationMember, error) {
	var member entity.OrganizationMember
	err := r.db.WithContext(ctx).
		Preload("Role.Permissions"). // Load permission agar middleware bisa cek akses
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		First(&member).Error
	return &member, err
}

func (r *orgMemberRepo) GetMembersByOrg(ctx context.Context, orgID uuid.UUID) ([]entity.OrganizationMember, error) {
	var members []entity.OrganizationMember
	err := r.db.WithContext(ctx).
		Preload("User"). // Load nama user
		Preload("Role"). // Load jabatan
		Where("organization_id = ?", orgID).
		Find(&members).Error
	return members, err
}

func (r *orgMemberRepo) GetMembersByUser(ctx context.Context, userID uuid.UUID) ([]entity.OrganizationMember, error) {
	var memberships []entity.OrganizationMember
	err := r.db.WithContext(ctx).
		Preload("Organization"). // Load nama organisasi
		Preload("Role").
		Where("user_id = ?", userID).
		Find(&memberships).Error
	return memberships, err
}