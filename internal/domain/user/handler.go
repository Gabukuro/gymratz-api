package user

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/user"
)

type (
	HTTPHandlerParams struct {
		App     *fiber.App
		Service *Service
	}

	httpHandler struct {
		app     *fiber.App
		service *Service
	}
)

func NewHTTPHandler(params HTTPHandlerParams) {
	httpHandler := &httpHandler{
		app:     params.App,
		service: params.Service,
	}

	params.App.Post("/register", httpHandler.RegisterUser)
	params.App.Post("/login", httpHandler.LoginUser)
}

func (h *httpHandler) RegisterUser(c *fiber.Ctx) error {
	var req user.RegisterUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if !req.Validate() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.service.CreateUser(c.Context(), req.Name, req.Email, req.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user created",
	})
}

func (h *httpHandler) LoginUser(c *fiber.Ctx) error {
	var req user.LoginUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if !req.Validate() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	token, err := h.service.LoginUser(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}
