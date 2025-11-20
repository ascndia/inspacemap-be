package main

import (
	"log"
	"os"

	"inspacemap/backend/config"
	"inspacemap/backend/internal/delivery/http/handler"
	"inspacemap/backend/internal/delivery/http/route"
	"inspacemap/backend/internal/repository"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 1. INIT CONFIG & DB
	// Pastikan .env atau environment variables sudah diset
	config.ConnectDB()
	db := config.DB

	// 2. INIT REPOSITORIES (Data Access Layer)
	userRepo := repository.NewUserRepository(db)
	orgRepo := repository.NewOrganizationRepository(db)
	orgMemberRepo := repository.NewOrganizationMemberRepository(db)
	invitationRepo := repository.NewInvitationRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permRepo := repository.NewPermissionRepository(db)

	venueRepo := repository.NewVenueRepository(db)
	floorRepo := repository.NewFloorRepository(db)
	areaRepo := repository.NewAreaRepository(db)

	graphRepo := repository.NewGraphRepository(db)
	revisionRepo := repository.NewGraphRevisionRepository(db)

	mediaRepo := repository.NewMediaRepository(db)
	venueGalleryRepo := repository.NewVenueGalleryRepository(db)
	areaGalleryRepo := repository.NewAreaGalleryRepository(db)

	auditRepo := repository.NewAuditRepository(db)

	// 3. INIT EXTERNAL SERVICES (Infrastructure)
	// Mengambil config dari Env
	minioEndpoint := getEnv("MINIO_ENDPOINT", "http://localhost:9000")
	minioAccess := getEnv("MINIO_ACCESS_KEY", "admin_inspacemap")
	minioSecret := getEnv("MINIO_SECRET_KEY", "password_rahasia_banget")
	minioRegion := getEnv("MINIO_REGION", "us-east-1")
	minioBucket := getEnv("MINIO_BUCKET", "panoramas")
	cdnURL := getEnv("CDN_BASE_URL", "http://localhost:9000/panoramas")

	storageProvider := storage.NewMinIOProvider(minioEndpoint, minioAccess, minioSecret, minioRegion)

	// 4. INIT SERVICES (Business Logic Layer)
	authService := service.NewAuthService(userRepo, orgRepo, orgMemberRepo, invitationRepo, roleRepo)
	mediaService := service.NewMediaService(mediaRepo, storageProvider, minioBucket, cdnURL)
	areaService := service.NewAreaService(areaRepo, areaGalleryRepo, graphRepo)
	graphService := service.NewGraphService(graphRepo, revisionRepo, floorRepo, venueRepo)
	venueService := service.NewVenueService(venueRepo)
	teamService := service.NewTeamService(userRepo, invitationRepo, orgMemberRepo, roleRepo)
	roleService := service.NewRoleService(roleRepo, permRepo)
	venueGalleryService := service.NewVenueGalleryService(venueGalleryRepo)
	areaGalleryService := service.NewAreaGalleryService(areaGalleryRepo)
	auditService := service.NewAuditService(auditRepo)

	// 5. INIT HANDLERS (HTTP Transport Layer)
	authHandler := handler.NewAuthHandler(authService)
	teamRoleHandler := handler.NewTeamRoleHandler(teamService, roleService)
	venueHandler := handler.NewVenueHandler(venueService)
	venueGalleryHandler := handler.NewVenueGalleryHandler(venueGalleryService) // Implementasi nanti
	graphHandler := handler.NewGraphHandler(graphService)
	mediaHandler := handler.NewMediaHandler(mediaService)
	areaHandler := handler.NewAreaHandler(areaService)
	areaGalleryHandler := handler.NewAreaGalleryHandler(areaGalleryService) // Implementasi nanti
	auditHandler := handler.NewAuditHandler(auditService)                   // Implementasi nanti
	// 6. SETUP FIBER APP
	app := fiber.New(fiber.Config{
		AppName: "InSpaceMap API v1",
	})

	// Middleware Global
	app.Use(logger.New())  // Logging request
	app.Use(recover.New()) // Mencegah crash jika panic
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Untuk development. Ubah domain spesifik saat prod.
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Tenant-ID",
	}))

	// 7. REGISTER ROUTES
	routeConfig := route.RouteConfig{
		App:                 app,
		TeamRoleHandler:     teamRoleHandler,
		AuthHandler:         authHandler,
		AreaHandler:         areaHandler,
		AreaGalleryHandler:  areaGalleryHandler,
		VenueHandler:        venueHandler,
		VenueGalleryHandler: venueGalleryHandler,
		GraphHandler:        graphHandler,
		MediaHandler:        mediaHandler,
		AuditHandler:        auditHandler,
	}
	routeConfig.Setup()

	// 8. START SERVER
	port := getEnv("PORT", "8080")
	log.Printf("🚀 Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

// Helper kecil untuk baca env di main
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
