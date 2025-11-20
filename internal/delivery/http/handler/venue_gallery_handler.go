package handler

import (
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VenueGalleryHandler struct {
	service service.VenueGalleryService
}

func NewVenueGalleryHandler(s service.VenueGalleryService) *VenueGalleryHandler {
	return &VenueGalleryHandler{service: s}
}

func (h *VenueGalleryHandler) AddItems(c *fiber.Ctx) error {
	var req models.AddGalleryVenueItemsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.AddGalleryItems(c.Context(), req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Items added to venue gallery")
}

func (h *VenueGalleryHandler) Reorder(c *fiber.Ctx) error {
	var req models.ReorderVenueGalleryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.ReorderGallery(c.Context(), req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Venue gallery reordered")
}

func (h *VenueGalleryHandler) UpdateItem(c *fiber.Ctx) error {
	var req models.UpdateVenueGalleryItemRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.UpdateGalleryItem(c.Context(), req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Item updated")
}

func (h *VenueGalleryHandler) RemoveItem(c *fiber.Ctx) error {
	venueID, err := uuid.Parse(c.Params("venue_id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Venue ID")
	}

	mediaID, err := uuid.Parse(c.Params("media_id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Media ID")
	}

	if err := h.service.RemoveGalleryItem(c.Context(), venueID, mediaID); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Item removed")
}
