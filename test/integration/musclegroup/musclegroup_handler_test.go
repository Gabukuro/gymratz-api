package musclegroup_test

import (
	"context"
	"net/http"
	"os"
	"slices"
	"testing"

	"github.com/Gabukuro/gymratz-api/internal/infra/database"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/musclegroup"
	"github.com/Gabukuro/gymratz-api/internal/pkg/response"
	"github.com/Gabukuro/gymratz-api/internal/pkg/setup"
	"github.com/Gabukuro/gymratz-api/internal/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestMuscleGroupHandler(t *testing.T) {
	t.Parallel()

	os.Setenv("GO_ENV", "test")
	setup, ctx := setup.Init()
	defer database.CloseTestDB(ctx)

	t.Run("should create a new muscle group", func(t *testing.T) {
		resp, err := testhelper.RunRequest(
			setup,
			http.MethodPost,
			"/muscle-groups",
			&musclegroup.CreateMuscleGroupRequest{
				Name: "Chest",
			},
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		responseParsed := testhelper.ParseSuccessResponseBody[musclegroup.Model](resp.Body)
		assert.Equal(t, response.StatusSuccess, responseParsed.Status)
		assert.NotNil(t, responseParsed.Data)
		assert.NotNil(t, responseParsed.Data.ID)
		assert.Equal(t, "Chest", responseParsed.Data.Name)
	})

	t.Run("should return a list of muscle groups", func(t *testing.T) {
		cleanUpDatabase(ctx)

		muscleGroups := []string{"Chest", "Back", "Legs", "Shoulders", "Arms"}
		for _, name := range muscleGroups {
			_ = createMuscleGroup(ctx, &musclegroup.Model{
				Name: name,
			})
		}

		resp, err := testhelper.RunRequest(
			setup,
			http.MethodGet,
			"/muscle-groups",
			nil,
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		responseParsed := testhelper.ParsePaginationResponseBody[[]musclegroup.Model](resp.Body)
		assert.Equal(t, response.StatusSuccess, responseParsed.Status)

		// Check the response pagination metadata
		assert.Equal(t, 1, responseParsed.Pagination.Page)
		assert.Equal(t, 10, responseParsed.Pagination.PerPage) // Default per page
		assert.Equal(t, 5, responseParsed.Pagination.TotalItems)
		assert.Equal(t, 1, responseParsed.Pagination.TotalPages)

		// Check the response data
		assert.NotEmpty(t, responseParsed.Data)
		assert.Len(t, responseParsed.Data, 5)
		for _, muscleGroup := range responseParsed.Data {
			assert.True(t, slices.Contains(muscleGroups, muscleGroup.Name))
		}
	})

	t.Run("should update a muscle group", func(t *testing.T) {
		cleanUpDatabase(ctx)

		muscleGroup := createMuscleGroup(ctx, &musclegroup.Model{
			Name: "Chest",
		})

		resp, err := testhelper.RunRequest(
			setup,
			http.MethodPut,
			"/muscle-groups/"+muscleGroup.ID.String(),
			&musclegroup.UpdateMuscleGroupRequest{
				Name: "Chest Updated",
			},
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		responseParsed := testhelper.ParseSuccessResponseBody[musclegroup.Model](resp.Body)
		assert.Equal(t, response.StatusSuccess, responseParsed.Status)
		assert.NotNil(t, responseParsed.Data)
		assert.Equal(t, "Chest Updated", responseParsed.Data.Name)
	})

	t.Run("should delete a muscle group", func(t *testing.T) {
		cleanUpDatabase(ctx)

		muscleGroup := createMuscleGroup(ctx, &musclegroup.Model{
			Name: "Chest",
		})

		resp, err := testhelper.RunRequest(
			setup,
			http.MethodDelete,
			"/muscle-groups/"+muscleGroup.ID.String(),
			nil,
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

func cleanUpDatabase(ctx context.Context) {
	_, err := database.DB().NewDelete().Model((*musclegroup.Model)(nil)).Where("1 = 1").Exec(ctx)
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
