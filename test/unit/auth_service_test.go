package unit

import (
	"context"
	"errors"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type AuthServiceTestSuite struct {
	suite.Suite
	ctrl           *gomock.Controller
	userRepo       *MockUserRepository
	orgRepo        *MockOrganizationRepository
	orgMemberRepo  *MockOrganizationMemberRepository
	invitationRepo *MockUserInvitationRepository
	roleRepo       *MockRoleRepository
	authService    service.AuthService
}

func (suite *AuthServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.userRepo = NewMockUserRepository(suite.ctrl)
	suite.orgRepo = NewMockOrganizationRepository(suite.ctrl)
	suite.orgMemberRepo = NewMockOrganizationMemberRepository(suite.ctrl)
	suite.invitationRepo = NewMockUserInvitationRepository(suite.ctrl)
	suite.roleRepo = NewMockRoleRepository(suite.ctrl)

	suite.authService = service.NewAuthService(
		suite.userRepo,
		suite.orgRepo,
		suite.orgMemberRepo,
		suite.invitationRepo,
		suite.roleRepo,
	)
}

func (suite *AuthServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (suite *AuthServiceTestSuite) TestRegister_EmailAlreadyExists() {
	ctx := context.Background()
	req := models.RegisterRequest{
		FullName:         "John Doe",
		Email:            "existing@example.com",
		Password:         "password123",
		OrganizationName: "Test Org",
	}

	existingUser := &entity.User{
		BaseEntity: entity.BaseEntity{
			ID: uuid.New(),
		},
		Email: req.Email,
	}

	// Mock expectations
	suite.userRepo.EXPECT().GetByEmail(ctx, req.Email).Return(existingUser, nil)

	// Execute
	result, err := suite.authService.Register(ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "email already registered", err.Error())
}

func (suite *AuthServiceTestSuite) TestRegister_OrganizationCreationFails() {
	ctx := context.Background()
	req := models.RegisterRequest{
		FullName:         "John Doe",
		Email:            "john@example.com",
		Password:         "password123",
		OrganizationName: "Test Org",
	}

	// Mock expectations
	suite.userRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, nil)
	suite.orgRepo.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("database error"))

	// Execute
	result, err := suite.authService.Register(ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "database error", err.Error())
}

func (suite *AuthServiceTestSuite) TestRegister_OwnerRoleNotFound() {
	ctx := context.Background()
	req := models.RegisterRequest{
		FullName:         "John Doe",
		Email:            "john@example.com",
		Password:         "password123",
		OrganizationName: "Test Org",
	}

	// Mock expectations
	suite.userRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, nil)
	suite.orgRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
	suite.roleRepo.EXPECT().GetByName(ctx, "Owner").Return(nil, errors.New("role not found"))

	// Execute
	result, err := suite.authService.Register(ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "system error: owner role not found", err.Error())
}

func (suite *AuthServiceTestSuite) TestRegister_Success() {
	ctx := context.Background()
	req := models.RegisterRequest{
		FullName:         "John Doe",
		Email:            "john@example.com",
		Password:         "password123",
		OrganizationName: "Test Org",
	}

	orgID := uuid.New()
	userID := uuid.New()
	roleID := uuid.New()

	// Mock expectations
	suite.userRepo.EXPECT().GetByEmail(ctx, req.Email).Return(nil, nil)
	suite.orgRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, org *entity.Organization) error {
		org.ID = orgID
		return nil
	})
	suite.roleRepo.EXPECT().GetByName(ctx, "Owner").Return(&entity.Role{
		BaseEntity: entity.BaseEntity{ID: roleID},
		Name:       "Owner",
		Permissions: []entity.Permission{
			{Key: "read:organization"},
		},
	}, nil)
	suite.userRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *entity.User) error {
		user.ID = userID
		return nil
	})
	suite.orgMemberRepo.EXPECT().AddMember(ctx, gomock.Any()).Return(nil)
	suite.userRepo.EXPECT().GetByEmail(ctx, req.Email).Return(&entity.User{
		BaseEntity: entity.BaseEntity{ID: userID},
		FullName:   req.FullName,
		Email:      req.Email,
		Memberships: []entity.OrganizationMember{
			{
				OrganizationID: orgID,
				Organization: entity.Organization{
					BaseEntity: entity.BaseEntity{ID: orgID},
					Name:       req.OrganizationName,
					Slug:       "generated-slug-" + uuid.NewString()[:8],
				},
				Role: entity.Role{
					BaseEntity: entity.BaseEntity{ID: roleID},
					Name:       "Owner",
					Permissions: []entity.Permission{
						{Key: "read:organization"},
					},
				},
			},
		},
	}, nil)

	// Execute
	result, err := suite.authService.Register(ctx, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "John Doe", result.User.FullName)
	assert.Equal(suite.T(), "john@example.com", result.User.Email)
	assert.NotEmpty(suite.T(), result.AccessToken)
	assert.Equal(suite.T(), int(86400), result.ExpiresIn)
}

