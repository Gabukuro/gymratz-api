package user

import (
	"github.com/Gabukuro/gymratz-api/internal/pkg/response"
)

var (
	// ErrorNameIsRequired is the error message for name is required
	ErrorNameIsRequired response.ErrorDetail = response.NewErrorDetail("name", "Name is required")

	// ErrorEmailIsRequired is the error message for email is required
	ErrorEmailIsRequired response.ErrorDetail = response.NewErrorDetail("email", "Email is required")

	// ErrorPasswordIsRequired is the error message for password is required
	ErrorPasswordIsRequired response.ErrorDetail = response.NewErrorDetail("password", "Password is required")
)
