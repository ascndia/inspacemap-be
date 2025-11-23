package handler

import (
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/utils"
	"strings"

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

// PUT /api/v1/editor/floors/:id
func (h *GraphHandler) UpdateFloor(c *fiber.Ctx) error {
	floorID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Floor ID")
	}

	var req models.UpdateFloorRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.UpdateFloor(c.Context(), floorID, req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, "Floor updated successfully")
}

// DELETE /api/v1/editor/floors/:id
func (h *GraphHandler) DeleteFloor(c *fiber.Ctx) error {
	floorID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Floor ID")
	}

	if err := h.service.DeleteFloor(c.Context(), floorID); err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, "Floor deleted successfully")
}

// GET /api/v1/editor/floors/:id
func (h *GraphHandler) GetFloor(c *fiber.Ctx) error {
	floorID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Floor ID")
	}

	floor, err := h.service.GetFloor(c.Context(), floorID)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, floor)
}

// GET /api/v1/editor/:venue_id/floors
func (h *GraphHandler) GetFloors(c *fiber.Ctx) error {
	venueID, err := uuid.Parse(c.Params("venue_id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Venue ID")
	}

	floors, err := h.service.GetFloors(c.Context(), venueID)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, floors)
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

// PUT /api/v1/editor/nodes/:id
func (h *GraphHandler) UpdateNode(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Node ID")
	}

	var req models.UpdateNodeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.UpdateNode(c.Context(), id, req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, "Node updated successfully")
}

// DELETE /api/v1/editor/nodes/:id
func (h *GraphHandler) DeleteNode(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))

	if err := h.service.DeleteNode(c.Context(), id); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Node deleted successfully")
}

// --- EDGES ---

// POST /api/v1/editor/connections
func (h *GraphHandler) ConnectNodes(c *fiber.Ctx) error {
	var req models.ConnectNodesRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	resp, err := h.service.ConnectNodes(c.Context(), req)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendCreated(c, resp)
}

// DELETE /api/v1/editor/connections
func (h *GraphHandler) DeleteConnection(c *fiber.Ctx) error {
	fromIDStr := c.Query("from_node_id")
	toIDStr := c.Query("to_node_id")

	fromID, err := uuid.Parse(fromIDStr)
	if err != nil {
		return utils.SendError(c, 400, "from_node_id query param required")
	}

	toID, err := uuid.Parse(toIDStr)
	if err != nil {
		return utils.SendError(c, 400, "to_node_id query param required")
	}

	if err := h.service.DeleteConnection(c.Context(), fromID, toID); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Connection deleted successfully")
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
		// Check for business logic validation errors
		if strings.Contains(err.Error(), "no draft revision found") ||
			strings.Contains(err.Error(), "cannot publish empty revision") {
			return utils.SendError(c, 400, err.Error())
		}
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Graph Published Successfully")
}

// --- GRAPH REVISION MANAGEMENT ---

// POST /api/v1/editor/:venue_id/revisions/draft
func (h *GraphHandler) CreateDraftRevision(c *fiber.Ctx) error {
	venueID, err := uuid.Parse(c.Params("venue_id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Venue ID")
	}

	resp, err := h.service.CreateDraftRevision(c.Context(), venueID)
	if err != nil {
		// Check for business logic validation errors
		if strings.Contains(err.Error(), "draft already exists") {
			return utils.SendError(c, 400, err.Error())
		}
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendCreated(c, resp)
}

// GET /api/v1/editor/:venue_id/revisions
func (h *GraphHandler) ListRevisions(c *fiber.Ctx) error {
	venueID, err := uuid.Parse(c.Params("venue_id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Venue ID")
	}

	revisions, err := h.service.ListRevisions(c.Context(), venueID)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, revisions)
}

// GET /api/v1/editor/revisions/:revision_id
func (h *GraphHandler) GetRevisionDetail(c *fiber.Ctx) error {
	revisionID, err := uuid.Parse(c.Params("revision_id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Revision ID")
	}

	detail, err := h.service.GetRevisionDetail(c.Context(), revisionID)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, detail)
}

// PUT /api/v1/editor/revisions/:revision_id
func (h *GraphHandler) UpdateRevision(c *fiber.Ctx) error {
	revisionID, err := uuid.Parse(c.Params("revision_id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Revision ID")
	}

	var req models.UpdateRevisionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.UpdateRevision(c.Context(), revisionID, req); err != nil {
		// Check for business logic validation errors
		if strings.Contains(err.Error(), "only draft revisions can be updated") {
			return utils.SendError(c, 400, err.Error())
		}
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, "Revision updated successfully")
}

// DELETE /api/v1/editor/revisions/:revision_id
func (h *GraphHandler) DeleteRevision(c *fiber.Ctx) error {
	revisionID, err := uuid.Parse(c.Params("revision_id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid Revision ID")
	}

	if err := h.service.DeleteRevision(c.Context(), revisionID); err != nil {
		// Check for business logic validation errors
		if strings.Contains(err.Error(), "only draft revisions can be deleted") {
			return utils.SendError(c, 400, err.Error())
		}
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, "Revision deleted successfully")
}
