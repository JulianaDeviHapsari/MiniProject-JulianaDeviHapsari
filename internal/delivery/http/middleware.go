package http

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"evermos-mini/internal/pkg/security"
)

// AuthMiddleware membaca JWT dari header `token` (sesuai collection Postman)
// atau `Authorization: Bearer <jwt>`, lalu menyimpan identitas ke context.
func AuthMiddleware(tokens *security.TokenManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := c.Get("token")
		if raw == "" {
			raw = c.Get("Authorization")
		}
		raw = strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))
		if raw == "" {
			return fail(c, fiber.StatusUnauthorized, "Unauthorized", "token tidak ditemukan")
		}
		claims, err := tokens.Parse(raw)
		if err != nil {
			return fail(c, fiber.StatusUnauthorized, "Unauthorized", "token tidak valid atau kedaluwarsa")
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("is_admin", claims.IsAdmin)
		return c.Next()
	}
}

// AdminMiddleware harus dipasang setelah AuthMiddleware.
func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if isAdmin, _ := c.Locals("is_admin").(bool); !isAdmin {
			return fail(c, fiber.StatusForbidden, "Forbidden", "hanya admin yang dapat mengelola kategori")
		}
		return c.Next()
	}
}
