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

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "inspacemap"),
		getEnv("DB_PORT", "5432"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),

		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("🚀 Database Connected Successfully")
	log.Println("Running Auto Migration...")

	log.Println("Creating base tables...")
	err = DB.AutoMigrate(
		&entity.User{},
		&entity.Permission{},
		&entity.Role{},
		&entity.Organization{},
	)
	if err != nil {
		log.Fatal("Migration Failed at base tables: ", err)
	}
	log.Println("✅ Base tables created")

	log.Println("Creating asset tables...")
	err = DB.AutoMigrate(
		&entity.MediaAsset{},
	)
	if err != nil {
		log.Fatal("Migration Failed at asset tables: ", err)
	}
	log.Println("✅ Asset tables created")

	log.Println("Creating venue & graph tables...")
	err = DB.AutoMigrate(
		&entity.Venue{},
		&entity.GraphRevision{},
		&entity.Floor{},
		&entity.GraphNode{},
		&entity.GraphEdge{},
	)
	if err != nil {
		log.Fatal("Migration Failed at venue & graph tables: ", err)
	}
	log.Println("✅ Venue & graph tables created")

	log.Println("Creating relation tables...")
	err = DB.AutoMigrate(
		&entity.OrganizationMember{},
		&entity.UserInvitation{},
		&entity.ApiKey{},
		&entity.RolePermission{},
	)
	if err != nil {
		log.Fatal("Migration Failed at relation tables: ", err)
	}
	log.Println("✅ Relation tables created")

	log.Println("Creating feature tables...")
	err = DB.AutoMigrate(
		&entity.Area{},
		&entity.VenueGalleryItem{},
		&entity.AreaGalleryItem{},
		&entity.AuditLog{},
	)
	if err != nil {
		log.Fatal("Migration Failed at feature tables: ", err)
	}
	log.Println("✅ Feature tables created")

	log.Println("✅ Database Migration Completed")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
