package handler

import (
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MediaHandler struct {
	service service.MediaService
}

func NewMediaHandler(s service.MediaService) *MediaHandler {
	return &MediaHandler{service: s}
}

// POST /api/v1/media/upload-init
func (h *MediaHandler) InitUpload(c *fiber.Ctx) error {
	var req models.PresignedUploadRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	// Ambil Org ID dari Token
	orgID := getOrgID(c)
	if orgID == uuid.Nil {
		// Fallback untuk testing atau superadmin, tapi sebaiknya wajib
		return utils.SendError(c, 401, "Organization Context Required")
	}

	resp, err := h.service.InitDirectUpload(c.Context(), orgID, req)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, resp)
}

// POST /api/v1/media/confirm
func (h *MediaHandler) ConfirmUpload(c *fiber.Ctx) error {
	var req models.ConfirmUploadRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.ConfirmUpload(c.Context(), req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, "Upload Confirmed")
}

// GET /api/v1/media
func (h *MediaHandler) ListAssets(c *fiber.Ctx) error {
	// Parse query parameters
	var query models.MediaAssetQuery
	if err := c.QueryParser(&query); err != nil {
		return utils.SendError(c, 400, "Invalid query parameters")
	}

	// Set organization ID from context
	orgID := getOrgID(c)
	if orgID != uuid.Nil {
		query.OrganizationID = &orgID
	}

	assets, total, err := h.service.ListAssets(c.Context(), query)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, fiber.Map{
		"assets": assets,
		"total":  total,
	})
}

// GET /api/v1/media/:id
func (h *MediaHandler) GetAsset(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.SendError(c, 400, "Invalid asset ID")
	}

	asset, err := h.service.GetAsset(c.Context(), id)
	if err != nil {
		return utils.SendError(c, 404, "Asset not found")
	}

	return utils.SendSuccess(c, asset)
}

// DELETE /api/v1/media/:id
func (h *MediaHandler) DeleteAsset(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return utils.SendError(c, 400, "Invalid asset ID")
	}

	if err := h.service.DeleteAsset(c.Context(), id); err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, "Asset deleted")
}
