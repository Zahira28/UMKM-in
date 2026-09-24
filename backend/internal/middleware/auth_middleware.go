package middleware

import (
	"strings"

	"backend/internal/utils"
	"backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func Protected(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Akses ditolak: Header Authorization tidak ditemukan", nil)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return response.Error(c, fiber.StatusUnauthorized, "Format token tidak valid, gunakan format 'Bearer <token>'", nil)
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Sesi login tidak valid atau sudah kedaluwarsa", nil)
		}

		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}
