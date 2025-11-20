package middleware

import (
	"inspacemap/backend/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	CtxUserID      = "user_id"
	CtxUserEmail   = "user_email"
	CtxOrgID       = "org_id"      // Organisasi aktif di token
	CtxPermissions = "permissions" // List []string
)

// Protected: Cek Token Validitas
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing authorization header"})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		// Simpan data penting ke Context agar bisa dipakai Handler
		c.Locals(CtxUserID, claims.UserID)
		c.Locals(CtxUserEmail, claims.Email)
		c.Locals(CtxOrgID, claims.OrganizationID)
		c.Locals(CtxPermissions, claims.Permissions) // Simpan list permission

		return c.Next()
	}
}

// RBAC Guard: Middleware Factory
// Contoh penggunaan: api.Post("/venues", RequirePermission("venue:create"), CreateVenue)
func RequirePermission(requiredPerm string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil permissions dari context (yang ditaruh oleh Protected middleware)
		userPerms, ok := c.Locals(CtxPermissions).([]string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "No permissions found in context"})
		}

		// Cek apakah user punya permission yang diminta
		hasPermission := false
		for _, p := range userPerms {
			if p == requiredPerm {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You do not have permission: " + requiredPerm,
			})
		}

		return c.Next()
	}
}

// TenantGuard: Memastikan ID di URL/Header cocok dengan ID di Token
// Mencegah user Org A mengakses data Org B
func TenantGuard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenOrgID := c.Locals(CtxOrgID).(uuid.UUID)

		// Cek header X-Tenant-ID dari request
		headerOrgIDStr := c.Get("X-Tenant-ID")
		if headerOrgIDStr == "" {
			// Jika header kosong, kita asumsikan user bekerja pada org default token
			return c.Next()
		}

		// Jika ada header, pastikan COCOK dengan token
		if headerOrgIDStr != tokenOrgID.String() {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Token organization mismatch. Please switch organization context.",
			})
		}

		return c.Next()
	}
}
