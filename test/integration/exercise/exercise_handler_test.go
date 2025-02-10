package exercise__test

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"slices"
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
			"/exercises",
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

		// Create a new exercise wit a muscle group
		testExercise, testMuscleGroup := createExerciseWithMuscleGroup(ctx,
			"pushup",
			"pushup description",
			"chest",
		)

		// Send a request to list exercises
		resp, err := testhelper.RunRequest(setup,
			http.MethodGet,
			"/exercises?page=1&per_page=10&name=pushup&muscle_group_names=chest&muscle_group_names=triceps",
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

		// Check the response data
		assert.Equal(t, 1, len(responseParsed.Data))
		assert.Equal(t, testExercise.ID, responseParsed.Data[0].ID)
		assert.Equal(t, testExercise.Name, responseParsed.Data[0].Name)
		assert.Equal(t, testExercise.Description, responseParsed.Data[0].Description)

		// Check the muscle group
		assert.Equal(t, 1, len(responseParsed.Data[0].MuscleGroups))
		assert.Equal(t, testMuscleGroup.Name, responseParsed.Data[0].MuscleGroups[0].Name)
	})

	t.Run("should update an exercise", func(t *testing.T) {
		// Clean up the database
		cleanUpDatabase(ctx)

		// Create a new exercise wit a muscle group
		testExercise, testMuscleGroup := createExerciseWithMuscleGroup(ctx,
			"pushup",
			"pushup description",
			"chest",
		)

		// Create a new muscle group
		newMuscleGroup := createMuscleGroup(ctx, &musclegroup.Model{
			Name: "triceps",
		})

		// Send a request to update the exercise
		updateExerciseRequest := exerciseEntity.UpdateExerciseRequest{
			Name:           "pushup updated",
			Description:    "pushup description updated",
			MuscleGroupIDs: []uuid.UUID{testMuscleGroup.ID, newMuscleGroup.ID},
		}

		rep, err := testhelper.RunRequest(setup,
			http.MethodPut,
			"/exercises/"+testExercise.ID.String(),
			updateExerciseRequest,
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rep.StatusCode)

		responseParsed := testhelper.ParseSuccessResponseBody[exerciseEntity.Model](rep.Body)
		assert.Equal(t, response.StatusSuccess, responseParsed.Status)
		assert.Equal(t, updateExerciseRequest.Name, responseParsed.Data.Name)
		assert.Equal(t, updateExerciseRequest.Description, responseParsed.Data.Description)

		// Check the muscle groups
		assert.Equal(t, 2, len(responseParsed.Data.MuscleGroups))
		muscleGroupsIDS := []uuid.UUID{testMuscleGroup.ID, newMuscleGroup.ID}
		for _, muscleGroup := range responseParsed.Data.MuscleGroups {
			assert.True(t, slices.Contains(muscleGroupsIDS, muscleGroup.ID))
		}
	})

	t.Run("should delete an exercise", func(t *testing.T) {
		// Clean up the database
		cleanUpDatabase(ctx)

		// Create a new exercise wit a muscle group
		testExercise, _ := createExerciseWithMuscleGroup(ctx,
			"pushup",
			"pushup description",
			"chest",
		)

		// Send a request to delete the exercise
		rep, err := testhelper.RunRequest(setup,
			http.MethodDelete,
			"/exercises/"+testExercise.ID.String(),
			nil,
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusNoContent, rep.StatusCode)

		// Check if the exercise was deleted
		exercise := getExerciseByID(ctx, testExercise.ID)
		assert.Nil(t, exercise)
	})
}

func cleanUpDatabase(ctx context.Context) {
	dropExercises(ctx)
	dropMuscleGroups(ctx)
}

func createExerciseWithMuscleGroup(ctx context.Context,
	name string,
	description string,
	muscleGroupName string,
) (exerciseEntity.Model, musclegroup.Model) {
	muscleGroup := createMuscleGroup(ctx, &musclegroup.Model{
		Name: muscleGroupName,
	})

	exerciseEntity := createExercise(ctx, &exerciseEntity.Model{
		Name:        name,
		Description: description,
	})

	createExerciseMuscleGroupAssociation(ctx, exerciseEntity, muscleGroup)

	return exerciseEntity, muscleGroup
}

func createExerciseMuscleGroupAssociation(ctx context.Context, exercise exerciseEntity.Model, muscleGroup musclegroup.Model) {
	_, err := database.DB().NewInsert().Model(&exerciseEntity.ExerciseMuscleGroupModel{
		ExerciseID:    exercise.ID,
		MuscleGroupID: muscleGroup.ID,
	}).Exec(ctx)
	if err != nil {
		panic(err)
	}
}

func createExercise(ctx context.Context, model *exerciseEntity.Model) exerciseEntity.Model {
	_, err := database.DB().NewInsert().Model(model).Exec(ctx)
	if err != nil {
		panic(err)
	}

	return *model
}

func getExerciseByID(ctx context.Context, id uuid.UUID) *exerciseEntity.Model {
	model := &exerciseEntity.Model{}
	err := database.DB().NewSelect().Model(model).Where("id = ?", id).Scan(ctx)
	if err != nil && err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}

	return model
}

func dropExercises(ctx context.Context) {
	_, err := database.DB().NewDelete().Model(&exerciseEntity.Model{}).Where("1 = 1").Exec(ctx)
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

func dropMuscleGroups(ctx context.Context) {
	_, err := database.DB().NewDelete().Model(&musclegroup.Model{}).Where("1 = 1").Exec(ctx)
	if err != nil {
		panic(err)
	}
}
