package main

import (
	"log"

	"inspacemap/backend/config"
	"inspacemap/backend/internal/entity"

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
}
