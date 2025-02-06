package response

import (
	"math"

	"github.com/gofiber/fiber/v2"
)

type (
	SuccessResponse[T any] struct {
		Status string `json:"status"` // Response status (success/error)
		Data   T      `json:"data"`   // API data
	}

	PaginationResponse[T any] struct {
		SuccessResponse[T]
		Pagination Pagination `json:"pagination"` // Pagination metadata
	}

	Pagination struct {
		Page       int `json:"page"`        // Current page number
		PerPage    int `json:"per_page"`    // Number of items per page
		TotalItems int `json:"total_items"` // Total number of items
		TotalPages int `json:"total_pages"` // Total number of pages
	}

	ErrorResponse struct {
		Status  string        `json:"status"`            // Response status (success/error)
		Message string        `json:"message"`           // Error message
		Code    int           `json:"code"`              // HTTP status code
		Details *ErrorDetails `json:"details,omitempty"` // Error details (optional)
	}

	ErrorDetail struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}

	ErrorDetails []ErrorDetail
)

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

func NewSuccessResponse(data any) SuccessResponse[any] {
	return SuccessResponse[any]{
		Status: StatusSuccess,
		Data:   data,
	}
}

func NewPaginationResponse(data any, pagination Pagination) PaginationResponse[any] {
	pagination.SetTotalPages()

	return PaginationResponse[any]{
		SuccessResponse: SuccessResponse[any]{
			Status: StatusSuccess,
			Data:   data,
		},
		Pagination: pagination,
	}
}

func (p *Pagination) SetTotalPages() {
	if p.TotalItems == 0 {
		p.TotalPages = 0
	} else {
		p.TotalPages = int(math.Ceil(float64(p.TotalItems) / float64(p.PerPage)))
	}
}

func NewErrorResponse(message string, code int, details *ErrorDetails) ErrorResponse {
	if details == nil {
		details = &ErrorDetails{}
	}

	return ErrorResponse{
		Status:  StatusError,
		Message: message,
		Code:    code,
		Details: details,
	}
}

func NewErrorDetail(field, message string) ErrorDetail {
	return ErrorDetail{
		Field:   field,
		Message: message,
	}
}

func NewErrorInvalidRequestBody(details *ErrorDetails) ErrorResponse {
	return NewErrorResponse("Invalid request body.", fiber.StatusBadRequest, details)
}

func NewErrorInvalidURLParam(details *ErrorDetails) ErrorResponse {
	return NewErrorResponse("Invalid URL parameter.", fiber.StatusBadRequest, details)
}
