package service

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"
	"inspacemap/backend/pkg/utils"
	"strings"

	"github.com/google/uuid"
)

type teamService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewTeamService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
) TeamService {
	return &teamService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *teamService) CreateUser(ctx context.Context, orgID uuid.UUID, req models.CreateUserRequest) error {
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return errors.New("email already registered")
	}

	role, err := s.roleRepo.GetByID(ctx, req.RoleID)
	if err != nil {
		return errors.New("invalid role")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	newUser := entity.User{
		FullName:       req.FullName,
		Email:          req.Email,
		PasswordHash:   hashedPassword,
		OrganizationID: orgID,
		RoleID:         role.ID,
	}
	return s.userRepo.Create(ctx, &newUser)
}

func (s *teamService) RemoveMember(ctx context.Context, orgID uuid.UUID, targetUserID uuid.UUID) error {

	targetUser, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return errors.New("user not found")
	}

	if targetUser.OrganizationID != orgID {
		return errors.New("user not in this organization")
	}

	if strings.EqualFold(targetUser.Role.Name, "Owner") {
		if err := s.ensureOrganizationHasOtherOwner(ctx, orgID, targetUserID); err != nil {
			return err
		}
	}

	return s.userRepo.Delete(ctx, targetUserID)
}

func (s *teamService) UpdateMemberRole(ctx context.Context, orgID uuid.UUID, req models.UpdateUserRoleRequest) error {

	if _, err := s.roleRepo.GetByID(ctx, req.NewRoleID); err != nil {
		return errors.New("role baru tidak valid")
	}

	targetUser, err := s.userRepo.GetByID(ctx, req.TargetUserID)
	if err != nil {
		return errors.New("user not found")
	}

	if targetUser.OrganizationID != orgID {
		return errors.New("user not in this organization")
	}

	if strings.EqualFold(targetUser.Role.Name, "Owner") && targetUser.RoleID != req.NewRoleID {
		if err := s.ensureOrganizationHasOtherOwner(ctx, orgID, req.TargetUserID); err != nil {
			return err
		}
	}

	return s.userRepo.UpdateRole(ctx, req.TargetUserID, req.NewRoleID)
}

func (s *teamService) GetMembersList(ctx context.Context, orgID uuid.UUID) ([]models.TeamMemberDetail, error) {
	users, err := s.userRepo.GetByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var details []models.TeamMemberDetail
	for _, user := range users {
		details = append(details, models.TeamMemberDetail{
			UserID:    user.ID,
			FullName:  user.FullName,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
			RoleID:    user.RoleID,
			RoleName:  user.Role.Name,
			JoinedAt:  user.CreatedAt, // Use CreatedAt as JoinedAt
		})
	}

	return details, nil
}

func (s *teamService) ensureOrganizationHasOtherOwner(ctx context.Context, orgID uuid.UUID, excludeUserID uuid.UUID) error {

	users, err := s.userRepo.GetByOrganizationID(ctx, orgID)
	if err != nil {
		return errors.New("failed to get users")
	}

	ownerCount := 0
	for _, user := range users {
		if strings.EqualFold(user.Role.Name, "Owner") && user.ID != excludeUserID {
			ownerCount++
		}
	}

	if ownerCount == 0 {
		return errors.New("organization must have at least one owner")
	}

	return nil
}
