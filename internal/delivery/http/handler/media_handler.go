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
