package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"inspacemap/backend/config"
	"inspacemap/backend/internal/delivery/http/handler"
	"inspacemap/backend/internal/delivery/http/route"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ComprehensiveHTTPTestSuite struct {
	suite.Suite
	app       *fiber.App
	db        *gorm.DB
	authSvc   service.AuthService
	venueSvc  service.VenueService
	graphSvc  service.GraphService
	teamSvc   service.TeamService
	auditSvc  service.AuditService
	userRepo  repository.UserRepository
	orgRepo   repository.OrganizationRepository
	testUser  *entity.User
	testOrg   *entity.Organization
	authToken string
}

func (suite *ComprehensiveHTTPTestSuite) SetupTest() {
	// Connect to test database
	config.ConnectDB()
	suite.db = config.DB

	// Initialize repositories
	suite.userRepo = repository.NewUserRepository(suite.db)
	suite.orgRepo = repository.NewOrganizationRepository(suite.db)
	orgMemberRepo := repository.NewOrganizationMemberRepository(suite.db)
	invitationRepo := repository.NewInvitationRepository(suite.db)
	roleRepo := repository.NewRoleRepository(suite.db)
	permRepo := repository.NewPermissionRepository(suite.db)

	venueRepo := repository.NewVenueRepository(suite.db)
	floorRepo := repository.NewFloorRepository(suite.db)
	areaRepo := repository.NewAreaRepository(suite.db)

	graphRepo := repository.NewGraphRepository(suite.db)
	revisionRepo := repository.NewGraphRevisionRepository(suite.db)

	venueGalleryRepo := repository.NewVenueGalleryRepository(suite.db)
	areaGalleryRepo := repository.NewAreaGalleryRepository(suite.db)

	auditRepo := repository.NewAuditRepository(suite.db)

	// Initialize services
	suite.authSvc = service.NewAuthService(suite.userRepo, suite.orgRepo, orgMemberRepo, invitationRepo, roleRepo)
	suite.venueSvc = service.NewVenueService(venueRepo)
	suite.graphSvc = service.NewGraphService(graphRepo, revisionRepo, floorRepo, venueRepo)
	// Skip media service for now due to storage provider complexity
	suite.teamSvc = service.NewTeamService(suite.userRepo, invitationRepo, orgMemberRepo, roleRepo)
	roleSvc := service.NewRoleService(roleRepo, permRepo)
	venueGallerySvc := service.NewVenueGalleryService(venueGalleryRepo)
	areaGallerySvc := service.NewAreaGalleryService(areaGalleryRepo)
	suite.auditSvc = service.NewAuditService(auditRepo)
	areaSvc := service.NewAreaService(areaRepo, areaGalleryRepo, graphRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(suite.authSvc)
	venueHandler := handler.NewVenueHandler(suite.venueSvc)
	areaHandler := handler.NewAreaHandler(areaSvc)
	graphHandler := handler.NewGraphHandler(suite.graphSvc)
	// Skip media handler for now
	teamRoleHandler := handler.NewTeamRoleHandler(suite.teamSvc, roleSvc)
	auditHandler := handler.NewAuditHandler(suite.auditSvc)
	venueGalleryHandler := handler.NewVenueGalleryHandler(venueGallerySvc)
	areaGalleryHandler := handler.NewAreaGalleryHandler(areaGallerySvc)

	// Initialize Fiber app
	suite.app = fiber.New()

	// Setup routes
	routeConfig := route.RouteConfig{
		App:                 suite.app,
		AuthHandler:         authHandler,
		VenueHandler:        venueHandler,
		AreaHandler:         areaHandler,
		GraphHandler:        graphHandler,
		MediaHandler:        nil, // Skip for now
		TeamRoleHandler:     teamRoleHandler,
		AuditHandler:        auditHandler,
		VenueGalleryHandler: venueGalleryHandler,
		AreaGalleryHandler:  areaGalleryHandler,
	}
	routeConfig.Setup()

	// Seed roles and permissions first
	suite.seedRolesAndPermissions()

	// Create test user and organization for authenticated tests
	suite.createTestUserAndOrg()
}

func (suite *ComprehensiveHTTPTestSuite) seedRolesAndPermissions() {
	// Always seed for each test since each test has fresh database
	// Create permissions
	permissions := []entity.Permission{
		{Key: "venue:create", Description: "Create Venue", Group: "Venue"},
		{Key: "venue:update", Description: "Update Venue", Group: "Venue"},
		{Key: "venue:delete", Description: "Delete Venue", Group: "Venue"},
		{Key: "graph:edit", Description: "Edit Graph", Group: "Graph"},
		{Key: "graph:publish", Description: "Publish Graph", Group: "Graph"},
		{Key: "org:settings", Description: "Organization Settings", Group: "Organization"},
		{Key: "org:billing", Description: "Organization Billing", Group: "Organization"},
		{Key: "team:invite", Description: "Invite Team Members", Group: "Team"},
		{Key: "team:manage", Description: "Manage Team", Group: "Team"},
		{Key: "media:upload", Description: "Upload Media", Group: "Media"},
	}

	for _, perm := range permissions {
		perm.ID = uuid.New()
		err := suite.db.Create(&perm).Error
		suite.NoError(err, "Failed to create permission: %s", perm.Key)
	}

	// Create roles
	roles := []entity.Role{
		{Name: "Owner"},
		{Name: "Editor"},
		{Name: "Viewer"},
	}

	for _, role := range roles {
		role.ID = uuid.New()
		err := suite.db.Create(&role).Error
		suite.NoError(err, "Failed to create role: %s", role.Name)
	}

	// Associate permissions with Owner role
	var ownerRole entity.Role
	err := suite.db.First(&ownerRole, "name = ?", "Owner").Error
	suite.NoError(err, "Failed to find Owner role")

	var allPerms []entity.Permission
	err = suite.db.Find(&allPerms).Error
	suite.NoError(err, "Failed to find permissions")

	for _, perm := range allPerms {
		rolePerm := entity.RolePermission{
			RoleID:       ownerRole.ID,
			PermissionID: perm.ID,
		}
		err = suite.db.Create(&rolePerm).Error
		suite.NoError(err, "Failed to create role permission association")
	}

	// Verify seeding worked
	var count int64
	suite.db.Model(&entity.Permission{}).Count(&count)
	suite.Greater(count, int64(0), "No permissions seeded")

	suite.db.Model(&entity.Role{}).Count(&count)
	suite.Greater(count, int64(0), "No roles seeded")
}

func (suite *ComprehensiveHTTPTestSuite) createTestUserAndOrg() {
	// Create test organization with unique slug
	orgSlug := "test-org-" + uuid.New().String()
	testOrg := &entity.Organization{
		BaseEntity: entity.BaseEntity{ID: uuid.New()},
		Name:       "Test Organization",
		Slug:       orgSlug,
		IsActive:   true,
	}
	err := suite.db.Create(testOrg).Error
	suite.NoError(err, "Failed to create test organization")
	suite.testOrg = testOrg

	// Create test user
	userID := uuid.New()
	hashedPassword, err := utils.HashPassword("password123")
	suite.NoError(err, "Failed to hash password")

	testUser := &entity.User{
		BaseEntity:   entity.BaseEntity{ID: userID},
		FullName:     "Test User",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	}
	err = suite.db.Create(testUser).Error
	suite.NoError(err, "Failed to create test user")
	suite.testUser = testUser

	// Create organization membership
	ownerRole := &entity.Role{}
	err = suite.db.First(ownerRole, "name = ?", "Owner").Error
	suite.NoError(err, "Failed to find Owner role")

	membership := &entity.OrganizationMember{
		UserID:         userID,
		OrganizationID: testOrg.ID,
		RoleID:         ownerRole.ID,
	}
	err = suite.db.Create(membership).Error
	suite.NoError(err, "Failed to create organization membership")

	// Generate JWT token
	token, err := utils.GenerateToken(userID, "test@example.com", testOrg.ID, "Owner", []string{"venue:create", "venue:update", "venue:delete"})
	suite.NoError(err, "Failed to generate JWT token")
	suite.authToken = token
}

func (suite *ComprehensiveHTTPTestSuite) TearDownTest() {
	// Clean up test data
	if suite.db != nil {
		suite.db.Exec("DELETE FROM audit_logs")
		suite.db.Exec("DELETE FROM venue_gallery_items")
		suite.db.Exec("DELETE FROM area_gallery_items")
		suite.db.Exec("DELETE FROM graph_edges")
		suite.db.Exec("DELETE FROM graph_nodes")
		suite.db.Exec("DELETE FROM floors")
		suite.db.Exec("DELETE FROM graph_revisions")
		suite.db.Exec("DELETE FROM areas")
		suite.db.Exec("DELETE FROM venues")
		suite.db.Exec("DELETE FROM media_assets")
		suite.db.Exec("DELETE FROM organization_members")
		suite.db.Exec("DELETE FROM user_invitations")
		suite.db.Exec("DELETE FROM api_keys")
		suite.db.Exec("DELETE FROM users")
		suite.db.Exec("DELETE FROM organizations")
		suite.db.Exec("DELETE FROM role_permissions")
		suite.db.Exec("DELETE FROM permissions")
		suite.db.Exec("DELETE FROM roles")
	}
}

// ============================================================================
// AUTH ENDPOINTS TESTS
// ============================================================================

func (suite *ComprehensiveHTTPTestSuite) TestAuthRegister() {
	req := models.RegisterRequest{
		FullName:         "New User",
		Email:            "newuser@test.com",
		Password:         "password123",
		OrganizationName: "New Org",
	}
	reqBody, _ := json.Marshal(req)

	httpReq := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := suite.app.Test(httpReq)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode)

	var response models.AuthResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.NotEmpty(suite.T(), response.AccessToken)
	assert.Equal(suite.T(), "New User", response.User.FullName)
}

