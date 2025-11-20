package handler

import (
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type VenueGalleryHandler struct {
	service service.VenueGalleryService
}

func NewVenueGalleryHandler(v service.VenueGalleryService) *VenueGalleryHandler {
	return &VenueGalleryHandler{service: v}
}

func (h *VenueGalleryHandler) AddVenueItems(c *fiber.Ctx) error {
	var req models.AddGalleryVenueItemsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}
	if err := h.service.AddGalleryItems(c.Context(), req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, nil)
}

func (h *VenueGalleryHandler) ReorderVenue(c *fiber.Ctx) error {
	var req models.ReorderVenueGalleryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}
	if err := h.service.ReorderGallery(c.Context(), req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, nil)
}
