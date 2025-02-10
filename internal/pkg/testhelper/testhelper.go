package testhelper

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/exercise"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/musclegroup"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/user"
	internalJWT "github.com/Gabukuro/gymratz-api/internal/pkg/jwt"
	"github.com/Gabukuro/gymratz-api/internal/pkg/response"
	"github.com/Gabukuro/gymratz-api/internal/pkg/setup"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func RunRequest(
	setup *setup.Setup,
	method string,
	path string,
	body any,
	header map[string]string,
) (*http.Response, error) {
	req := httptest.NewRequest(
		method,
		path,
		parseBodyToStringReader(body),
	)

	if _, ok := header["Authorization"]; !ok {
		req.Header.Add("Authorization", generateAuthToken(setup.EnvVariables.JWTSecret))
	}

	req.Header.Add("Content-Type", "application/json")
	for key, value := range header {
		req.Header.Add(key, value)
	}

	return setup.App.Test(req, -1)
}

func parseBodyToStringReader(requestBody any) *strings.Reader {
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}
	return strings.NewReader(string(jsonBody))
}

func generateAuthToken(secret string) string {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := internalJWT.Claims{
		Email: "test@email.com",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "gymratz-api-test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}

	return "Bearer " + signedToken
}

func ParseSuccessResponseBody[data any](body io.ReadCloser) response.SuccessResponse[data] {
	var responseBody response.SuccessResponse[data]

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		panic(err)
	}

	return responseBody
}

func ParsePaginationResponseBody[data any](body io.ReadCloser) response.PaginationResponse[data] {
	var responseBody response.PaginationResponse[data]

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		panic(err)
	}

	return responseBody
}

func ParseErrorResponseBody(body io.ReadCloser) response.ErrorResponse {
	var responseBody response.ErrorResponse

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		panic(err)
	}

	return responseBody
}

func GetPointer[data any](value data) *data {
	return &value
}

func CleanUpDatabase(ctx context.Context, database *bun.DB) {
	dropExercises(ctx, database)
	dropMuscleGroups(ctx, database)
	dropUsers(ctx, database)
}

func CreateUser(ctx context.Context, database *bun.DB, model *user.Model) user.Model {
	model.HashPassword()
	_, err := database.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		panic(err)
	}

	return *model
}

func dropUsers(ctx context.Context, database *bun.DB) {
	_, err := database.NewDelete().Model(&user.Model{}).Where("1 = 1").Exec(ctx)
	if err != nil {
		panic(err)
	}
}

func CreateExerciseWithMuscleGroup(
	ctx context.Context,
	database *bun.DB,
	name string,
	description string,
	muscleGroupName string,
) (exercise.Model, musclegroup.Model) {
	muscleGroup := createMuscleGroup(ctx, database, &musclegroup.Model{
		Name: muscleGroupName,
	})

	exerciseEntity := createExercise(ctx, database, &exercise.Model{
		Name:        name,
		Description: description,
	})

	createExerciseMuscleGroupAssociation(ctx, database, exerciseEntity, muscleGroup)

	return exerciseEntity, muscleGroup
}

func createExerciseMuscleGroupAssociation(
	ctx context.Context,
	database *bun.DB,
	exerciseModel exercise.Model,
	muscleGroup musclegroup.Model,
) {
	_, err := database.NewInsert().Model(&exercise.ExerciseMuscleGroupModel{
		ExerciseID:    exerciseModel.ID,
		MuscleGroupID: muscleGroup.ID,
	}).Exec(ctx)
	if err != nil {
		panic(err)
	}
}

func createExercise(ctx context.Context, database *bun.DB, model *exercise.Model) exercise.Model {
	_, err := database.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		panic(err)
	}

	return *model
}

func getExerciseByID(ctx context.Context, database *bun.DB, id uuid.UUID) *exercise.Model {
	model := &exercise.Model{}
	err := database.NewSelect().Model(model).Where("id = ?", id).Scan(ctx)
	if err != nil && err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}

	return model
}

func dropExercises(ctx context.Context, database *bun.DB) {
	_, err := database.NewDelete().Model(&exercise.Model{}).Where("1 = 1").Exec(ctx)
	if err != nil {
		panic(err)
	}
}

func createMuscleGroup(ctx context.Context, database *bun.DB, model *musclegroup.Model) musclegroup.Model {
	_, err := database.NewInsert().Model(model).Exec(ctx)
	if err != nil {
		panic(err)
	}

	return *model
}

func dropMuscleGroups(ctx context.Context, database *bun.DB) {
	_, err := database.NewDelete().Model(&musclegroup.Model{}).Where("1 = 1").Exec(ctx)
	if err != nil {
		panic(err)
	}
}
