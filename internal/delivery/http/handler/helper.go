package handler

import (
	"inspacemap/backend/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func getUserID(c *fiber.Ctx) uuid.UUID {
	id, ok := c.Locals(middleware.CtxUserID).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return id
}

func getOrgID(c *fiber.Ctx) uuid.UUID {
	id, ok := c.Locals(middleware.CtxOrgID).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return id
}
