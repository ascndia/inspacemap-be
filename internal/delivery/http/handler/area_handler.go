package handler

import (
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"
	"inspacemap/backend/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AreaHandler struct {
	service      service.AreaService
	venueService service.VenueService
}

func NewAreaHandler(s service.AreaService, vs service.VenueService) *AreaHandler {
	return &AreaHandler{
		service:      s,
		venueService: vs,
	}
}

func (h *AreaHandler) CreateArea(c *fiber.Ctx) error {
	var req models.CreateAreaRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	// Get floor_id from URL params
	floorID, err := uuid.Parse(c.Params("floor_id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid floor_id")
	}
	req.FloorID = &floorID

	resp, err := h.service.CreateArea(c.Context(), req)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendCreated(c, resp)
}

func (h *AreaHandler) SetStartNode(c *fiber.Ctx) error {
	areaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid area_id")
	}

	var req models.SetStartNodeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	resp, err := h.service.SetAreaStartNode(c.Context(), areaID, req)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, resp)
}

func (h *AreaHandler) UpdateArea(c *fiber.Ctx) error {
	areaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid area_id")
	}

	var req models.UpdateAreaRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, 400, "Invalid JSON")
	}

	if err := h.service.UpdateArea(c.Context(), areaID, req); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Area updated")
}

func (h *AreaHandler) DeleteArea(c *fiber.Ctx) error {
	areaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.SendError(c, 400, "Invalid area_id")
	}

	if err := h.service.DeleteArea(c.Context(), areaID); err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, "Area deleted")
}

func (h *AreaHandler) GetDetail(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	resp, err := h.service.GetAreaDetail(c.Context(), id)
	if err != nil {
		return utils.SendError(c, 404, "Area not found")
	}
	return utils.SendSuccess(c, resp)
}

func (h *AreaHandler) GetVenueAreas(c *fiber.Ctx) error {
	venueID, _ := uuid.Parse(c.Params("venue_id"))
	resp, err := h.service.GetVenueAreas(c.Context(), venueID)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}
	return utils.SendSuccess(c, resp)
}

func (h *AreaHandler) ListAreas(c *fiber.Ctx) error {
	// Default values
	defaultLimit := 20
	defaultOffset := 0

	query := models.AreaQuery{
		Limit:  &defaultLimit,
		Offset: &defaultOffset,
	}

	// Parse venue_id from URL (required for tenant filtering)
	var venueID uuid.UUID
	if venueIDStr := c.Params("venue_id"); venueIDStr != "" {
		if vid, err := uuid.Parse(venueIDStr); err == nil {
			venueID = vid
			query.VenueID = &venueID
		}
	}

	// Parse query parameters
	if revisionIDStr := c.Query("revision_id"); revisionIDStr != "" {
		if revisionID, err := uuid.Parse(revisionIDStr); err == nil {
			query.RevisionID = &revisionID
		}
	}

	if status := c.Query("status"); status != "" {
		if status == "published" || status == "draft" || status == "all" {
			query.Status = &status
		}
	} else if query.RevisionID == nil {
		// If no revision_id and no status provided, default to the venue's live revision
		if venueID != uuid.Nil {
			venueDetail, err := h.venueService.GetVenueDetail(c.Context(), venueID)
			if err == nil && venueDetail.LiveRevisionID != uuid.Nil {
				query.RevisionID = &venueDetail.LiveRevisionID
			}
		}
	}

	if floorIDStr := c.Query("floor_id"); floorIDStr != "" {
		if floorID, err := uuid.Parse(floorIDStr); err == nil {
			query.FloorID = &floorID
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			query.Limit = &limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			query.Offset = &offset
		}
	}

	if name := c.Query("name"); name != "" {
		query.Name = &name
	}

	if category := c.Query("category"); category != "" {
		query.Category = &category
	}

	areas, total, err := h.service.ListAreas(c.Context(), query)
	if err != nil {
		return utils.SendError(c, 500, err.Error())
	}

	return utils.SendSuccess(c, fiber.Map{
		"areas":  areas,
		"total":  total,
		"limit":  *query.Limit,
		"offset": *query.Offset,
	})
}
