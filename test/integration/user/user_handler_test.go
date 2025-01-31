package user_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Gabukuro/gymratz-api/internal/infra/database"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/user"
	"github.com/Gabukuro/gymratz-api/internal/pkg/response"
	"github.com/Gabukuro/gymratz-api/internal/pkg/setup"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler(t *testing.T) {
	t.Parallel()

	os.Setenv("GO_ENV", "test")
	setup, ctx := setup.Init()
	defer database.CloseTestDB(ctx)

	t.Run("should create a new user", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodPost,
			"/register",
			strings.NewReader(`{
                "name": "test",
                "email": "test@email.com",
                "password": "password"
            }`),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, err := setup.App.Test(req, -1)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("should not create a new user with the same email", func(t *testing.T) {
		// Create a user with the same email
		createUser(ctx, user.Model{
			Name:     "test",
			Email:    "email@email.com",
			Password: "password",
		})

		// Try to create a new user with the same email
		req := httptest.NewRequest(
			http.MethodPost,
			"/register",
			strings.NewReader(`{
				"name": "test",
				"email": "email@email.com",
				"password": "password"
			}`),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, err := setup.App.Test(req, -1)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errorResponse response.ErrorResponse

		body, _ := io.ReadAll(resp.Body)
		assert.Nil(t, json.Unmarshal(body, &errorResponse))
		assert.Equal(t, http.StatusBadRequest, errorResponse.Code)
		assert.Equal(t, response.StatusError, errorResponse.Status)
		assert.Equal(t, "Invalid request body.", errorResponse.Message)

		assert.Len(t, *errorResponse.Details, 1)
		details := *errorResponse.Details
		assert.Equal(t, "email", details[0].Field)
		assert.Equal(t, "It looks like this email is already registered on our platform", details[0].Message)
	})

	t.Run("should login a user", func(t *testing.T) {
		createUser(ctx, user.Model{
			Name:     "test",
			Email:    "fake@email.com",
			Password: "password123",
		})

		req := httptest.NewRequest(
			http.MethodPost,
			"/login",
			strings.NewReader(`{
				"email": "fake@email.com",
				"password": "password123"
			}`),
		)

		req.Header.Set("Content-Type", "application/json")

		resp, err := setup.App.Test(req, -1)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var successResponse response.SuccessResponse[user.LoginUserResponse]

		body, _ := io.ReadAll(resp.Body)
		assert.Nil(t, json.Unmarshal(body, &successResponse))
		assert.Equal(t, response.StatusSuccess, successResponse.Status)
		assert.NotEmpty(t, successResponse.Data.Token)
	})
}

func createUser(ctx context.Context, userModel user.Model) {
	err := userModel.HashPassword()
	if err != nil {
		panic(err)
	}

	_, err = database.DB().NewInsert().Model(&userModel).Exec(ctx)
	if err != nil {
		panic(err)
	}
}
