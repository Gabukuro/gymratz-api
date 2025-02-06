package exercise__test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/Gabukuro/gymratz-api/internal/infra/database"
	exerciseEntity "github.com/Gabukuro/gymratz-api/internal/pkg/entity/exercise"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/musclegroup"
	"github.com/Gabukuro/gymratz-api/internal/pkg/response"
	"github.com/Gabukuro/gymratz-api/internal/pkg/setup"
	"github.com/Gabukuro/gymratz-api/internal/pkg/testhelper"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestExerciseHandler(t *testing.T) {
	t.Parallel()

	os.Setenv("GO_ENV", "test")
	setup, ctx := setup.Init()
	defer database.CloseTestDB(ctx)

	t.Run("should create a new exercise", func(t *testing.T) {
		testMuscleGroup := createMuscleGroup(ctx, &musclegroup.Model{
			Name: "test muscle group",
		})

		requestBody := exerciseEntity.CreateExerciseRequest{
			Name:           "test exercise",
			Description:    "test description",
			MuscleGroupIDs: []uuid.UUID{testMuscleGroup.ID},
		}

		resp, err := testhelper.RunRequest(setup,
			http.MethodPost,
			"/exercise",
			requestBody,
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		responseParsed := testhelper.ParseSuccessResponseBody[exerciseEntity.Model](resp.Body)
		assert.Equal(t, response.StatusSuccess, responseParsed.Status)
		assert.Equal(t, requestBody.Name, responseParsed.Data.Name)
		assert.Equal(t, requestBody.Description, responseParsed.Data.Description)
		assert.Equal(t, 1, len(responseParsed.Data.MuscleGroups))
		assert.Equal(t, testMuscleGroup.ID, responseParsed.Data.MuscleGroups[0].ID)
	})

	t.Run("should list exercises", func(t *testing.T) {
		// Clean up the database
		cleanUpDatabase(ctx)

		// Create a new exercise
		testExercise := createExercise(ctx, &exerciseEntity.Model{
			Name:        "test exercise",
			Description: "test description",
		})

		// Send a request to list exercises
		resp, err := testhelper.RunRequest(setup,
			http.MethodGet,
			"/exercise?name=test",
			nil,
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		responseParsed := testhelper.ParsePaginationResponseBody[[]*exerciseEntity.Model](resp.Body)
		assert.Equal(t, response.StatusSuccess, responseParsed.Status)

		// Check the response pagination metadata
		assert.Equal(t, 1, responseParsed.Pagination.Page)
		assert.Equal(t, 10, responseParsed.Pagination.PerPage) // Default per page
		assert.Equal(t, 1, responseParsed.Pagination.TotalItems)
		assert.Equal(t, 1, responseParsed.Pagination.TotalPages)

		assert.Equal(t, 1, len(responseParsed.Data))
		assert.Equal(t, testExercise.ID, responseParsed.Data[0].ID)
		assert.Equal(t, testExercise.Name, responseParsed.Data[0].Name)
		assert.Equal(t, testExercise.Description, responseParsed.Data[0].Description)
	})
}

func createExercise(ctx context.Context, model *exerciseEntity.Model) exerciseEntity.Model {
	_, err := database.DB().NewInsert().Model(model).Exec(ctx)
	if err != nil {
		panic(err)
	}

	return *model
}

func cleanUpDatabase(ctx context.Context) {
	dropExercises(ctx)
	dropMuscleGroups(ctx)
}

func dropExercises(ctx context.Context) {
	_, err := database.DB().NewDelete().Model(&exerciseEntity.Model{}).Where("1 = 1").Exec(ctx)
	if err != nil {
		panic(err)
	}
}

func dropMuscleGroups(ctx context.Context) {
	_, err := database.DB().NewDelete().Model(&musclegroup.Model{}).Where("1 = 1").Exec(ctx)
	if err != nil {
		panic(err)
	}
}

func createMuscleGroup(ctx context.Context, model *musclegroup.Model) musclegroup.Model {
	_, err := database.DB().NewInsert().Model(model).Exec(ctx)
	if err != nil {
		panic(err)
	}

	return *model
}
