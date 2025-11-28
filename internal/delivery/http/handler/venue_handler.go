package handler

import (
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/utils"
	"strconv"

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

	// Get organization ID from context (set by Protected middleware)
	orgID, ok := c.Locals("org_id").(uuid.UUID)
	if !ok {
		return utils.SendError(c, 400, "Organization ID not found in context")
	}

	resp, err := h.service.CreateVenue(c.Context(), orgID, req)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendCreated(c, resp)
}

// GET /api/v1/venues/:orgSlug/:venueSlug/manifest (Mobile App Read)
func (h *VenueHandler) GetManifest(c *fiber.Ctx) error {
	orgSlug := c.Params("orgSlug")
	venueSlug := c.Params("venueSlug")

	if orgSlug == "" || venueSlug == "" {
		return utils.SendError(c, 400, "Org slug and venue slug are required")
	}

	manifest, err := h.service.GetMobileManifest(c.Context(), orgSlug, venueSlug)
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

	detail, err := h.service.GetVenueDetail(c.Context(), id)
	if err != nil {
		return utils.SendError(c, 404, "Venue not found")
	}

	return utils.SendSuccess(c, detail)
}

// PUT /api/v1/venues/:id (Admin Update)
func (h *VenueHandler) UpdateVenue(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.SendError(c, 400, "Invalid UUID")
	}

	var req models.UpdateVenueRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.UpdateVenue(c.Context(), id, req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, "Venue updated successfully")
}

// GET /api/v1/venues (Admin List)
func (h *VenueHandler) ListVenues(c *fiber.Ctx) error {
	// Default values
	defaultLimit := 20
	defaultOffset := 0

	query := models.VenueQuery{
		Limit:  &defaultLimit,
		Offset: &defaultOffset,
	}

	// Get organization ID from context (set by Protected middleware)
	orgID, ok := c.Locals("org_id").(uuid.UUID)
	if ok {
		query.OrganizationID = &orgID
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			query.Limit = &limit
		}
	}

	// Parse offset
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			query.Offset = &offset
		}
	}

	// Parse filters
	if name := c.Query("name"); name != "" {
		query.Name = &name
	}
	if city := c.Query("city"); city != "" {
		query.City = &city
	}
	if visibility := c.Query("visibility"); visibility != "" {
		query.Visibility = &visibility
	}

	venues, total, err := h.service.ListVenues(c.Context(), query)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, fiber.Map{
		"venues": venues,
		"total":  total,
		"limit":  *query.Limit,
		"offset": *query.Offset,
	})
}
