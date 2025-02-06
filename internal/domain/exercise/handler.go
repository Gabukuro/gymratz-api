package exercise

import (
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/exercise"
	"github.com/Gabukuro/gymratz-api/internal/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type (
	HTTPHandlerParams struct {
		App     *fiber.App
		Service *Service
	}

	httpHandler struct {
		service *Service
	}
)

func NewHTTPHandler(params HTTPHandlerParams) {
	httpHandler := &httpHandler{
		service: params.Service,
	}

	exerciseGroup := params.App.Group("/exercise")
	exerciseGroup.Post("/", httpHandler.CreateExercise)
	exerciseGroup.Get("/", httpHandler.ListExercises)
	exerciseGroup.Put("/:id", httpHandler.UpdateExercise)
}

func (h *httpHandler) CreateExercise(c *fiber.Ctx) error {
	var reqParams exercise.CreateExerciseRequest

	if err := c.BodyParser(&reqParams); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidRequestBody(nil))
	}

	exercise, err := h.service.CreateExercise(c.Context(), reqParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(err.Error(), fiber.StatusInternalServerError, nil))
	}

	return c.Status(fiber.StatusCreated).JSON(response.NewSuccessResponse(exercise))
}

func (h *httpHandler) ListExercises(c *fiber.Ctx) error {
	var reqParams exercise.ListExercisesQueryParams

	if err := c.QueryParser(&reqParams); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidRequestBody(nil))
	}

	if reqParams.PerPage == 0 || reqParams.Page == 0 {
		reqParams.PerPage = 10
		reqParams.Page = 1
	}

	if reqParams.PerPage > 100 {
		reqParams.PerPage = 100
	}

	exercises, total, err := h.service.ListExercises(c.Context(), reqParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(err.Error(), fiber.StatusInternalServerError, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewPaginationResponse(exercises, response.Pagination{
		Page:       reqParams.Page,
		PerPage:    reqParams.PerPage,
		TotalItems: total,
	}))
}

func (h *httpHandler) UpdateExercise(c *fiber.Ctx) error {
	exerciseID, err := uuid.ParseBytes([]byte(c.Params("id")))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidURLParam(&response.ErrorDetails{
				response.NewErrorDetail("id", "Invalid UUID format"),
			}))
	}

	var bodyRequest exercise.UpdateExerciseRequest

	if err := c.BodyParser(&bodyRequest); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidRequestBody(nil))
	}

	exercise, err := h.service.UpdateExercise(c.Context(), exerciseID, bodyRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(err.Error(), fiber.StatusInternalServerError, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(exercise))
}
