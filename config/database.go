package config

import (
	"fmt"
	"log"
	"os"

	"inspacemap/backend/internal/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	// Ambil konfigurasi dari Environment Variables
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"), // Default password lokal
		getEnv("DB_NAME", "inspacemap"),
		getEnv("DB_PORT", "5432"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log query SQL (untuk dev)
	})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("🚀 Database Connected Successfully")

	// Auto Migrate: Membuat tabel otomatis
	log.Println("Running Auto Migration...")
	err = DB.AutoMigrate(
		&entity.Organization{},
		&entity.OrganizationMember{},
		&entity.User{},
		&entity.UserInvitation{},
		&entity.ApiKey{},

		&entity.Venue{},
		&entity.VenueGalleryItem{},
		&entity.Area{},
		&entity.AreaGalleryItem{},
		&entity.MediaAsset{},

		&entity.GraphRevision{},
		&entity.Floor{},
		&entity.GraphNode{},
		&entity.GraphEdge{},

		&entity.Role{},
		&entity.Permission{},
		&entity.RolePermission{},

		&entity.AuditLog{},
	)

	if err != nil {
		log.Fatal("Migration Failed: ", err)
	}
	log.Println("✅ Database Migration Completed")
}

// Helper untuk baca ENV dengan fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
