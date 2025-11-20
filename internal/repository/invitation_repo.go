package repository

import (
	"context"
	"inspacemap/backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type invitationRepo struct {
	BaseRepository[entity.UserInvitation, uuid.UUID]
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) UserInvitationRepository {
	return &invitationRepo{
		BaseRepository: NewBaseRepository[entity.UserInvitation, uuid.UUID](db),
		db:             db,
	}
}


func (r *invitationRepo) GetByToken(ctx context.Context, token string) (*entity.UserInvitation, error) {
	var invite entity.UserInvitation
	

	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Organization").
		Where("token = ? AND status = ?", token, "pending"). 
		First(&invite).Error

	if err != nil {
		return nil, err
	}
	return &invite, nil
}

func (r *invitationRepo) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]entity.UserInvitation, error) {
	var invites []entity.UserInvitation
	err := r.db.WithContext(ctx).
		Preload("Role"). 
		Where("organization_id = ?", orgID).
		Order("created_at desc").
		Find(&invites).Error
	return invites, err
}

func (r *invitationRepo) GetByEmail(ctx context.Context, email string) ([]entity.UserInvitation, error) {
	var invites []entity.UserInvitation
	err := r.db.WithContext(ctx).
		Where("email = ? AND status = ?", email, "pending").
		Find(&invites).Error
	return invites, err
}

func (r *invitationRepo) GetByStatus(ctx context.Context, orgID uuid.UUID, status string) ([]entity.UserInvitation, error) {
	var invites []entity.UserInvitation
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND status = ?", orgID, status).
		Find(&invites).Error
	return invites, err
}

func (r *invitationRepo) RevokeInvitation(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.UserInvitation{}).
		Where("id = ?", id).
		Update("status", "revoked").Error
}

func (r *invitationRepo) GetByRoleID(ctx context.Context, roleID uuid.UUID) ([]entity.UserInvitation, error) {
	var invites []entity.UserInvitation
	err := r.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Find(&invites).Error
	return invites, err
}

func (r *invitationRepo) GetByInviterID(ctx context.Context, inviterID uuid.UUID) ([]entity.UserInvitation, error) {
	var invites []entity.UserInvitation
	err := r.db.WithContext(ctx).
		Where("invited_by_user_id = ?", inviterID).
		Find(&invites).Error
	return invites, err
}