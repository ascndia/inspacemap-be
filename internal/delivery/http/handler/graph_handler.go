package handler

import (
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GraphHandler struct {
	service service.GraphService
}

func NewGraphHandler(s service.GraphService) *GraphHandler {
	return &GraphHandler{service: s}
}

// --- EDITOR DATA ---

// GET /api/v1/editor/:venue_id
func (h *GraphHandler) GetEditorData(c *fiber.Ctx) error {
	venueID, err := uuid.Parse(c.Params("venue_id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Venue ID")
	}

	data, err := h.service.GetEditorData(c.Context(), venueID)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, data)
}

// POST /api/v1/editor/floors
func (h *GraphHandler) CreateFloor(c *fiber.Ctx) error {
	var req models.CreateFloorRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	// Ambil VenueID dari query param atau context (disini asumsi ada di query ?venue_id=...)
	// Atau DTO CreateFloorRequest bisa diupdate untuk menerima VenueID
	venueIDStr := c.Query("venue_id")
	venueID, err := uuid.Parse(venueIDStr)
	if err != nil {
		return utils.SendError(c, 400, "venue_id query param required")
	}

	resp, err := h.service.CreateFloor(c.Context(), venueID, req)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendCreated(c, resp)
}

// --- NODES ---

// POST /api/v1/editor/nodes
func (h *GraphHandler) CreateNode(c *fiber.Ctx) error {
	var req models.CreateNodeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	resp, err := h.service.CreateNode(c.Context(), req)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendCreated(c, resp)
}

// PUT /api/v1/editor/nodes/:id/position
func (h *GraphHandler) UpdateNodePosition(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	var req models.UpdateNodePositionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.UpdateNodePosition(c.Context(), id, req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, nil)
}

// PUT /api/v1/editor/nodes/:id/calibration
func (h *GraphHandler) CalibrateNode(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	var req models.UpdateNodeCalibrationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.UpdateNodeCalibration(c.Context(), id, req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, nil)
}

// --- EDGES ---

// POST /api/v1/editor/connections
func (h *GraphHandler) ConnectNodes(c *fiber.Ctx) error {
	var req models.ConnectNodesRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.ConnectNodes(c.Context(), req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, nil)
}

// --- PUBLISH ---

// POST /api/v1/editor/:venue_id/publish
func (h *GraphHandler) Publish(c *fiber.Ctx) error {
	venueID, _ := uuid.Parse(c.Params("venue_id"))
	var req models.PublishDraftRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.PublishChanges(c.Context(), venueID, req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Graph Published Successfully")
}
