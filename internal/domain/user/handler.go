package user

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/user"
	"github.com/Gabukuro/gymratz-api/internal/pkg/jwt"
	"github.com/Gabukuro/gymratz-api/internal/pkg/middleware"
	"github.com/Gabukuro/gymratz-api/internal/pkg/response"
)

type (
	HTTPHandlerParams struct {
		App       *fiber.App
		Service   *Service
		JWTSecret string
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

	userGroup := params.App.Group("/users", middleware.AuthMiddleware(params.JWTSecret))
	userGroup.Get("/profile", httpHandler.GetUserProfile)
}

func (h *httpHandler) RegisterUser(c *fiber.Ctx) error {
	var req user.RegisterUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidRequestBody(nil))
	}

	if validationErr := req.Validate(); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorInvalidRequestBody(validationErr))
	}

	if err := h.service.CreateUser(c.Context(), req.Name, req.Email, req.Password); err != nil {
		if strings.Contains(err.Error(), "idx_users_email") {
			return c.Status(fiber.StatusBadRequest).JSON(
				response.NewErrorInvalidRequestBody(&response.ErrorDetails{
					response.NewErrorDetail("email", "It looks like this email is already registered on our platform"),
				}))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(err.Error(), fiber.StatusInternalServerError, nil))
	}

	return c.Status(fiber.StatusCreated).JSON(response.NewSuccessResponse(
		user.RegisterUserResponse{
			Message: "User created successfully",
		}))
}

func (h *httpHandler) LoginUser(c *fiber.Ctx) error {
	var req user.LoginUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorInvalidRequestBody(nil))
	}

	if validationErr := req.Validate(); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.NewErrorInvalidRequestBody(validationErr))
	}

	token, err := h.service.LoginUser(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			response.NewErrorResponse("Unauthorized", fiber.StatusUnauthorized, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(user.LoginUserResponse{Token: token}))
}

func (h *httpHandler) GetUserProfile(c *fiber.Ctx) error {
	claims := c.Locals("session").(*jwt.Claims)

	userModel, err := h.service.GetUserProfile(c.Context(), claims.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(err.Error(), fiber.StatusInternalServerError, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(user.GetUserProfileResponse{
		ID:    userModel.ID,
		Name:  userModel.Name,
		Email: userModel.Email,
	}))
}