func (suite *ComprehensiveHTTPTestSuite) TestAuthLogin() {
	// First register a user
	suite.TestAuthRegister()

	req := models.LoginRequest{
		Email:    "newuser@test.com",
		Password: "password123",
	}
	reqBody, _ := json.Marshal(req)

	httpReq := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := suite.app.Test(httpReq)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode)

	var response models.AuthResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.NotEmpty(suite.T(), response.AccessToken)
}

// ============================================================================
// PUBLIC ENDPOINTS TESTS
// ============================================================================

func (suite *ComprehensiveHTTPTestSuite) TestGetVenueManifest() {
	httpReq := httptest.NewRequest("GET", "/api/v1/venues/non-existent/manifest", nil)
	resp, err := suite.app.Test(httpReq)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 404, resp.StatusCode)
}

func (suite *ComprehensiveHTTPTestSuite) TestGetAreaDetail() {
	httpReq := httptest.NewRequest("GET", "/api/v1/areas/"+uuid.New().String(), nil)
	resp, err := suite.app.Test(httpReq)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 404, resp.StatusCode)
}

// ============================================================================
// PROTECTED ENDPOINTS TESTS
// ============================================================================

func (suite *ComprehensiveHTTPTestSuite) TestListSystemRoles() {
	httpReq := httptest.NewRequest("GET", "/api/v1/roles", nil)
	httpReq.Header.Set("Authorization", "Bearer "+suite.authToken)
	resp, err := suite.app.Test(httpReq)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode)

	var apiResp struct {
		Success bool                `json:"success"`
		Data    []models.RoleDetail `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	suite.NoError(err)
	suite.True(apiResp.Success)
	assert.Greater(suite.T(), len(apiResp.Data), 0)
}

func (suite *ComprehensiveHTTPTestSuite) TestListPermissions() {
	httpReq := httptest.NewRequest("GET", "/api/v1/permissions", nil)
	httpReq.Header.Set("Authorization", "Bearer "+suite.authToken)
	resp, err := suite.app.Test(httpReq)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode)

	var apiResp struct {
		Success bool                    `json:"success"`
		Data    []models.PermissionNode `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	suite.NoError(err)
	suite.True(apiResp.Success)
	assert.Greater(suite.T(), len(apiResp.Data), 0)
}

