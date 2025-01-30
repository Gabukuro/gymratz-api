package response

import "github.com/gofiber/fiber/v2"

type (
	SuccessResponse struct {
		Status  string `json:"status"`             // Status da resposta (success/error)
		Data    any    `json:"data"`               // Dados retornados pela API
		TraceID string `json:"trace_id,omitempty"` // ID de rastreamento (opcional)
	}

	ErrorResponse struct {
		Status  string        `json:"status"`             // Status da resposta (success/error)
		Message string        `json:"message"`            // Mensagem de erro
		Code    int           `json:"code"`               // Código de erro (HTTP status code)
		Details *ErrorDetails `json:"details,omitempty"`  // Detalhes do erro (opcional)
		TraceID string        `json:"trace_id,omitempty"` // ID de rastreamento (opcional)
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

func NewSuccessResponse(data any, traceID string) SuccessResponse {
	return SuccessResponse{
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
	return NewErrorResponse("invalid request body", fiber.StatusBadRequest, details, traceID)
}
