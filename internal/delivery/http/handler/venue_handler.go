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

	detail, err := h.service.GetVenueDetail(c.Context(), id)
	if err != nil {
		return utils.SendError(c, 404, "Venue not found")
	}

	return utils.SendSuccess(c, detail)
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
