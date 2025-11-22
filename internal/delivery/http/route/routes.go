package route

import (
	"inspacemap/backend/internal/delivery/http/handler"
	"inspacemap/backend/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                 *fiber.App
	TeamRoleHandler     *handler.TeamRoleHandler
	AuthHandler         *handler.AuthHandler
	VenueHandler        *handler.VenueHandler
	VenueGalleryHandler *handler.VenueGalleryHandler
	AreaHandler         *handler.AreaHandler
	AreaGalleryHandler  *handler.AreaGalleryHandler
	GraphHandler        *handler.GraphHandler
	MediaHandler        *handler.MediaHandler
	AuditHandler        *handler.AuditHandler
}

func (c *RouteConfig) Setup() {
	api := c.App.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/login", c.AuthHandler.Login)
	auth.Post("/register", c.AuthHandler.Register)
	auth.Post("/invite/accept", c.AuthHandler.AcceptInvite)

	api.Get("/venues/:slug/manifest", c.VenueHandler.GetManifest)
	api.Get("/areas/:id", c.AreaHandler.GetDetail)

	protected := api.Group("/", middleware.Protected())
	protected.Get("/roles", c.TeamRoleHandler.ListSystemRoles)
	protected.Get("/permissions", c.TeamRoleHandler.ListPermissions)

	tenant := protected.Group("/", middleware.TenantGuard())

	media := tenant.Group("/media")
	media.Post("/upload-init", c.MediaHandler.InitUpload)
	media.Post("/confirm", c.MediaHandler.ConfirmUpload)
	media.Get("/", c.MediaHandler.ListAssets)
	media.Get("/:id", c.MediaHandler.GetAsset)
	media.Delete("/:id", c.MediaHandler.DeleteAsset)

	venues := tenant.Group("/venues")
	venues.Post("/", c.VenueHandler.CreateVenue)
	venues.Get("/", c.VenueHandler.ListVenues)
	venues.Get("/:id", c.VenueHandler.GetDetail)

	// Venue Gallery endpoints
	venues.Post("/:id/gallery", c.VenueGalleryHandler.AddItems)
	venues.Put("/:id/gallery/reorder", c.VenueGalleryHandler.Reorder)
	venues.Patch("/:id/gallery/:media_id", c.VenueGalleryHandler.UpdateItem)
	venues.Delete("/:id/gallery/:media_id", c.VenueGalleryHandler.RemoveItem)

	areas := tenant.Group("/areas")
	areas.Post("/", c.AreaHandler.CreateArea)

	tenant.Get("/venues/:venue_id/areas", c.AreaHandler.GetVenueAreas)

	aGallery := tenant.Group("/gallery/area")
	aGallery.Post("/", c.AreaGalleryHandler.AddItems)
	aGallery.Put("/reorder", c.AreaGalleryHandler.Reorder)
	aGallery.Patch("/item", c.AreaGalleryHandler.UpdateItem)
	aGallery.Delete("/:area_id/:media_id", c.AreaGalleryHandler.RemoveItem)

	orgs := tenant.Group("/orgs/:org_id")
	orgs.Get("/members", c.TeamRoleHandler.ListMembers)
	orgs.Post("/invite", c.TeamRoleHandler.InviteMember)
	orgs.Patch("/members", c.TeamRoleHandler.UpdateMemberRole)
	orgs.Delete("/members/:user_id", c.TeamRoleHandler.RemoveMember)

	tenant.Get("/audit-logs", c.AuditHandler.GetLogs)

	editor := tenant.Group("/editor")

	editor.Get("/:venue_id", c.GraphHandler.GetEditorData)

	
	editor.Post("/floors", c.GraphHandler.CreateFloor)
	editor.Get("/:venue_id/floors", c.GraphHandler.GetFloors)
	editor.Get("/floors/:id", c.GraphHandler.GetFloor)
	editor.Put("/floors/:id", c.GraphHandler.UpdateFloor)
	editor.Delete("/floors/:id", c.GraphHandler.DeleteFloor)

	editor.Post("/nodes", c.GraphHandler.CreateNode)
	editor.Put("/nodes/:id/position", c.GraphHandler.UpdateNodePosition)
	editor.Put("/nodes/:id/calibration", c.GraphHandler.CalibrateNode)
	editor.Delete("/nodes/:id", c.GraphHandler.DeleteNode)

	editor.Post("/connections", c.GraphHandler.ConnectNodes)
	editor.Delete("/connections", c.GraphHandler.DeleteConnection)

	editor.Post("/:venue_id/publish", c.GraphHandler.Publish)

	// Graph Revision Management
	editor.Post("/:venue_id/revisions/draft", c.GraphHandler.CreateDraftRevision)
	editor.Get("/:venue_id/revisions", c.GraphHandler.ListRevisions)
	editor.Get("/revisions/:revision_id", c.GraphHandler.GetRevisionDetail)
	editor.Put("/revisions/:revision_id", c.GraphHandler.UpdateRevision)
	editor.Delete("/revisions/:revision_id", c.GraphHandler.DeleteRevision)
}
