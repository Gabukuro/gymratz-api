package workout

import (
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/workout"
	"github.com/Gabukuro/gymratz-api/internal/pkg/jwt"
	"github.com/Gabukuro/gymratz-api/internal/pkg/middleware"
	"github.com/Gabukuro/gymratz-api/internal/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type (
	HTTPHandlerParams struct {
		App       *fiber.App
		Service   *Service
		JWTSecret string
	}

	httpHandler struct {
		service *Service
	}
)

func NewHTTPHandler(params HTTPHandlerParams) {
	httpHandler := &httpHandler{
		service: params.Service,
	}

	workoutGroup := params.App.Group("/workouts", middleware.AuthMiddleware(params.JWTSecret))
	workoutGroup.Post("/", httpHandler.CreateWorkout)
	workoutGroup.Get("/", httpHandler.GetUserWorkoutPaginated)
}

func (h *httpHandler) CreateWorkout(c *fiber.Ctx) error {
	var req workout.CreateWorkoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidRequestBody(nil))
	}

	workoutModel, err := h.service.CreateWorkout(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(err.Error(), fiber.StatusInternalServerError, nil))
	}

	return c.Status(fiber.StatusCreated).JSON(response.NewSuccessResponse(workoutModel))
}

func (h *httpHandler) GetUserWorkoutPaginated(c *fiber.Ctx) error {
	claims := c.Locals("session").(*jwt.Claims)

	var reqQuery workout.ListWorkoutsQueryParams
	if err := c.QueryParser(&reqQuery); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidRequestBody(nil))
	}

	workouts, total, err := h.service.ListUserWorkouts(c.Context(), claims.Email, reqQuery)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(err.Error(), fiber.StatusInternalServerError, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewPaginationResponse(workouts, response.Pagination{
		Page:       reqQuery.Page,
		PerPage:    reqQuery.PerPage,
		TotalItems: total,
	}))
}
