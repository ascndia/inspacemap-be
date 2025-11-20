package handler

import (
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TeamHandler struct {
	teamService service.TeamService
	roleService service.RoleService
}

func NewTeamHandler(ts service.TeamService, rs service.RoleService) *TeamHandler {
	return &TeamHandler{teamService: ts, roleService: rs}
}

// --- TEAM ---

// POST /api/v1/orgs/:org_id/invite
func (h *TeamHandler) InviteMember(c *fiber.Ctx) error {
	orgID, _ := uuid.Parse(c.Params("org_id"))
	userID := getUserID(c) // Pengundang

	var req models.InviteUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.teamService.InviteMember(c.Context(), orgID, userID, req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Invitation sent")
}

// GET /api/v1/orgs/:org_id/members
func (h *TeamHandler) ListMembers(c *fiber.Ctx) error {
	orgID, _ := uuid.Parse(c.Params("org_id"))
	list, err := h.teamService.GetMembersList(c.Context(), orgID)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, list)
}

// PATCH /api/v1/orgs/:org_id/members
func (h *TeamHandler) UpdateMemberRole(c *fiber.Ctx) error {
	orgID, _ := uuid.Parse(c.Params("org_id"))
	var req models.UpdateUserRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}
	if err := h.teamService.UpdateMemberRole(c.Context(), orgID, req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Role updated")
}

// DELETE /api/v1/orgs/:org_id/members/:user_id
func (h *TeamHandler) RemoveMember(c *fiber.Ctx) error {
	orgID, _ := uuid.Parse(c.Params("org_id"))
	targetUserID, _ := uuid.Parse(c.Params("user_id"))

	if err := h.teamService.RemoveMember(c.Context(), orgID, targetUserID); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Member removed")
}

// --- ROLES ---

// GET /api/v1/orgs/:org_id/roles
func (h *TeamHandler) ListRoles(c *fiber.Ctx) error {
	orgID, _ := uuid.Parse(c.Params("org_id"))
	roles, err := h.roleService.GetRolesByOrganization(c.Context(), orgID)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, roles)
}

// GET /api/v1/permissions (List semua izin yang tersedia untuk bikin role)
func (h *TeamHandler) ListPermissions(c *fiber.Ctx) error {
	perms, err := h.roleService.GetAvailablePermissions(c.Context())
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, perms)
}

// POST /api/v1/orgs/:org_id/roles (Custom Role)
func (h *TeamHandler) CreateRole(c *fiber.Ctx) error {
	orgID, _ := uuid.Parse(c.Params("org_id"))
	var req models.CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}
	resp, err := h.roleService.CreateRole(c.Context(), orgID, req)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendCreated(c, resp)
}
