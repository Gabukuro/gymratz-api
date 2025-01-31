package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TraceMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := uuid.New().String()
		c.Set("X-Request-ID", requestID)

		c.Response().Header.Set("X-Request-ID", requestID)

		return c.Next()
	}
}
