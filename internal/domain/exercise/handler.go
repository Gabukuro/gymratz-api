package exercise

import (
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/exercise"
	"github.com/Gabukuro/gymratz-api/internal/pkg/response"
	"github.com/gofiber/fiber/v2"
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
	var reqParams exercise.ListExercisesRequest

	if err := c.BodyParser(&reqParams); err != nil {
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
		TotalPages: total / reqParams.PerPage,
	}))
}
