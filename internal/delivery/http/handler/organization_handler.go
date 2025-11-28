package handler

import (
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrganizationHandler struct {
	service service.OrganizationService
}

func NewOrganizationHandler(s service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{service: s}
}

// GET /api/v1/orgs/:id (Get Organization Detail)
func (h *OrganizationHandler) GetDetail(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.SendError(c, 400, "Invalid UUID")
	}

	detail, err := h.service.GetDetailByID(c.Context(), id)
	if err != nil {
		return utils.SendError(c, 404, "Organization not found")
	}

	return utils.SendSuccess(c, detail)
}

// PUT /api/v1/orgs/:id (Update Organization Profile)
func (h *OrganizationHandler) UpdateProfile(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.SendError(c, 400, "Invalid UUID")
	}

	var req models.UpdateOrganizationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.UpdateProfile(c.Context(), id, req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, "Organization updated successfully")
}

// GET /api/v1/orgs (List Organizations - Admin Only)
func (h *OrganizationHandler) ListOrganizations(c *fiber.Ctx) error {
	// Default values
	defaultLimit := 20
	defaultOffset := 0

	query := models.OrganizationQuery{
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
	if slug := c.Query("slug"); slug != "" {
		query.Slug = &slug
	}

	orgs, total, err := h.service.ListOrganizations(c.Context(), query)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, fiber.Map{
		"organizations": orgs,
		"total":         total,
		"limit":         *query.Limit,
		"offset":        *query.Offset,
	})
}