// ============================================================================
// VENUE ENDPOINTS TESTS
// ============================================================================

func (suite *ComprehensiveHTTPTestSuite) TestCreateVenue() {
	req := models.CreateVenueRequest{
		Name: "Test Venue",
		Slug: "test-venue",
	}
	reqBody, _ := json.Marshal(req)

	httpReq := httptest.NewRequest("POST", "/api/v1/venues", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+suite.authToken)
	resp, err := suite.app.Test(httpReq)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 201, resp.StatusCode)

	var response models.IDResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.NotEqual(suite.T(), uuid.Nil, response.ID)
}

// ============================================================================
// TEAM MANAGEMENT ENDPOINTS TESTS
// ============================================================================

func (suite *ComprehensiveHTTPTestSuite) TestListMembers() {
	orgID := suite.testOrg.ID.String()
	httpReq := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/orgs/%s/members", orgID), nil)
	httpReq.Header.Set("Authorization", "Bearer "+suite.authToken)
	resp, err := suite.app.Test(httpReq)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode)

	var apiResp struct {
		Success bool                      `json:"success"`
		Data    []models.TeamMemberDetail `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	suite.NoError(err)
	suite.True(apiResp.Success)
	assert.Greater(suite.T(), len(apiResp.Data), 0)
}

// ============================================================================
// GRAPH EDITOR ENDPOINTS TESTS
// ============================================================================

func (suite *ComprehensiveHTTPTestSuite) TestGetEditorData() {
	// Create a venue first
	req := models.CreateVenueRequest{
		Name: "Editor Test Venue",
		Slug: "editor-test-venue-" + uuid.New().String()[:8],
	}
	reqBody, _ := json.Marshal(req)

	venueHttpReq := httptest.NewRequest("POST", "/api/v1/venues", bytes.NewReader(reqBody))
	venueHttpReq.Header.Set("Content-Type", "application/json")
	venueHttpReq.Header.Set("Authorization", "Bearer "+suite.authToken)
	venueResp, err := suite.app.Test(venueHttpReq)
	suite.NoError(err)
	suite.Equal(201, venueResp.StatusCode)

	var apiResp struct {
		Success bool              `json:"success"`
		Data    models.IDResponse `json:"data"`
	}
	err = json.NewDecoder(venueResp.Body).Decode(&apiResp)
	suite.NoError(err)
	suite.True(apiResp.Success)
	suite.NotNil(apiResp.Data.ID)

	venueID, ok := apiResp.Data.ID.(string)
	suite.True(ok, "ID should be a string")

	// Now test editor data
	httpReq := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/editor/%s", venueID), nil)
	httpReq.Header.Set("Authorization", "Bearer "+suite.authToken)
	resp, err := suite.app.Test(httpReq)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode)
}

// ============================================================================
// ERROR HANDLING TESTS
// ============================================================================

func (suite *ComprehensiveHTTPTestSuite) TestUnauthorizedAccess() {
	httpReq := httptest.NewRequest("GET", "/api/v1/roles", nil)
	// No Authorization header
	resp, err := suite.app.Test(httpReq)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 401, resp.StatusCode)
}

func (suite *ComprehensiveHTTPTestSuite) TestInvalidJSON() {
	httpReq := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := suite.app.Test(httpReq)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 400, resp.StatusCode)
}

func TestComprehensiveHTTPTestSuite(t *testing.T) {
	suite.Run(t, new(ComprehensiveHTTPTestSuite))
}
