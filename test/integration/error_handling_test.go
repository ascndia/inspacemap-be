package integration_test

import (
	"context"
	"testing"

	"inspacemap/backend/internal/entity"
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/repository"

	"github.com/google/uuid"
)

// =================================================================
// ERROR HANDLING & EDGE CASES TESTS
// =================================================================

// TestRegisterDuplicateEmail: Verifikasi error handling untuk email duplikat
func TestRegisterDuplicateEmail(t *testing.T) {
	ctx := context.Background()

	// 1. ARRANGE: Seed role Owner
	ownerRole := entity.Role{
		BaseEntity: entity.BaseEntity{ID: uuid.New()},
		Name:       "Owner",
	}
	if err := testDB.Create(&ownerRole).Error; err != nil {
		t.Fatalf("Failed to seed owner role: %v", err)
	}

	// 2. ACT: Register pertama (harus berhasil)
	req := models.RegisterRequest{
		FullName:         "First User",
		Email:            "duplicate@test.com",
		Password:         "password123",
		OrganizationName: "First Org",
	}
	_, err1 := authSvc.Register(ctx, req)
	if err1 != nil {
		t.Fatalf("First register should succeed: %v", err1)
	}

	// 3. ACT: Register kedua dengan email sama (harus gagal)
	req2 := models.RegisterRequest{
		FullName:         "Second User",
		Email:            "duplicate@test.com", // Same email
		Password:         "password456",
		OrganizationName: "Second Org",
	}
	_, err2 := authSvc.Register(ctx, req2)

	// 4. ASSERT: Harus ada error
	if err2 == nil {
		t.Error("Expected error for duplicate email, but got none")
	} else {
		t.Logf("✅ Correctly rejected duplicate email: %v", err2)
	}
}

// TestConnectNonExistentNodes: Verifikasi error handling untuk node yang tidak ada
func TestConnectNonExistentNodes(t *testing.T) {
	ctx := context.Background()

	// 1. ARRANGE: Buat venue dan floor kosong (tanpa nodes)
	orgID := uuid.New()
	venueID := uuid.New()

	testDB.Create(&entity.Organization{BaseEntity: entity.BaseEntity{ID: orgID}, Name: "TestErrorOrg"})
	venue := entity.Venue{
		BaseEntity:     entity.BaseEntity{ID: venueID},
		OrganizationID: orgID,
		Name:           "Test Error Venue",
	}
	if err := testDB.Create(&venue).Error; err != nil {
		t.Fatalf("Failed to create test venue: %v", err)
	}

	floorReq := models.CreateFloorRequest{
		Name:           "Test Floor",
		LevelIndex:     1,
		MapWidth:       1000,
		MapHeight:      500,
		PixelsPerMeter: 10.0,
	}
	floorResp, err := graphSvc.CreateFloor(ctx, venueID, floorReq)
	if err != nil {
		t.Fatalf("Failed to create floor: %v", err)
	}
	floorID := floorResp.ID.(uuid.UUID)

	// 2. ACT: Coba connect node yang tidak ada
	fakeNodeID1 := uuid.New()
	fakeNodeID2 := uuid.New()
	connReq := models.ConnectNodesRequest{
		FromNodeID: fakeNodeID1,
		ToNodeID:   fakeNodeID2,
	}
	err = graphSvc.ConnectNodes(ctx, connReq)

	// 3. ASSERT: Harus ada error
	if err == nil {
		t.Error("Expected error when connecting non-existent nodes, but got none")
	} else {
		t.Logf("✅ Correctly rejected connection of non-existent nodes: %v", err)
	}

	// 4. ACT: Coba connect node yang ada ke node yang tidak ada
	// Buat satu node dulu
	nodeReq := models.CreateNodeRequest{
		FloorID:         floorID,
		X:               100,
		Y:               100,
		PanoramaAssetID: uuid.New(),
		Label:           "Test Node",
	}
	nodeResp, err := graphSvc.CreateNode(ctx, nodeReq)
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
	realNodeID := nodeResp.ID.(uuid.UUID)

	// Connect real node ke fake node
	connReq2 := models.ConnectNodesRequest{
		FromNodeID: realNodeID,
		ToNodeID:   fakeNodeID2,
	}
	err2 := graphSvc.ConnectNodes(ctx, connReq2)

	// 5. ASSERT: Harus ada error
	if err2 == nil {
		t.Error("Expected error when connecting real node to non-existent node, but got none")
	} else {
		t.Logf("✅ Correctly rejected connection to non-existent node: %v", err2)
	}
}

