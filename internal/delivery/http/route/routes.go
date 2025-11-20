package route

import (
	"inspacemap/backend/internal/delivery/http/handler"
	"inspacemap/backend/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App          *fiber.App
	AuthHandler  *handler.AuthHandler
	VenueHandler *handler.VenueHandler
	GraphHandler *handler.GraphHandler
	MediaHandler *handler.MediaHandler
}

func (c *RouteConfig) Setup() {
	api := c.App.Group("/api/v1")

	// --- PUBLIC ROUTES ---
	api.Post("/auth/login", c.AuthHandler.Login)
	api.Post("/auth/register", c.AuthHandler.Register)
	api.Post("/auth/invite/accept", c.AuthHandler.AcceptInvite)

	// Mobile App Fetch Manifest (Public atau API Key Protected)
	api.Get("/venues/:slug/manifest", c.VenueHandler.GetManifest)

	// --- PROTECTED ROUTES (Butuh Login) ---
	protected := api.Group("/", middleware.Protected())

	// --- ORGANIZATION SCOPED ---
	tenant := protected.Group("/", middleware.TenantGuard())

	// Media
	tenant.Post("/media/upload-init", c.MediaHandler.InitUpload)
	tenant.Post("/media/confirm", c.MediaHandler.ConfirmUpload)

	// Venue Management
	tenant.Post("/venues", c.VenueHandler.CreateVenue)

	// Graph Editor
	editor := tenant.Group("/editor")
	editor.Get("/:venue_id", c.GraphHandler.GetEditorData)
	editor.Post("/floors", c.GraphHandler.CreateFloor)
	editor.Post("/nodes", c.GraphHandler.CreateNode)
	editor.Put("/nodes/:id/position", c.GraphHandler.UpdateNodePosition)
	editor.Put("/nodes/:id/calibration", c.GraphHandler.CalibrateNode)
	editor.Post("/connections", c.GraphHandler.ConnectNodes)
	editor.Post("/:venue_id/publish", c.GraphHandler.Publish)
}
