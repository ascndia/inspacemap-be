package service

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"
	"time"

	"github.com/google/uuid"
)

type authService struct {
	userRepo       repository.UserRepository
	authRepo       repository.AuthRepository 
	orgRepo        repository.OrganizationRepository
	orgMemberRepo  repository.OrganizationMemberRepository
	invitationRepo repository.UserInvitationRepository
	roleRepo       repository.RoleRepository
}

func NewAuthService(
	userRepo repository.UserRepository,
	authRepo repository.AuthRepository,
	orgRepo repository.OrganizationRepository,
	orgMemberRepo repository.OrganizationMemberRepository,
	invitationRepo repository.UserInvitationRepository,
	roleRepo repository.RoleRepository,
) AuthService {
	return &authService{
		userRepo:       userRepo,
		authRepo:       authRepo,
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
	if req.Password != user.PasswordHash { 
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

	newUser := entity.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: req.Password, 
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
		newUser := entity.User{
			FullName:     req.FullName,
			Email:        invite.Email,
			PasswordHash: req.Password, 
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


func (s *authService) generateAuthResponse(user *entity.User) (*models.AuthResponse, error) {
	var orgs []models.OrgMemberDetail
	for _, m := range user.Memberships {
		orgs = append(orgs, models.OrgMemberDetail{
			OrganizationID: m.OrganizationID,
			Name:           m.Organization.Name,
			Slug:           m.Organization.Slug,
			RoleName:       m.Role.Name,
		})
	}

	userDetail := models.UserDetail{
		ID:            user.ID,
		Email:         user.Email,
		FullName:      user.FullName,
		AvatarURL:     user.AvatarURL,
		Organizations: orgs,
	}

	
	token := "mock-jwt-token-xyz" 
	
	return &models.AuthResponse{
		AccessToken: token,
		RefreshToken: "mock-refresh-token",
		ExpiresIn: 3600,
		User: userDetail,
	}, nil
}