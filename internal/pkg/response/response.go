package response

import "github.com/gofiber/fiber/v2"

type (
	SuccessResponse[T any] struct {
		Status  string `json:"status"`             // Response status (success/error)
		Data    T      `json:"data"`               // API data
		TraceID string `json:"trace_id,omitempty"` // Trace ID (optional)
	}

	ErrorResponse struct {
		Status  string        `json:"status"`             // Response status (success/error)
		Message string        `json:"message"`            // Error message
		Code    int           `json:"code"`               // HTTP status code
		Details *ErrorDetails `json:"details,omitempty"`  // Error details (optional)
		TraceID string        `json:"trace_id,omitempty"` // Trace ID (optional)
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

func NewSuccessResponse(data any, traceID string) SuccessResponse[any] {
	return SuccessResponse[any]{
		Status:  StatusSuccess,
		Data:    data,
		TraceID: traceID,
	}
}

func NewErrorResponse(message string, code int, details *ErrorDetails, traceID string) ErrorResponse {
	if details == nil {
		details = &ErrorDetails{}
	}

	return ErrorResponse{
		Status:  StatusError,
		Message: message,
		Code:    code,
		Details: details,
		TraceID: traceID,
	}
}

func NewErrorDetail(field, message string) ErrorDetail {
	return ErrorDetail{
		Field:   field,
		Message: message,
	}
}

func NewErrorInvalidRequestBody(details *ErrorDetails, traceID string) ErrorResponse {
	return NewErrorResponse("Invalid request body.", fiber.StatusBadRequest, details, traceID)
}