// TestPublishEmptyDraft: Verifikasi error handling untuk publish draft kosong
func TestPublishEmptyDraft(t *testing.T) {
	ctx := context.Background()

	// 1. ARRANGE: Buat venue tanpa floor/node apapun
	orgID := uuid.New()
	venueID := uuid.New()

	testDB.Create(&entity.Organization{BaseEntity: entity.BaseEntity{ID: orgID}, Name: "TestEmptyOrg"})
	venue := entity.Venue{
		BaseEntity:     entity.BaseEntity{ID: venueID},
		OrganizationID: orgID,
		Name:           "Test Empty Venue",
	}
	if err := testDB.Create(&venue).Error; err != nil {
		t.Fatalf("Failed to create test venue: %v", err)
	}

	// 2. ACT: Coba publish draft kosong (tanpa floor/node)
	publishReq := models.PublishDraftRequest{Note: "Empty draft test"}
	err := graphSvc.PublishChanges(ctx, venueID, publishReq)

	// 3. ASSERT: Harus ada error atau warning
	// (Bergantung implementasi - bisa error atau berhasil dengan draft kosong)
	if err != nil {
		t.Logf("✅ Publish empty draft correctly failed: %v", err)
	} else {
		// Jika berhasil, verifikasi bahwa live revision kosong
		var updatedVenue entity.Venue
		testDB.First(&updatedVenue, venueID)

		if updatedVenue.LiveRevisionID != uuid.Nil {
			// Cek apakah live revision benar-benar kosong
			var liveRev entity.GraphRevision
			if err := testDB.Preload("Floors.Nodes").First(&liveRev, updatedVenue.LiveRevisionID).Error; err != nil {
				t.Errorf("Live revision not found: %v", err)
			} else {
				totalNodes := 0
				for _, floor := range liveRev.Floors {
					totalNodes += len(floor.Nodes)
				}
				if totalNodes != 0 {
					t.Errorf("Expected 0 nodes in empty draft publish, got %d", totalNodes)
				} else {
					t.Log("✅ Publish empty draft succeeded with 0 nodes (acceptable)")
				}
			}
		} else {
			t.Error("Publish empty draft did not update LiveRevisionID")
		}
	}
}

// TestCreateNodeInvalidCoordinates: Edge case untuk koordinat negatif
func TestCreateNodeInvalidCoordinates(t *testing.T) {
	ctx := context.Background()

	// 1. ARRANGE: Buat venue dan floor
	orgID := uuid.New()
	venueID := uuid.New()

	testDB.Create(&entity.Organization{BaseEntity: entity.BaseEntity{ID: orgID}, Name: "TestCoordOrg"})
	venue := entity.Venue{
		BaseEntity:     entity.BaseEntity{ID: venueID},
		OrganizationID: orgID,
		Name:           "Test Coord Venue",
	}
	if err := testDB.Create(&venue).Error; err != nil {
		t.Fatalf("Failed to create test venue: %v", err)
	}

	floorReq := models.CreateFloorRequest{
		Name:           "Test Floor",
		LevelIndex:     1,
		MapWidth:       1000,
		MapHeight:      500,
		PixelsPerMeter: 10.0,
	}
	floorResp, err := graphSvc.CreateFloor(ctx, venueID, floorReq)
	if err != nil {
		t.Fatalf("Failed to create floor: %v", err)
	}
	floorID := floorResp.ID.(uuid.UUID)

	// 2. ACT: Coba buat node dengan koordinat negatif
	nodeReq := models.CreateNodeRequest{
		FloorID:         floorID,
		X:               -100, // Invalid: negative X
		Y:               -50,  // Invalid: negative Y
		PanoramaAssetID: uuid.New(),
		Label:           "Invalid Node",
	}
	_, err = graphSvc.CreateNode(ctx, nodeReq)

	// 3. ASSERT: Harus ada error
	if err == nil {
		t.Error("Expected error for negative coordinates, but got none")
	} else {
		t.Logf("✅ Correctly rejected negative coordinates: %v", err)
	}
}

