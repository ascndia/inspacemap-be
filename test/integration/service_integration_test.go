package integration_test

import (
	"context"
	"log"
	"os"
	"testing"

	"inspacemap/backend/config"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"
	"inspacemap/backend/internal/service"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var authSvc service.AuthService
var graphSvc service.GraphService
var userRepo repository.UserRepository
var orgRepo repository.OrganizationRepository

func TestMain(m *testing.M) {

	log.Println("🔄 Starting database connection and migration...")
	config.ConnectDB()
	testDB = config.DB
	log.Println("✅ Database connected and migrated successfully")

	log.Println("🔄 Initializing repositories...")
	userRepo = repository.NewUserRepository(testDB)
	orgRepo = repository.NewOrganizationRepository(testDB)
	log.Println("✅ Repositories initialized")

	log.Println("🔄 Initializing services...")
	authSvc = service.NewAuthService(
		userRepo,
		orgRepo,
		repository.NewOrganizationMemberRepository(testDB),
		repository.NewInvitationRepository(testDB),
		repository.NewRoleRepository(testDB),
	)
	log.Println("✅ Auth service initialized")

	graphSvc = service.NewGraphService(
		repository.NewGraphRepository(testDB), repository.NewGraphRevisionRepository(testDB),
		repository.NewFloorRepository(testDB), repository.NewVenueRepository(testDB),
	)
	log.Println("✅ Graph service initialized")

	log.Println("🚀 Starting tests...")

	code := m.Run()

	os.Exit(code)
}

func TestAuthAndMembershipFlow(t *testing.T) {
	ctx := context.Background()
	email := "owner@test.com"

	ownerRole := entity.Role{
		BaseEntity: entity.BaseEntity{ID: uuid.New()},
		Name:       "Owner",
	}

	if err := testDB.Create(&ownerRole).Error; err != nil {
		t.Fatalf("Failed to seed owner role: %v", err)
	}
	t.Logf("Owner role created with ID: %s", ownerRole.ID)

	roleRepo := repository.NewRoleRepository(testDB)
	foundRole, err := roleRepo.GetByName(ctx, "Owner")
	if err != nil {
		t.Fatalf("Failed to find owner role: %v", err)
	}
	t.Logf("Owner role found: %s", foundRole.Name)

	req := models.RegisterRequest{
		FullName:         "Test Owner",
		Email:            email,
		Password:         "password",
		OrganizationName: "Test Org Inc",
	}
	authResp, err := authSvc.Register(ctx, req)

	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if authResp.User.Email != email {
		t.Errorf("Expected email %s, got %s", email, authResp.User.Email)
	}
	if len(authResp.User.Organizations) != 1 {
		t.Fatalf("Expected 1 organization in membership, got %d", len(authResp.User.Organizations))
	}

	orgID := authResp.User.Organizations[0].OrganizationID
	t.Logf("Organization ID from response: %s", orgID)

	var org entity.Organization
	if err := testDB.First(&org, "id = ?", orgID).Error; err != nil {
		t.Fatalf("Org was not created in DB: %v", err)
	}
	t.Logf("Org '%s' created successfully.", org.Name)
}

func TestGraphPublishIntegrity(t *testing.T) {
	ctx := context.Background()

	orgID := uuid.New()
	venueID := uuid.New()

	testDB.Create(&entity.Organization{BaseEntity: entity.BaseEntity{ID: orgID}, Name: "TestPublishOrg"})
	venue := entity.Venue{
		BaseEntity:     entity.BaseEntity{ID: venueID},
		OrganizationID: orgID,
		Name:           "Test Graph Venue",
	}
	if err := testDB.Create(&venue).Error; err != nil {
		t.Fatalf("Failed to create test venue: %v", err)
	}

	floorReq := models.CreateFloorRequest{
		Name:           "Ground Floor",
		LevelIndex:     1,
		MapWidth:       1000,
		MapHeight:      500,
		PixelsPerMeter: 10.0,
	}

	floorResp, err := graphSvc.CreateFloor(ctx, venueID, floorReq)
	if err != nil {
		t.Fatalf("Failed to create floor/draft: %v", err)
	}
	floorID := floorResp.ID.(uuid.UUID)

	nodeReq := models.CreateNodeRequest{
		FloorID: floorID, X: 100, Y: 100,
		PanoramaAssetID: uuid.New(),
		Label:           "Lobby",
	}
	nodeResp, err := graphSvc.CreateNode(ctx, nodeReq)
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
	nodeID1 := nodeResp.ID.(uuid.UUID)

	nodeReq2 := models.CreateNodeRequest{
		FloorID: floorID, X: 200, Y: 200,
		PanoramaAssetID: uuid.New(),
		Label:           "Hall",
	}
	nodeResp2, err := graphSvc.CreateNode(ctx, nodeReq2)
	if err != nil {
		t.Fatalf("Failed to create second node: %v", err)
	}
	nodeID2 := nodeResp2.ID.(uuid.UUID)

	connReq := models.ConnectNodesRequest{FromNodeID: nodeID1, ToNodeID: nodeID2}
	if err := graphSvc.ConnectNodes(ctx, connReq); err != nil {
		t.Fatalf("Failed to connect nodes: %v", err)
	}

	testDB.First(&venue, venueID)
	t.Logf("Venue DraftRevisionID before publish: %s", venue.DraftRevisionID)

	publishReq := models.PublishDraftRequest{Note: "Initial stable release"}
	if err := graphSvc.PublishChanges(ctx, venueID, publishReq); err != nil {
		t.Fatalf("PublishChanges failed (TRANSACTION ROLLBACK): %v", err)
	}

	var updatedVenue entity.Venue
	testDB.First(&updatedVenue, venueID)
	t.Logf("DraftRevisionID after publish: %s", updatedVenue.DraftRevisionID)

	if updatedVenue.LiveRevisionID == uuid.Nil {
		t.Fatalf("LiveRevisionID was NOT updated after publish.")
	}
	liveRevID := updatedVenue.LiveRevisionID

	var liveNodesCount int64
	testDB.Table("graph_nodes").
		Joins("JOIN floors ON floors.id = graph_nodes.floor_id").
		Joins("JOIN graph_revisions gr ON gr.id = floors.graph_revision_id").
		Where("gr.id = ?", liveRevID).
		Count(&liveNodesCount)

	if liveNodesCount != 2 {
		t.Errorf("Expected 2 nodes in LIVE revision, got %d", liveNodesCount)
	}

	var draft entity.GraphRevision
	if err := testDB.First(&draft, "id = ?", venue.DraftRevisionID).Error; err != nil {
		t.Errorf("Original Draft record missing: %v", err)
	}

	t.Log("✅ Graph integrity test passed. Deep Copy successful.")
}
