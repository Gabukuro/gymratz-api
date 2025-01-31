package response

import "github.com/gofiber/fiber/v2"

type (
	SuccessResponse[T any] struct {
		Status string `json:"status"` // Response status (success/error)
		Data   T      `json:"data"`   // API data
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