// TestRolePermissionsValidation: Verifikasi role permissions bekerja
func TestRolePermissionsValidation(t *testing.T) {
	ctx := context.Background()

	// 1. ARRANGE: Buat role dengan permissions
	roleID := uuid.New()
	testRole := entity.Role{
		BaseEntity: entity.BaseEntity{ID: roleID},
		Name:       "TestRole",
	}
	if err := testDB.Create(&testRole).Error; err != nil {
		t.Fatalf("Failed to create test role: %v", err)
	}

	// Buat permissions
	perm1 := entity.Permission{
		BaseEntity:  entity.BaseEntity{ID: uuid.New()},
		Key:         "venue.create",
		Description: "Create Venue",
		Group:       "Venue",
	}
	perm2 := entity.Permission{
		BaseEntity:  entity.BaseEntity{ID: uuid.New()},
		Key:         "venue.delete",
		Description: "Delete Venue",
		Group:       "Venue",
	}
	if err := testDB.Create(&perm1).Error; err != nil {
		t.Fatalf("Failed to create permission 1: %v", err)
	}
	if err := testDB.Create(&perm2).Error; err != nil {
		t.Fatalf("Failed to create permission 2: %v", err)
	}

	// Associate permissions dengan role
	rolePerm1 := entity.RolePermission{RoleID: roleID, PermissionID: perm1.ID}
	rolePerm2 := entity.RolePermission{RoleID: roleID, PermissionID: perm2.ID}
	if err := testDB.Create(&rolePerm1).Error; err != nil {
		t.Fatalf("Failed to associate permission 1: %v", err)
	}
	if err := testDB.Create(&rolePerm2).Error; err != nil {
		t.Fatalf("Failed to associate permission 2: %v", err)
	}

	// 2. ACT: Ambil permissions untuk role ini
	roleRepo := repository.NewRoleRepository(testDB)
	permissions, err := roleRepo.GetPermissions(ctx, roleID)
	if err != nil {
		t.Fatalf("Failed to get role permissions: %v", err)
	}

	// 3. ASSERT: Harus ada 2 permissions
	if len(permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(permissions))
	} else {
		// Verifikasi permission keys
		permKeys := make(map[string]bool)
		for _, p := range permissions {
			permKeys[p.Key] = true
		}
		if !permKeys["venue.create"] || !permKeys["venue.delete"] {
			t.Errorf("Expected permissions 'venue.create' and 'venue.delete', got %v", permKeys)
		} else {
			t.Log("✅ Role permissions correctly associated and retrieved")
		}
	}
}

// TestGraphRevisionHistory: Verifikasi multiple publishes menyimpan history
func TestGraphRevisionHistory(t *testing.T) {
	ctx := context.Background()

	// 1. ARRANGE: Buat venue
	orgID := uuid.New()
	venueID := uuid.New()

	testDB.Create(&entity.Organization{BaseEntity: entity.BaseEntity{ID: orgID}, Name: "TestHistoryOrg"})
	venue := entity.Venue{
		BaseEntity:     entity.BaseEntity{ID: venueID},
		OrganizationID: orgID,
		Name:           "Test History Venue",
	}
	if err := testDB.Create(&venue).Error; err != nil {
		t.Fatalf("Failed to create test venue: %v", err)
	}

	// 2. ACT: Publish pertama (empty draft)
	publishReq1 := models.PublishDraftRequest{Note: "First publish - empty"}
	err1 := graphSvc.PublishChanges(ctx, venueID, publishReq1)
	if err1 != nil {
		t.Logf("First publish failed (expected for empty draft): %v", err1)
	}

	// Buat floor dan node
	floorReq := models.CreateFloorRequest{
		Name:           "Ground Floor",
		LevelIndex:     1,
		MapWidth:       1000,
		MapHeight:      500,
		PixelsPerMeter: 10.0,
	}
	floorResp, err := graphSvc.CreateFloor(ctx, venueID, floorReq)
	if err != nil {
		t.Fatalf("Failed to create floor: %v", err)
	}
	floorID := floorResp.ID.(uuid.UUID)

	nodeReq := models.CreateNodeRequest{
		FloorID:         floorID,
		X:               100,
		Y:               100,
		PanoramaAssetID: uuid.New(),
		Label:           "Entry Node",
	}
	_, err = graphSvc.CreateNode(ctx, nodeReq)
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}

	// Publish kedua
	publishReq2 := models.PublishDraftRequest{Note: "Second publish - with content"}
	err2 := graphSvc.PublishChanges(ctx, venueID, publishReq2)
	if err2 != nil {
		t.Fatalf("Second publish should succeed: %v", err2)
	}

	// 3. ASSERT: Cek history revisions
	revisionRepo := repository.NewGraphRevisionRepository(testDB)
	revisions, err := revisionRepo.GetByVenueID(ctx, venueID)
	if err != nil {
		t.Fatalf("Failed to get revision history: %v", err)
	}

	// Harus ada minimal 1 revision (live), mungkin lebih jika ada draft
	if len(revisions) == 0 {
		t.Error("Expected at least 1 revision in history")
	} else {
		t.Logf("✅ Found %d revisions in history", len(revisions))

		// Verifikasi venue menunjuk ke live revision terbaru
		var updatedVenue entity.Venue
		testDB.First(&updatedVenue, venueID)

		if updatedVenue.LiveRevisionID == uuid.Nil {
			t.Error("Venue should have live revision after publish")
		} else {
			t.Log("✅ Venue correctly points to latest live revision")
		}
	}
}

