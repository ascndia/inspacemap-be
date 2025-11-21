package service

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"
	"inspacemap/backend/pkg/utils"
	"time"

	"github.com/google/uuid"
)

type authService struct {
	userRepo       repository.UserRepository
	orgRepo        repository.OrganizationRepository
	orgMemberRepo  repository.OrganizationMemberRepository
	invitationRepo repository.UserInvitationRepository
	roleRepo       repository.RoleRepository
}

func NewAuthService(
	userRepo repository.UserRepository,
	orgRepo repository.OrganizationRepository,
	orgMemberRepo repository.OrganizationMemberRepository,
	invitationRepo repository.UserInvitationRepository,
	roleRepo repository.RoleRepository,
) AuthService {
	return &authService{
		userRepo:       userRepo,
		orgRepo:        orgRepo,
		orgMemberRepo:  orgMemberRepo,
		invitationRepo: invitationRepo,
		roleRepo:       roleRepo,
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
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}
	if err := s.userRepo.Create(ctx, &newUser); err != nil {
		return nil, err
	}

	member := entity.OrganizationMember{
		OrganizationID: newOrg.ID,
		UserID:         newUser.ID,
		RoleID:         ownerRole.ID,
	}
	if err := s.orgMemberRepo.AddMember(ctx, &member); err != nil {
		return nil, err
	}

	fullUser, _ := s.userRepo.GetByEmail(ctx, newUser.Email)
	return s.generateAuthResponse(fullUser)
}

func (s *authService) AcceptInvitation(ctx context.Context, req models.AcceptInviteRequest) (*models.AuthResponse, error) {

	invite, err := s.invitationRepo.GetByToken(ctx, req.Token)
	if err != nil {
		return nil, errors.New("invalid or expired invitation token")
	}

	if invite.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("invitation expired")
	}

	var targetUserID uuid.UUID
	existingUser, _ := s.userRepo.GetByEmail(ctx, invite.Email)

	if existingUser != nil {

		targetUserID = existingUser.ID
	} else {
		// Hash password for new user
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		
		newUser := entity.User{
			FullName:        req.FullName,
			Email:           invite.Email,
			PasswordHash:    hashedPassword,
			IsEmailVerified: true,
		}
		if err := s.userRepo.Create(ctx, &newUser); err != nil {
			return nil, err
		}
		targetUserID = newUser.ID
	}
	member := entity.OrganizationMember{
		OrganizationID: invite.OrganizationID,
		UserID:         targetUserID,
		RoleID:         invite.RoleID,
	}
	if err := s.orgMemberRepo.AddMember(ctx, &member); err != nil {
		return nil, err
	}
	fullUser, _ := s.userRepo.GetByEmail(ctx, invite.Email)
	return s.generateAuthResponse(fullUser)
}

// func (s *authService) generateAuthResponse(user *entity.User) (*models.AuthResponse, error) {
// 	var orgs []models.OrgMemberDetail
// 	for _, m := range user.Memberships {
// 		orgs = append(orgs, models.OrgMemberDetail{
// 			OrganizationID: m.OrganizationID,
// 			Name:           m.Organization.Name,
// 			Slug:           m.Organization.Slug,
// 			RoleName:       m.Role.Name,
// 		})
// 	}

// 	userDetail := models.UserDetail{
// 		ID:            user.ID,
// 		Email:         user.Email,
// 		FullName:      user.FullName,
// 		AvatarURL:     user.AvatarURL,
// 		Organizations: orgs,
// 	}

// 	token := "mock-jwt-token-xyz"

//	return &models.AuthResponse{
//		AccessToken: token,
//		RefreshToken: "mock-refresh-token",
//		ExpiresIn: 3600,
//		User: userDetail,
//	}, nil
func (s *authService) generateAuthResponse(user *entity.User) (*models.AuthResponse, error) {
	var orgs []models.OrgMemberDetail

	// Cari Org Default (yang pertama) untuk dijadikan konteks token awal
	var defaultOrgID uuid.UUID
	var defaultRoleName string
	var permissions []string

	for i, m := range user.Memberships {
		orgs = append(orgs, models.OrgMemberDetail{
			OrganizationID: m.OrganizationID,
			Name:           m.Organization.Name,
			Slug:           m.Organization.Slug,
			RoleName:       m.Role.Name,
		})

		// Set default ke index 0
		if i == 0 {
			defaultOrgID = m.OrganizationID
			defaultRoleName = m.Role.Name

			// Ambil Permission Strings dari Role
			for _, p := range m.Role.Permissions {
				permissions = append(permissions, p.Key)
			}
		}
	}

	// Generate Token yang sudah "dibumbui" Permission
	token, err := utils.GenerateToken(user.ID, user.Email, defaultOrgID, defaultRoleName, permissions)
	if err != nil {
		return nil, err
	}

	userDetail := models.UserDetail{
		ID:            user.ID,
		Email:         user.Email,
		FullName:      user.FullName,
		AvatarURL:     user.AvatarURL,
		Organizations: orgs,
	}

	return &models.AuthResponse{
		AccessToken:  token,
		RefreshToken: "mock-refresh",
		ExpiresIn:    3600 * 24,
		User:         userDetail,
	}, nil
}
