package handler

import (
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AreaGalleryHandler struct {
	service service.AreaGalleryService
}

func NewAreaGalleryHandler(s service.AreaGalleryService) *AreaGalleryHandler {
	return &AreaGalleryHandler{service: s}
}

// POST /gallery/area
func (h *AreaGalleryHandler) AddItems(c *fiber.Ctx) error {
	var req models.AddAreaGalleryItemsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.AddGalleryItems(c.Context(), req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Items added to area gallery")
}

// PUT /gallery/area/reorder
func (h *AreaGalleryHandler) Reorder(c *fiber.Ctx) error {
	var req models.ReorderAreaGalleryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.ReorderGallery(c.Context(), req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Area gallery reordered")
}

// PATCH /gallery/area/item
func (h *AreaGalleryHandler) UpdateItem(c *fiber.Ctx) error {
	var req models.UpdateAreaGalleryItemRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.UpdateGalleryItem(c.Context(), req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Item updated")
}

// DELETE /gallery/area/:area_id/:media_id
func (h *AreaGalleryHandler) RemoveItem(c *fiber.Ctx) error {
	areaID, _ := uuid.Parse(c.Params("area_id"))
	mediaID, _ := uuid.Parse(c.Params("media_id"))

	if err := h.service.RemoveGalleryItem(c.Context(), areaID, mediaID); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Item removed")
}
