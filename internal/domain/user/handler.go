package user

import "github.com/gofiber/fiber/v2"

type (
	HTTPHandlerParams struct {
		App *fiber.App
	}
)

func NewHTTPHandler(params HTTPHandlerParams) {
	params.App.Get("/users", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})
}
