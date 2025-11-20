package handler

import (
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AuditHandler struct {
	service service.AuditService
}

func NewAuditHandler(s service.AuditService) *AuditHandler {
	return &AuditHandler{service: s}
}

// GET /api/v1/audit-logs
func (h *AuditHandler) GetLogs(c *fiber.Ctx) error {
	orgID := getOrgID(c) // Dari middleware

	// Parse Query Params ke DTO
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	cursor := c.Query("cursor")

	query := models.AuditLogQueryCursor{
		AuditLogFilter: models.AuditLogFilter{
			OrganizationID: orgID.String(),
			Action:         c.Query("action"),
			UserID:         c.Query("user_id"),
		},
		Limit:  &limit,
		Cursor: &cursor,
	}

	resp, err := h.service.GetActivityLogs(c.Context(), orgID, query)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, resp)
}
