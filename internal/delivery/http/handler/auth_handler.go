package handler

import (
	"inspacemap/backend/internal/models"
	"inspacemap/backend/internal/service"

	// Import hash utils
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validasi Input (Library Validator bisa ditambahkan nanti)
	if req.Email == "" || req.Password == "" || req.OrganizationName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email, Password, and Org Name are required"})
	}

	// [IMPROVISASI] Hash password di Handler sebelum masuk Service?
	// Atau di Service? Di diskusi sebelumnya kita taruh di Service/Repo.
	// Tapi Best Practice: Hashing di Service Layer.
	// Namun di kode Service auth_service.go Anda sebelumnya, Anda belum panggil Hash.
	// Jadi kita hash di sini sebelum kirim ke Service, atau update Service-nya.
	// SAYA SARANKAN HASH DI SERVICE.
	// Untuk sekarang, kita biarkan plain text masuk service (asumsi service akan diupdate pakai utils.HashPassword)

	// Update Service Anda agar menggunakan utils.HashPassword!
	// (Saya asumsikan Anda sudah update service auth_service.go untuk pakai utils)

	resp, err := h.service.Register(c.Context(), req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(resp)
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	resp, err := h.service.Login(c.Context(), req)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

// POST /api/v1/auth/invite/accept
func (h *AuthHandler) AcceptInvite(c *fiber.Ctx) error {
	var req models.AcceptInviteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	resp, err := h.service.AcceptInvitation(c.Context(), req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}
