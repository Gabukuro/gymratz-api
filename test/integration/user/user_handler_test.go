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
                "email": "fake@email.com",
                "password": "password"
            }`),
		)

		resp, err := runRequest(req, setup)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("should not create a new user with the same email", func(t *testing.T) {
		// Create a user with the same email
		_ = createUser(ctx, user.Model{
			Name:     "test",
			Email:    "fake1@email.com",
			Password: "password",
		})

		// Try to create a new user with the same email
		req := httptest.NewRequest(
			http.MethodPost,
			"/register",
			strings.NewReader(`{
				"name": "test",
				"email": "fake1@email.com",
				"password": "password"
			}`),
		)

		resp, err := runRequest(req, setup)

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
		_ = createUser(ctx, user.Model{
			Name:     "test",
			Email:    "fake2@email.com",
			Password: "password123",
		})

		req := httptest.NewRequest(
			http.MethodPost,
			"/login",
			strings.NewReader(`{
				"email": "fake2@email.com",
				"password": "password123"
			}`),
		)

		resp, err := runRequest(req, setup)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var successResponse response.SuccessResponse[user.LoginUserResponse]

		body, _ := io.ReadAll(resp.Body)
		assert.Nil(t, json.Unmarshal(body, &successResponse))
		assert.Equal(t, response.StatusSuccess, successResponse.Status)
		assert.NotEmpty(t, successResponse.Data.Token)
	})

	t.Run("should not login a user with invalid credentials", func(t *testing.T) {
		_ = createUser(ctx, user.Model{
			Name:     "test",
			Email:    "fake3@email.com",
			Password: "password123",
		})

		req := httptest.NewRequest(
			http.MethodPost,
			"/login",
			strings.NewReader(`{
				"email": "fake3@email.com",
				"password": "wrongpassword"
			}`),
		)

		resp, err := runRequest(req, setup)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var errorResponse response.ErrorResponse

		body, _ := io.ReadAll(resp.Body)
		assert.Nil(t, json.Unmarshal(body, &errorResponse))
		assert.Equal(t, http.StatusUnauthorized, errorResponse.Code)
		assert.Equal(t, response.StatusError, errorResponse.Status)
		assert.Equal(t, "Unauthorized", errorResponse.Message)

		assert.Empty(t, *errorResponse.Details)
	})

	t.Run("should get user profile", func(t *testing.T) {
		userModel := createUser(ctx, user.Model{
			Name:     "test",
			Email:    "fake4@email.com",
			Password: "password123",
		})

		// Login the user
		req := httptest.NewRequest(
			http.MethodPost,
			"/login",
			strings.NewReader(`{
				"email": "fake4@email.com",
				"password": "password123"
			}`),
		)

		resp, err := runRequest(req, setup)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var loginResponse response.SuccessResponse[user.LoginUserResponse]

		body, _ := io.ReadAll(resp.Body)
		assert.Nil(t, json.Unmarshal(body, &loginResponse))
		assert.Equal(t, response.StatusSuccess, loginResponse.Status)

		// Get the user profile
		req = httptest.NewRequest(
			http.MethodGet,
			"/user/profile",
			nil,
		)
		req.Header.Add("Authorization", "Bearer "+*loginResponse.Data.Token)

		resp, err = runRequest(req, setup)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var profileResponse response.SuccessResponse[user.GetUserProfileResponse]

		body, _ = io.ReadAll(resp.Body)
		assert.Nil(t, json.Unmarshal(body, &profileResponse))
		assert.Equal(t, response.StatusSuccess, profileResponse.Status)

		assert.Equal(t, userModel.ID, profileResponse.Data.ID)
		assert.Equal(t, userModel.Name, profileResponse.Data.Name)
		assert.Equal(t, userModel.Email, profileResponse.Data.Email)
	})
}

func createUser(ctx context.Context, userModel user.Model) user.Model {
	err := userModel.HashPassword()
	if err != nil {
		panic(err)
	}

	_, err = database.DB().NewInsert().Model(&userModel).Exec(ctx)
	if err != nil {
		panic(err)
	}

	return userModel
}

func runRequest(req *http.Request, setup *setup.Setup) (*http.Response, error) {
	req.Header.Add("Content-Type", "application/json")

	return setup.App.Test(req, -1)
}
