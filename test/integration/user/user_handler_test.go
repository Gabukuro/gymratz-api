package user_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Gabukuro/gymratz-api/internal/infra/database"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/user"
	"github.com/Gabukuro/gymratz-api/internal/pkg/response"
	"github.com/Gabukuro/gymratz-api/internal/pkg/setup"
	"github.com/Gabukuro/gymratz-api/internal/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

const (
	registerUserPath = "/register"
	loginUserPath    = "/login"
)

func TestUserHandler(t *testing.T) {
	t.Parallel()

	os.Setenv("GO_ENV", "test")
	setup, ctx := setup.Init()
	defer database.CloseTestDB(ctx)

	t.Run("should create a new user", func(t *testing.T) {
		resp, err := testhelper.RunRequest(
			setup,
			http.MethodPost,
			registerUserPath,
			user.RegisterUserRequest{
				Name:     "test",
				Email:    "fake@email.com",
				Password: "password",
			},
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("should not create a new user with the same email", func(t *testing.T) {
		fakeEmail := "fake1@email.com"

		// Create a user with the same email
		_ = createUser(ctx, user.Model{
			Name:     "test",
			Email:    fakeEmail,
			Password: "password",
		})

		resp, err := testhelper.RunRequest(
			setup,
			http.MethodPost,
			registerUserPath,
			user.RegisterUserRequest{
				Name:     "test",
				Email:    fakeEmail,
				Password: "password",
			},
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		errorResponse := testhelper.ParseErrorResponseBody(resp.Body)

		assert.Equal(t, http.StatusBadRequest, errorResponse.Code)
		assert.Equal(t, response.StatusError, errorResponse.Status)
		assert.Equal(t, "Invalid request body.", errorResponse.Message)

		assert.Len(t, *errorResponse.Details, 1)
		details := *errorResponse.Details
		assert.Equal(t, "email", details[0].Field)
		assert.Equal(t, "It looks like this email is already registered on our platform", details[0].Message)
	})

	t.Run("should login a user", func(t *testing.T) {
		userCreated := createUser(ctx, user.Model{
			Name:     "test",
			Email:    "fake2@email.com",
			Password: "password123",
		})

		resp, err := testhelper.RunRequest(
			setup,
			http.MethodPost,
			loginUserPath,
			user.LoginUserRequest{
				Email:    userCreated.Email,
				Password: "password123",
			},
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		successResponse := testhelper.ParseSuccessResponseBody[user.LoginUserResponse](resp.Body)
		assert.Equal(t, response.StatusSuccess, successResponse.Status)
		assert.NotEmpty(t, successResponse.Data.Token)
	})

	t.Run("should not login a user with invalid credentials", func(t *testing.T) {
		userCreated := createUser(ctx, user.Model{
			Name:     "test",
			Email:    "fake3@email.com",
			Password: "password123",
		})

		resp, err := testhelper.RunRequest(
			setup,
			http.MethodPost,
			loginUserPath,
			user.LoginUserRequest{
				Email:    userCreated.Email,
				Password: "wrongpassword",
			},
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		errorResponse := testhelper.ParseErrorResponseBody(resp.Body)

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
		resp, err := testhelper.RunRequest(
			setup,
			http.MethodPost,
			loginUserPath,
			user.LoginUserRequest{
				Email:    userModel.Email,
				Password: "password123",
			},
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		loginResponse := testhelper.ParseSuccessResponseBody[user.LoginUserResponse](resp.Body)

		assert.Equal(t, response.StatusSuccess, loginResponse.Status)
		assert.NotNil(t, loginResponse.Data.Token)

		// Get the user profile
		resp, err = testhelper.RunRequest(
			setup,
			http.MethodGet,
			"/users/profile",
			nil,
			map[string]string{
				"Authorization": fmt.Sprintf("Bearer %s", *loginResponse.Data.Token),
			},
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		profileResponse := testhelper.ParseSuccessResponseBody[user.GetUserProfileResponse](resp.Body)

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
