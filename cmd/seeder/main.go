package main

import (
	"log"
	"os"

	"inspacemap/backend/config"
	"inspacemap/backend/internal/entity"
	"inspacemap/backend/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {

	config.ConnectDB()
	db := config.DB

	log.Println("🌱 Seeding Database...")

	permissions := []entity.Permission{

		{Key: "venue:create", Description: "Create new venues", Group: "CMS"},
		{Key: "venue:update", Description: "Update venue details", Group: "CMS"},
		{Key: "venue:delete", Description: "Delete venues", Group: "CMS"},

		{Key: "graph:edit", Description: "Edit nodes and edges", Group: "Graph"},
		{Key: "graph:publish", Description: "Publish draft to live", Group: "Graph"},

		{Key: "org:settings", Description: "Manage organization profile", Group: "Org"},
		{Key: "org:billing", Description: "Manage billing and subscription", Group: "Org"},
		{Key: "team:invite", Description: "Invite new members", Group: "Team"},
		{Key: "team:manage", Description: "Change member roles", Group: "Team"},

		{Key: "media:upload", Description: "Upload new assets", Group: "Media"},
	}

	for _, p := range permissions {

		var count int64
		db.Model(&entity.Permission{}).Where("key = ?", p.Key).Count(&count)
		if count == 0 {
			p.ID = uuid.New()
			db.Create(&p)
			log.Printf("Created Permission: %s", p.Key)
		}
	}

	roles := []struct {
		Name        string
		Description string
		PermKeys    []string
	}{
		{
			Name:        "Owner",
			Description: "Pemilik Organisasi (Super Admin)",
			PermKeys:    []string{"venue:create", "venue:update", "venue:delete", "graph:edit", "graph:publish", "org:settings", "org:billing", "team:invite", "team:manage", "media:upload"},
		},
		{
			Name:        "Editor",
			Description: "Pengelola Konten (Peta & Info)",
			PermKeys:    []string{"venue:create", "venue:update", "graph:edit", "graph:publish", "media:upload"},
		},
		{
			Name:        "Viewer",
			Description: "Hanya bisa melihat (Read Only)",
			PermKeys:    []string{},
		},
	}

	for _, r := range roles {
		var role entity.Role

		if err := db.Where("name = ?", r.Name).First(&role).Error; err != nil {
			if err == gorm.ErrRecordNotFound {

				role = entity.Role{

					Name:        r.Name,
					Description: r.Description,
				}

				role.ID = uuid.New()
				db.Create(&role)
				log.Printf("Created Role: %s", r.Name)
			} else {
				log.Fatalf("Error checking role: %v", err)
			}
		}

		var perms []entity.Permission
		if len(r.PermKeys) > 0 {
			db.Where("key IN ?", r.PermKeys).Find(&perms)

			for _, p := range perms {
				pivot := entity.RolePermission{
					RoleID:       role.ID,
					PermissionID: p.ID,
				}
				db.Clauses(clause.OnConflict{DoNothing: true}).Create(&pivot)
			}
		}
	}

	log.Println("✅ Seeding Completed!")

	// Additional development seeding
	if os.Getenv("APP_ENV") == "development" {
		log.Println("🌱 Seeding Development Data...")
		seedDevelopmentData(db)
		log.Println("✅ Development Seeding Completed!")
	}
}

func seedDevelopmentData(db *gorm.DB) {
	// Get roles
	var ownerRole, editorRole, viewerRole entity.Role
	db.Where("name = ?", "Owner").First(&ownerRole)
	db.Where("name = ?", "Editor").First(&editorRole)
	db.Where("name = ?", "Viewer").First(&viewerRole)

	// Sample users
	users := []struct {
		email    string
		password string
		fullName string
		role     entity.UserRole
	}{
		{"admin@inspacemap.dev", "admin123", "Admin User", entity.RoleOwner},
		{"editor@inspacemap.dev", "editor123", "Editor User", entity.RoleEditor},
		{"viewer@inspacemap.dev", "viewer123", "Viewer User", entity.RoleViewer},
	}

	for _, u := range users {
		var count int64
		db.Model(&entity.User{}).Where("email = ?", u.email).Count(&count)
		if count == 0 {
			hash, err := utils.HashPassword(u.password)
			if err != nil {
				log.Fatalf("Error hashing password: %v", err)
			}
			user := entity.User{
				Email:        u.email,
				PasswordHash: hash,
				FullName:     u.fullName,
				IsEmailVerified: true,
			}
			user.ID = uuid.New()
			db.Create(&user)
			log.Printf("Created User: %s", u.email)
		}
	}

	// Sample organization
	var orgCount int64
	db.Model(&entity.Organization{}).Where("slug = ?", "demo-org").Count(&orgCount)
	if orgCount == 0 {
		org := entity.Organization{
			Name:     "Demo Organization",
			Slug:     "demo-org",
			Website:  "https://demo.inspacemap.dev",
			IsActive: true,
		}
		org.ID = uuid.New()
		db.Create(&org)
		log.Printf("Created Organization: %s", org.Name)

		// Add users to organization
		var adminUser, editorUser, viewerUser entity.User
		db.Where("email = ?", "admin@inspacemap.dev").First(&adminUser)
		db.Where("email = ?", "editor@inspacemap.dev").First(&editorUser)
		db.Where("email = ?", "viewer@inspacemap.dev").First(&viewerUser)

		members := []struct {
			user entity.User
			role entity.Role
		}{
			{adminUser, ownerRole},
			{editorUser, editorRole},
			{viewerUser, viewerRole},
		}

		for _, m := range members {
			member := entity.OrganizationMember{
				OrganizationID: org.ID,
				UserID:         m.user.ID,
				RoleID:         m.role.ID,
			}
			member.ID = uuid.New()
			db.Clauses(clause.OnConflict{DoNothing: true}).Create(&member)
		}
		log.Printf("Added members to organization: %s", org.Name)
	}

	// Sample venue
	var venueCount int64
	db.Model(&entity.Venue{}).Where("slug = ?", "demo-venue").Count(&venueCount)
	if venueCount == 0 {
		var org entity.Organization
		db.Where("slug = ?", "demo-org").First(&org)

		var adminUser entity.User
		db.Where("email = ?", "admin@inspacemap.dev").First(&adminUser)

		// Create a draft revision first
		draftRevision := entity.GraphRevision{
			OrganizationID: org.ID,
			CreatedByID:    adminUser.ID,
			VenueID:        uuid.New(), // Will be set after venue creation
			Status:         entity.StatusDraft,
			Note:           "Auto-generated initial draft",
		}
		draftRevision.ID = uuid.New()
		// Don't create yet, will create after venue

		venue := entity.Venue{
			OrganizationID:  org.ID,
			Name:           "Demo Venue",
			Slug:           "demo-venue",
			Description:    "A sample venue for development testing",
			Address:        "123 Demo Street",
			City:           "Demo City",
			Province:       "Demo Province",
			PostalCode:     "12345",
			Latitude:       -6.2088, // Jakarta coordinates
			Longitude:      106.8456,
		}
		venue.ID = uuid.New()
		draftRevision.VenueID = venue.ID
		db.Create(&draftRevision)
		venue.LiveRevisionID = draftRevision.ID
		venue.DraftRevisionID = &draftRevision.ID
		db.Create(&venue)
		log.Printf("Created Venue: %s", venue.Name)
	}
}