func (suite *AuthServiceTestSuite) TestAcceptInvitation_Success_ExistingUser() {
	ctx := context.Background()
	req := models.AcceptInviteRequest{
		Token:    "valid-invitation-token",
		Password: "newpassword123",
		FullName: "John Doe",
	}

	orgID := uuid.New()
	userID := uuid.New()
	roleID := uuid.New()

	invitation := &entity.UserInvitation{
		BaseEntity:     entity.BaseEntity{ID: uuid.New()},
		OrganizationID: orgID,
		Email:          "john@example.com",
		RoleID:         roleID,
		Token:          req.Token,
		ExpiresAt:      time.Now().Add(time.Hour), // Not expired
	}

	existingUser := &entity.User{
		BaseEntity: entity.BaseEntity{ID: userID},
		Email:      invitation.Email,
		FullName:   "Existing User",
	}

	// Mock expectations
	suite.invitationRepo.EXPECT().GetByToken(ctx, req.Token).Return(invitation, nil)
	suite.userRepo.EXPECT().GetByEmail(ctx, invitation.Email).Return(existingUser, nil)
	suite.orgMemberRepo.EXPECT().AddMember(ctx, gomock.Any()).Return(nil)
	suite.userRepo.EXPECT().GetByEmail(ctx, invitation.Email).Return(existingUser, nil)

	// Execute
	result, err := suite.authService.AcceptInvitation(ctx, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "Existing User", result.User.FullName)
	assert.Equal(suite.T(), "john@example.com", result.User.Email)
	assert.NotEmpty(suite.T(), result.AccessToken)
}

func (suite *AuthServiceTestSuite) TestAcceptInvitation_Success_NewUser() {
	ctx := context.Background()
	req := models.AcceptInviteRequest{
		Token:    "valid-invitation-token",
		Password: "newpassword123",
		FullName: "John Doe",
	}

	orgID := uuid.New()
	newUserID := uuid.New()
	roleID := uuid.New()

	invitation := &entity.UserInvitation{
		BaseEntity:     entity.BaseEntity{ID: uuid.New()},
		OrganizationID: orgID,
		Email:          "john@example.com",
		RoleID:         roleID,
		Token:          req.Token,
		ExpiresAt:      time.Now().Add(time.Hour), // Not expired
	}

	// Mock expectations
	suite.invitationRepo.EXPECT().GetByToken(ctx, req.Token).Return(invitation, nil)
	suite.userRepo.EXPECT().GetByEmail(ctx, invitation.Email).Return(nil, nil) // No existing user
	suite.userRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *entity.User) error {
		user.ID = newUserID
		user.IsEmailVerified = true
		return nil
	})
	suite.orgMemberRepo.EXPECT().AddMember(ctx, gomock.Any()).Return(nil)
	suite.userRepo.EXPECT().GetByEmail(ctx, invitation.Email).Return(&entity.User{
		BaseEntity:      entity.BaseEntity{ID: newUserID},
		FullName:        req.FullName,
		Email:           invitation.Email,
		IsEmailVerified: true,
		Memberships: []entity.OrganizationMember{
			{
				OrganizationID: orgID,
				Organization: entity.Organization{
					BaseEntity: entity.BaseEntity{ID: orgID},
					Name:       "Test Org",
					Slug:       "test-org",
				},
				Role: entity.Role{
					BaseEntity: entity.BaseEntity{ID: roleID},
					Name:       "Editor",
					Permissions: []entity.Permission{
						{Key: "read:organization"},
					},
				},
			},
		},
	}, nil)

	// Execute
	result, err := suite.authService.AcceptInvitation(ctx, req)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "John Doe", result.User.FullName)
	assert.Equal(suite.T(), "john@example.com", result.User.Email)
	assert.NotEmpty(suite.T(), result.AccessToken)
}

func (suite *AuthServiceTestSuite) TestAcceptInvitation_InvalidToken() {
	ctx := context.Background()
	req := models.AcceptInviteRequest{
		Token:    "invalid-token",
		Password: "password123",
		FullName: "John Doe",
	}

	// Mock expectations
	suite.invitationRepo.EXPECT().GetByToken(ctx, req.Token).Return(nil, errors.New("invitation not found"))

	// Execute
	result, err := suite.authService.AcceptInvitation(ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "invalid or expired invitation token", err.Error())
}

func (suite *AuthServiceTestSuite) TestAcceptInvitation_ExpiredToken() {
	ctx := context.Background()
	req := models.AcceptInviteRequest{
		Token:    "expired-token",
		Password: "password123",
		FullName: "John Doe",
	}

	expiredInvitation := &entity.UserInvitation{
		BaseEntity: entity.BaseEntity{ID: uuid.New()},
		Token:      req.Token,
		ExpiresAt:  time.Now().Add(-time.Hour), // Already expired
	}

	// Mock expectations
	suite.invitationRepo.EXPECT().GetByToken(ctx, req.Token).Return(expiredInvitation, nil)

	// Execute
	result, err := suite.authService.AcceptInvitation(ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "invitation expired", err.Error())
}

func (suite *AuthServiceTestSuite) TestAcceptInvitation_AddMemberFails() {
	ctx := context.Background()
	req := models.AcceptInviteRequest{
		Token:    "valid-token",
		Password: "password123",
		FullName: "John Doe",
	}

	orgID := uuid.New()
	roleID := uuid.New()

	invitation := &entity.UserInvitation{
		BaseEntity:     entity.BaseEntity{ID: uuid.New()},
		OrganizationID: orgID,
		Email:          "john@example.com",
		RoleID:         roleID,
		Token:          req.Token,
		ExpiresAt:      time.Now().Add(time.Hour), // Not expired
	}

	existingUser := &entity.User{
		BaseEntity: entity.BaseEntity{ID: uuid.New()},
		Email:      invitation.Email,
	}

	// Mock expectations
	suite.invitationRepo.EXPECT().GetByToken(ctx, req.Token).Return(invitation, nil)
	suite.userRepo.EXPECT().GetByEmail(ctx, invitation.Email).Return(existingUser, nil)
	suite.orgMemberRepo.EXPECT().AddMember(ctx, gomock.Any()).Return(errors.New("database error"))

	// Execute
	result, err := suite.authService.AcceptInvitation(ctx, req)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "database error", err.Error())
}