// TestVenueNonExistentOperations: Edge case untuk operasi pada venue yang tidak ada
func TestVenueNonExistentOperations(t *testing.T) {
	ctx := context.Background()

	// 1. ARRANGE: Fake venue ID
	fakeVenueID := uuid.New()

	// 2. ACT: Coba berbagai operasi pada venue yang tidak ada
	floorReq := models.CreateFloorRequest{
		Name:           "Fake Floor",
		LevelIndex:     1,
		MapWidth:       1000,
		MapHeight:      500,
		PixelsPerMeter: 10.0,
	}

	_, err1 := graphSvc.CreateFloor(ctx, fakeVenueID, floorReq)

	publishReq := models.PublishDraftRequest{Note: "Fake publish"}
	err2 := graphSvc.PublishChanges(ctx, fakeVenueID, publishReq)

	// 3. ASSERT: CreateFloor mungkin berhasil (karena lazy init draft),
	// tapi PublishChanges harus error
	if err1 != nil {
		t.Logf("✅ CreateFloor correctly failed for non-existent venue: %v", err1)
	} else {
		t.Log("ℹ️  CreateFloor succeeded (lazy draft creation) - this may be acceptable")
	}

	if err2 == nil {
		t.Error("Expected error when publishing non-existent venue")
	} else {
		t.Logf("✅ Correctly rejected publish: %v", err2)
	}
}

// TestSelfConnectionPrevention: Verifikasi pencegahan self-connection
func TestSelfConnectionPrevention(t *testing.T) {
	ctx := context.Background()

	// 1. ARRANGE: Buat venue, floor, dan node
	orgID := uuid.New()
	venueID := uuid.New()

	testDB.Create(&entity.Organization{BaseEntity: entity.BaseEntity{ID: orgID}, Name: "TestSelfOrg"})
	venue := entity.Venue{
		BaseEntity:     entity.BaseEntity{ID: venueID},
		OrganizationID: orgID,
		Name:           "Test Self Venue",
	}
	if err := testDB.Create(&venue).Error; err != nil {
		t.Fatalf("Failed to create test venue: %v", err)
	}

	floorReq := models.CreateFloorRequest{
		Name:           "Test Floor",
		LevelIndex:     1,
		MapWidth:       1000,
		MapHeight:      500,
		PixelsPerMeter: 10.0,
	}
	floorResp, err := graphSvc.CreateFloor(ctx, venueID, floorReq)
	if err != nil {
		t.Fatalf("Failed to create floor: %v", err)
	}
	floorID := floorResp.ID.(uuid.UUID)

	nodeReq := models.CreateNodeRequest{
		FloorID:         floorID,
		X:               100,
		Y:               100,
		PanoramaAssetID: uuid.New(),
		Label:           "Test Node",
	}
	nodeResp, err := graphSvc.CreateNode(ctx, nodeReq)
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
	nodeID := nodeResp.ID.(uuid.UUID)

	// 2. ACT: Coba connect node ke dirinya sendiri
	connReq := models.ConnectNodesRequest{
		FromNodeID: nodeID,
		ToNodeID:   nodeID, // Same node
	}
	err = graphSvc.ConnectNodes(ctx, connReq)

	// 3. ASSERT: Harus ada error
	if err == nil {
		t.Error("Expected error when connecting node to itself, but got none")
	} else {
		t.Logf("✅ Correctly prevented self-connection: %v", err)
	}
}
