package musclegroup

import (
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/musclegroup"
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

	muscleGroupGroup := params.App.Group("/muscle-groups", middleware.AuthMiddleware(params.JWTSecret))
	muscleGroupGroup.Post("/", httpHandler.CreateMuscleGroup)
	muscleGroupGroup.Get("/", httpHandler.ListMuscleGroups)
	muscleGroupGroup.Put("/:id", httpHandler.UpdateMuscleGroup)
	muscleGroupGroup.Delete("/:id", httpHandler.DeleteMuscleGroup)
}

func (h *httpHandler) CreateMuscleGroup(c *fiber.Ctx) error {
	var reqParams musclegroup.CreateMuscleGroupRequest

	if err := c.BodyParser(&reqParams); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidRequestBody(nil))
	}

	muscleGroup, err := h.service.CreateMuscleGroup(c.Context(), reqParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(err.Error(), fiber.StatusInternalServerError, nil))
	}

	return c.Status(fiber.StatusCreated).JSON(response.NewSuccessResponse(muscleGroup))
}

func (h *httpHandler) ListMuscleGroups(c *fiber.Ctx) error {
	var reqParams musclegroup.ListMuscleGroupsQueryParams

	if err := c.QueryParser(&reqParams); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidRequestBody(nil))
	}

	reqParams.ValidateAndSetDefaults()

	muscleGroups, total, err := h.service.ListMuscleGroups(c.Context(), reqParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(err.Error(), fiber.StatusInternalServerError, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewPaginationResponse(muscleGroups, response.Pagination{
		Page:       reqParams.Page,
		PerPage:    reqParams.PerPage,
		TotalItems: total,
	}))
}

func (h *httpHandler) UpdateMuscleGroup(c *fiber.Ctx) error {
	muscleGroupID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidRequestBody(&response.ErrorDetails{
				response.NewErrorDetail("id", "Invalid UUID format"),
			}))
	}

	var bodyRequest musclegroup.UpdateMuscleGroupRequest
	if err := c.BodyParser(&bodyRequest); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidRequestBody(nil))
	}

	muscleGroup, err := h.service.UpdateMuscleGroup(c.Context(), muscleGroupID, bodyRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(err.Error(), fiber.StatusInternalServerError, nil))
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(muscleGroup))
}

func (h *httpHandler) DeleteMuscleGroup(c *fiber.Ctx) error {
	muscleGroupID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			response.NewErrorInvalidRequestBody(&response.ErrorDetails{
				response.NewErrorDetail("id", "Invalid UUID format"),
			}))
	}

	err = h.service.DeleteMuscleGroup(c.Context(), muscleGroupID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.NewErrorResponse(err.Error(), fiber.StatusInternalServerError, nil))
	}

	return c.SendStatus(fiber.StatusNoContent)
}
