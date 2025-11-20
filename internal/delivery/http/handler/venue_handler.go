package handler

import (
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VenueHandler struct {
	service service.VenueService
}

func NewVenueHandler(s service.VenueService) *VenueHandler {
	return &VenueHandler{service: s}
}

// POST /api/v1/venues (Admin Create)
func (h *VenueHandler) CreateVenue(c *fiber.Ctx) error {
	var req models.CreateVenueRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	resp, err := h.service.CreateVenue(c.Context(), req)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendCreated(c, resp)
}

// GET /api/v1/venues/:slug/manifest (Mobile App Read)
func (h *VenueHandler) GetManifest(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return utils.SendError(c, 400, "Slug is required")
	}

	manifest, err := h.service.GetMobileManifest(c.Context(), slug)
	if err != nil {
		return utils.SendError(c, 404, "Venue not found or not published")
	}

	// Khusus manifest, return raw struct agar strukturnya sesuai persis dengan DTO
	return c.JSON(manifest)
}

// GET /api/v1/venues/:id (Admin Detail)
func (h *VenueHandler) GetDetail(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.SendError(c, 400, "Invalid UUID")
	}

	// TODO: Implement GetVenueDetail di service dulu (ini placeholder)
	// detail, err := h.service.GetVenueDetail(c.Context(), id)

	return utils.SendSuccess(c, fiber.Map{"id": id, "status": "todo"})
}
