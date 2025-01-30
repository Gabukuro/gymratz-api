package user_test

import (
	"testing"

	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/user"
	"github.com/Gabukuro/gymratz-api/internal/pkg/response"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUserRequestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		request  user.RegisterUserRequest
		expected *response.ErrorDetails
	}{
		{
			name: "valid request",
			request: user.RegisterUserRequest{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "password123",
			},
			expected: nil,
		},
		{
			name: "missing name",
			request: user.RegisterUserRequest{
				Email:    "john.doe@example.com",
				Password: "password123",
			},
			expected: &response.ErrorDetails{user.ErrorNameIsRequired},
		},
		{
			name: "missing email",
			request: user.RegisterUserRequest{
				Name:     "John Doe",
				Password: "password123",
			},
			expected: &response.ErrorDetails{user.ErrorEmailIsRequired},
		},
		{
			name: "missing password",
			request: user.RegisterUserRequest{
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
			expected: &response.ErrorDetails{user.ErrorPasswordIsRequired},
		},
		{
			name:    "missing all fields",
			request: user.RegisterUserRequest{},
			expected: &response.ErrorDetails{
				user.ErrorNameIsRequired,
				user.ErrorEmailIsRequired,
				user.ErrorPasswordIsRequired,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.request.Validate()
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestLoginUserRequestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		request  user.LoginUserRequest
		expected *response.ErrorDetails
	}{
		{
			name: "valid request",
			request: user.LoginUserRequest{
				Email:    "john.doe@example.com",
				Password: "password123",
			},
			expected: nil,
		},
		{
			name: "missing email",
			request: user.LoginUserRequest{
				Password: "password123",
			},
			expected: &response.ErrorDetails{user.ErrorEmailIsRequired},
		},
		{
			name: "missing password",
			request: user.LoginUserRequest{
				Email: "john.doe@example.com",
			},
			expected: &response.ErrorDetails{user.ErrorPasswordIsRequired},
		},
		{
			name:    "missing all fields",
			request: user.LoginUserRequest{},
			expected: &response.ErrorDetails{
				user.ErrorEmailIsRequired,
				user.ErrorPasswordIsRequired,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.request.Validate()
			assert.Equal(t, tt.expected, err)
		})
	}
}
