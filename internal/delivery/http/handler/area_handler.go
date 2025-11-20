package handler

import (
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AreaHandler struct {
	service service.AreaService
}

func NewAreaHandler(s service.AreaService) *AreaHandler {
	return &AreaHandler{service: s}
}

func (h *AreaHandler) CreateArea(c *fiber.Ctx) error {
	var req models.CreateAreaRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}
	resp, err := h.service.CreateArea(c.Context(), req)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendCreated(c, resp)
}

func (h *AreaHandler) GetDetail(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	resp, err := h.service.GetAreaDetail(c.Context(), id)
	if err != nil {
		return utils.SendError(c, 404, "Area not found")
	}
	return utils.SendSuccess(c, resp)
}

func (h *AreaHandler) GetVenueAreas(c *fiber.Ctx) error {
	venueID, _ := uuid.Parse(c.Params("venue_id"))
	resp, err := h.service.GetVenueAreas(c.Context(), venueID)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, resp)
}
