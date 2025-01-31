package middleware

import (
	"strings"

	"github.com/Gabukuro/gymratz-api/internal/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")

		if token == "" && !strings.HasPrefix(token, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		tokenStr := strings.TrimPrefix(token, "Bearer ")

		claims, err := jwt.ValidateToken(tokenStr, secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		c.Locals("session", claims)

		return c.Next()
	}
}
