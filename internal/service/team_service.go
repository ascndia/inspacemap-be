package service

import (
	"context"
	"errors"
	"fmt"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"
	"strings"
	"time"

	"github.com/google/uuid"
)

type teamService struct {
	userRepo       repository.UserRepository
	invitationRepo repository.UserInvitationRepository
	orgMemberRepo  repository.OrganizationMemberRepository
	roleRepo       repository.RoleRepository
}

func NewTeamService(
	userRepo repository.UserRepository,
	invitationRepo repository.UserInvitationRepository,
	orgMemberRepo repository.OrganizationMemberRepository,
	roleRepo repository.RoleRepository,
) TeamService {
	return &teamService{
		userRepo:       userRepo,
		invitationRepo: invitationRepo,
		orgMemberRepo:  orgMemberRepo,
		roleRepo:       roleRepo,
	}
}

func (s *teamService) InviteMember(ctx context.Context, orgID uuid.UUID, inviterID uuid.UUID, req models.InviteUserRequest) error {

	role, err := s.roleRepo.GetByID(ctx, req.RoleID)
	if err != nil {
		return errors.New("role tidak valid atau tidak ditemukan")
	}

	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {

		memberCheck, _ := s.orgMemberRepo.GetMember(ctx, orgID, existingUser.ID)
		if memberCheck != nil {
			return fmt.Errorf("user %s sudah menjadi anggota tim", req.Email)
		}
	}

	existingInvites, _ := s.invitationRepo.GetByEmail(ctx, req.Email)
	for _, inv := range existingInvites {
		if inv.OrganizationID == orgID && inv.Status == "pending" {

			return errors.New("user ini sudah diundang dan statusnya masih pending")
		}
	}

	token := uuid.NewString()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	invite := entity.UserInvitation{
		OrganizationID:  orgID,
		Email:           req.Email,
		RoleID:          role.ID,
		Token:           token,
		ExpiresAt:       expiresAt,
		InvitedByUserID: inviterID,
		Status:          "pending",
	}

	if err := s.invitationRepo.Create(ctx, &invite); err != nil {
		return errors.New("gagal membuat undangan")
	}

	return nil
}

func (s *teamService) RemoveMember(ctx context.Context, orgID uuid.UUID, targetUserID uuid.UUID) error {

	targetMember, err := s.orgMemberRepo.GetMember(ctx, orgID, targetUserID)
	if err != nil {
		return errors.New("member tidak ditemukan")
	}

	if strings.EqualFold(targetMember.Role.Name, "Owner") {
		if err := s.ensureOrganizationHasOtherOwner(ctx, orgID, targetUserID); err != nil {
			return err
		}
	}

	return s.orgMemberRepo.RemoveMember(ctx, orgID, targetUserID)
}

func (s *teamService) UpdateMemberRole(ctx context.Context, orgID uuid.UUID, req models.UpdateUserRoleRequest) error {

	if _, err := s.roleRepo.GetByID(ctx, req.NewRoleID); err != nil {
		return errors.New("role baru tidak valid")
	}

	targetMember, err := s.orgMemberRepo.GetMember(ctx, orgID, req.TargetUserID)
	if err != nil {
		return errors.New("member tidak ditemukan")
	}

	if strings.EqualFold(targetMember.Role.Name, "Owner") && targetMember.RoleID != req.NewRoleID {
		if err := s.ensureOrganizationHasOtherOwner(ctx, orgID, req.TargetUserID); err != nil {
			return err
		}
	}

	return s.orgMemberRepo.UpdateRole(ctx, orgID, req.TargetUserID, req.NewRoleID)
}

func (s *teamService) GetMembersList(ctx context.Context, orgID uuid.UUID) ([]models.TeamMemberDetail, error) {
	memberships, err := s.orgMemberRepo.GetMembersByOrg(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var details []models.TeamMemberDetail
	for _, m := range memberships {
		details = append(details, models.TeamMemberDetail{
			UserID:    m.UserID,
			FullName:  m.User.FullName,
			Email:     m.User.Email,
			AvatarURL: m.User.AvatarURL,
			RoleID:    m.RoleID,
			RoleName:  m.Role.Name,
			JoinedAt:  m.JoinedAt,
		})
	}

	return details, nil
}

func (s *teamService) ensureOrganizationHasOtherOwner(ctx context.Context, orgID uuid.UUID, excludeUserID uuid.UUID) error {

	members, err := s.orgMemberRepo.GetMembersByOrg(ctx, orgID)
	if err != nil {
		return errors.New("gagal memvalidasi owner")
	}

	ownerCount := 0
	for _, m := range members {

		if strings.EqualFold(m.Role.Name, "Owner") && m.UserID != excludeUserID {
			ownerCount++
		}
	}

	if ownerCount == 0 {
		return errors.New("tindakan ditolak: organisasi harus memiliki setidaknya satu owner")
	}

	return nil
}
