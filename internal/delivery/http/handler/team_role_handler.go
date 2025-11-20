package handler

import (
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TeamRoleHandler struct {
	teamService service.TeamService
	roleService service.RoleService
}

func NewTeamRoleHandler(ts service.TeamService, rs service.RoleService) *TeamRoleHandler {
	return &TeamRoleHandler{
		teamService: ts,
		roleService: rs,
	}
}

func (h *TeamRoleHandler) ListSystemRoles(c *fiber.Ctx) error {
	roles, err := h.roleService.GetSystemRoles(c.Context())
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, roles)
}

func (h *TeamRoleHandler) ListPermissions(c *fiber.Ctx) error {
	perms, err := h.roleService.GetAvailablePermissions(c.Context())
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, perms)
}

func (h *TeamRoleHandler) ListMembers(c *fiber.Ctx) error {
	orgID, err := uuid.Parse(c.Params("org_id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Organization ID")
	}

	list, err := h.teamService.GetMembersList(c.Context(), orgID)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, list)
}

func (h *TeamRoleHandler) InviteMember(c *fiber.Ctx) error {
	orgID, err := uuid.Parse(c.Params("org_id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Organization ID")
	}

	inviterID := getUserID(c)

	var req models.InviteUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.teamService.InviteMember(c.Context(), orgID, inviterID, req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, "Invitation sent successfully")
}

func (h *TeamRoleHandler) UpdateMemberRole(c *fiber.Ctx) error {
	orgID, err := uuid.Parse(c.Params("org_id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Organization ID")
	}

	var req models.UpdateUserRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.teamService.UpdateMemberRole(c.Context(), orgID, req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, "Member role updated")
}

func (h *TeamRoleHandler) RemoveMember(c *fiber.Ctx) error {
	orgID, err := uuid.Parse(c.Params("org_id"))
	targetUserID, err2 := uuid.Parse(c.Params("user_id"))

	if err != nil || err2 != nil {
		return utils.SendError(c, 400, "Invalid ID format")
	}

	if err := h.teamService.RemoveMember(c.Context(), orgID, targetUserID); err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, "Member removed from organization")
}
