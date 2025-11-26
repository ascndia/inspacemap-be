package service

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"
	"inspacemap/backend/pkg/utils"

	"github.com/google/uuid"
)

type authService struct {
	userRepo repository.UserRepository
	orgRepo  repository.OrganizationRepository
	roleRepo repository.RoleRepository
}

func NewAuthService(
	userRepo repository.UserRepository,
	orgRepo repository.OrganizationRepository,
	roleRepo repository.RoleRepository,
) AuthService {
	return &authService{
		userRepo: userRepo,
		orgRepo:  orgRepo,
		roleRepo: roleRepo,
	}
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password hash
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	token, _ := s.generateAuthResponse(user)
	return token, nil
}

func (s *authService) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {

	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}
	newOrg := entity.Organization{
		Name: req.OrganizationName,
		Slug: "generated-slug-" + uuid.NewString()[:8],
	}
	if err := s.orgRepo.Create(ctx, &newOrg); err != nil {
		return nil, err
	}
	ownerRole, err := s.roleRepo.GetByName(ctx, "Owner")
	if err != nil {
		return nil, errors.New("system error: owner role not found")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	newUser := entity.User{
		FullName:       req.FullName,
		Email:          req.Email,
		PasswordHash:   hashedPassword,
		OrganizationID: newOrg.ID,
		RoleID:         ownerRole.ID,
	}
	if err := s.userRepo.Create(ctx, &newUser); err != nil {
		return nil, err
	}

	fullUser, _ := s.userRepo.GetByEmail(ctx, newUser.Email)
	return s.generateAuthResponse(fullUser)
}

func (s *authService) generateAuthResponse(user *entity.User) (*models.AuthResponse, error) {
	var permissions []string

	// Get permissions from user's role
	for _, p := range user.Role.Permissions {
		permissions = append(permissions, p.Key)
	}

	// Generate Token with user's org and role
	token, err := utils.GenerateToken(user.ID, user.Email, user.OrganizationID, user.Role.Name, permissions)
	if err != nil {
		return nil, err
	}

	userDetail := models.UserDetail{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		AvatarURL: user.AvatarURL,
		Organization: models.OrgMemberDetail{ // Changed from Organizations array to single Organization
			OrganizationID: user.OrganizationID,
			Name:           user.Organization.Name,
			Slug:           user.Organization.Slug,
			RoleName:       user.Role.Name,
		},
	}

	return &models.AuthResponse{
		AccessToken:  token,
		RefreshToken: "mock-refresh",
		ExpiresIn:    3600 * 24,
		User:         userDetail,
	}, nil
}
