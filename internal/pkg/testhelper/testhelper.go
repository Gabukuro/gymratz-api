package testhelper

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/exercise"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/musclegroup"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/user"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/workout"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/workoutexercise"
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
		req.Header.Add("Authorization", GenerateAuthToken(setup.EnvVariables.JWTSecret, nil))
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

func GenerateAuthToken(secret string, email *string) string {
	if email == nil {
		email = GetPointer("test@email.com")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := internalJWT.Claims{
		Email: *email,
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
	if model == nil {
		model = &user.Model{
			Name:     "John Doe",
			Email:    "john@doe.com",
			Password: "password",
		}
	}

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

func CreateUserManyWorkout(
	ctx context.Context,
	database *bun.DB,
	userID uuid.UUID,
	n int,
) map[string]workout.Model {
	workoutModels := make(map[string]workout.Model, n)

	for i := 0; i < n; i++ {
		workoutModel := CreateSingleWorkout(ctx, database, userID, i)

		workoutModels[workoutModel.ID.String()] = workoutModel
	}

	return workoutModels
}

func CreateSingleWorkout(
	ctx context.Context,
	database *bun.DB,
	userID uuid.UUID,
	index int,
) workout.Model {
	workoutModel := workout.Model{
		UserID: userID,
		Name:   fmt.Sprintf("Test workout #%d", index),
	}

	_, err := database.NewInsert().Model(&workoutModel).Exec(ctx)
	if err != nil {
		panic(err)
	}

	exercise, _ := CreateExerciseWithMuscleGroup(ctx, database,
		fmt.Sprintf("Test exercise #%d", index),
		fmt.Sprintf("Test exercise description #%d", index),
		fmt.Sprintf("Test muscle group #%d", index),
	)

	_, err = database.NewInsert().Model(&workoutexercise.Model{
		WorkoutID:   workoutModel.ID,
		ExerciseID:  exercise.ID,
		Sets:        3,
		Repetitions: GetPointer(10),
		Weight:      GetPointer(20.0),
		Duration:    GetPointer(0),
		RestTime:    60,
		Notes:       GetPointer("This is a note"),
	}).Exec(ctx)

	if err != nil {
		panic(err)
	}

	_ = database.NewSelect().Model(&workoutModel).Relation("WorkoutExercises.Exercise").Where("id = ?", workoutModel.ID).Scan(ctx)

	return workoutModel
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
