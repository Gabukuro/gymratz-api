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
}

func createMuscleGroup(ctx context.Context, model *musclegroup.Model) musclegroup.Model {
	_, err := database.DB().NewInsert().Model(model).Exec(ctx)
	if err != nil {
		panic(err)
	}

	return *model
}
