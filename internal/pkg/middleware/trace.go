package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func TraceMiddleware(c *fiber.Ctx) error {
	c.Set("X-Request-ID", uuid.New().String())
	return c.Next()
}
