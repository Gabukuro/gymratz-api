package workout

import (
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/workout"
	"github.com/Gabukuro/gymratz-api/internal/pkg/jwt"
	"github.com/Gabukuro/gymratz-api/internal/pkg/middleware"
	"github.com/Gabukuro/gymratz-api/internal/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	workoutGroup.Put("/:id", httpHandler.UpdateWorkout)
	workoutGroup.Put("/:workoutID/exercises/:workoutExerciseID", httpHandler.UpdateWorkoutExercise)
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

func (h *httpHandler) UpdateWorkout(c *fiber.Ctx) error {
	workoutID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidURLParam(&response.ErrorDetails{
				response.NewErrorDetail("id", "Invalid UUID format"),
			}))
	}

	var reqParams workout.UpdateWorkoutRequest
	if err := c.BodyParser(&reqParams); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidRequestBody(nil))
	}

	workoutModel, err := h.service.UpdateWorkout(c.Context(), workoutID, reqParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(err.Error(), fiber.StatusInternalServerError, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(workoutModel))
}

func (h *httpHandler) UpdateWorkoutExercise(c *fiber.Ctx) error {
	workoutID, err := uuid.Parse(c.Params("workoutID"))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidURLParam(&response.ErrorDetails{
				response.NewErrorDetail("workoutID", "Invalid UUID format"),
			}))
	}

	workoutExerciseID, err := uuid.Parse(c.Params("workoutExerciseID"))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidURLParam(&response.ErrorDetails{
				response.NewErrorDetail("workoutExerciseID", "Invalid UUID format"),
			}))
	}

	var reqParams workout.UpdateWorkoutExerciseRequest
	if err := c.BodyParser(&reqParams); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidRequestBody(nil))
	}

	workoutExercise, err := h.service.UpdateWorkoutExercise(c.Context(), workoutID, workoutExerciseID, reqParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(err.Error(), fiber.StatusInternalServerError, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(workoutExercise))
}
