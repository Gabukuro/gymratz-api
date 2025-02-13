package workout_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/Gabukuro/gymratz-api/internal/infra/database"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/workout"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/workoutexercise"
	"github.com/Gabukuro/gymratz-api/internal/pkg/response"
	"github.com/Gabukuro/gymratz-api/internal/pkg/setup"
	"github.com/Gabukuro/gymratz-api/internal/pkg/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestWorkoutHandler(t *testing.T) {
	t.Parallel()

	os.Setenv("GO_ENV", "test")
	setup, ctx := setup.Init()
	defer database.CloseTestDB(ctx)

	t.Run("should create a new workout", func(t *testing.T) {
		testhelper.CleanUpDatabase(ctx, database.DB())

		user := testhelper.CreateUser(ctx, database.DB(), nil)

		exercise, _ := testhelper.CreateExerciseWithMuscleGroup(ctx, database.DB(), "Barbell Bench Press", "Barbell Bench Press Description", "Chest")
		resp, err := testhelper.RunRequest(
			setup,
			http.MethodPost,
			"/workouts",
			&workout.CreateWorkoutRequest{
				Name:   "Chest Day",
				UserID: user.ID,
				Exercises: []workout.WorkoutExercise{
					{
						ExerciseID:  exercise.ID,
						Sets:        3,
						Repetitions: testhelper.GetPointer(10),
						Weight:      testhelper.GetPointer(float64(20)),
						Duration:    testhelper.GetPointer(0),
						RestTime:    60,
						Notes:       testhelper.GetPointer("This is a note"),
					},
				},
			},
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		responseParsed := testhelper.ParseSuccessResponseBody[workout.Model](resp.Body)
		assert.Equal(t, response.StatusSuccess, responseParsed.Status)

		assert.Equal(t, "Chest Day", responseParsed.Data.Name)
		assert.Equal(t, user.ID, responseParsed.Data.UserID)
		assert.Equal(t, 1, len(responseParsed.Data.WorkoutExercises))
		assert.Equal(t, exercise.ID, responseParsed.Data.WorkoutExercises[0].ExerciseID)
		assert.Equal(t, 3, responseParsed.Data.WorkoutExercises[0].Sets)
		assert.Equal(t, 10, *responseParsed.Data.WorkoutExercises[0].Repetitions)
		assert.Equal(t, 20.0, *responseParsed.Data.WorkoutExercises[0].Weight)
		assert.Equal(t, 0, *responseParsed.Data.WorkoutExercises[0].Duration)
		assert.Equal(t, 60, responseParsed.Data.WorkoutExercises[0].RestTime)
		assert.Equal(t, "This is a note", *responseParsed.Data.WorkoutExercises[0].Notes)
	})

	t.Run("should list all user workouts", func(t *testing.T) {
		testhelper.CleanUpDatabase(ctx, database.DB())

		user := testhelper.CreateUser(ctx, database.DB(), nil)
		testWorkouts := testhelper.CreateUserManyWorkout(ctx, database.DB(), user.ID, 5)

		resp, err := testhelper.RunRequest(
			setup,
			http.MethodGet,
			"/workouts",
			nil,
			map[string]string{
				"Authorization": testhelper.GenerateAuthToken(setup.EnvVariables.JWTSecret, &user.Email),
			},
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		responseParsed := testhelper.ParsePaginationResponseBody[[]workout.Model](resp.Body)
		assert.Equal(t, response.StatusSuccess, responseParsed.Status)

		assert.Len(t, responseParsed.Data, len(testWorkouts))
		for _, item := range responseParsed.Data {
			testWorkout, ok := testWorkouts[item.ID.String()]
			assert.True(t, ok)

			assert.Equal(t, testWorkout.ID, item.ID)
			assert.Equal(t, testWorkout.Name, item.Name)
			assert.Equal(t, testWorkout.UserID, item.UserID)

			assert.Len(t, item.WorkoutExercises, len(testWorkout.WorkoutExercises))
			for i, exercise := range item.WorkoutExercises {
				assert.Equal(t, testWorkout.WorkoutExercises[i].ExerciseID, exercise.ExerciseID)
				assert.Equal(t, testWorkout.WorkoutExercises[i].Sets, exercise.Sets)
				assert.Equal(t, testWorkout.WorkoutExercises[i].Repetitions, exercise.Repetitions)
				assert.Equal(t, testWorkout.WorkoutExercises[i].Weight, exercise.Weight)
				assert.Equal(t, testWorkout.WorkoutExercises[i].Duration, exercise.Duration)
				assert.Equal(t, testWorkout.WorkoutExercises[i].RestTime, exercise.RestTime)
				assert.Equal(t, testWorkout.WorkoutExercises[i].Notes, exercise.Notes)
			}
		}
	})

	t.Run("should update a workout", func(t *testing.T) {
		testhelper.CleanUpDatabase(ctx, database.DB())

		user := testhelper.CreateUser(ctx, database.DB(), nil)
		testWorkout := testhelper.CreateSingleWorkout(ctx, database.DB(), user.ID, 1)

		resp, err := testhelper.RunRequest(
			setup,
			http.MethodPut,
			"/workouts/"+testWorkout.ID.String(),
			&workout.UpdateWorkoutRequest{
				Name: "Updated workout name",
			},
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		responseParsed := testhelper.ParseSuccessResponseBody[workout.Model](resp.Body)
		assert.Equal(t, response.StatusSuccess, responseParsed.Status)

		assert.Equal(t, "Updated workout name", responseParsed.Data.Name)
		assert.Equal(t, user.ID, responseParsed.Data.UserID)
		assert.Equal(t, testWorkout.ID, responseParsed.Data.ID)
	})

	t.Run("should update a workout exercise", func(t *testing.T) {
		testhelper.CleanUpDatabase(ctx, database.DB())

		user := testhelper.CreateUser(ctx, database.DB(), nil)
		testWorkout := testhelper.CreateSingleWorkout(ctx, database.DB(), user.ID, 1)
		workoutExercise := testWorkout.WorkoutExercises[0]

		updateWorkoutExerciseRequest := workout.UpdateWorkoutExerciseRequest{
			Sets:        5,
			Repetitions: testhelper.GetPointer(15),
			Weight:      testhelper.GetPointer(float64(30)),
			Duration:    testhelper.GetPointer(0),
			RestTime:    90,
			Notes:       testhelper.GetPointer("This is an updated note"),
		}

		resp, err := testhelper.RunRequest(
			setup,
			http.MethodPut,
			"/workouts/"+testWorkout.ID.String()+"/exercises/"+workoutExercise.ID.String(),
			&updateWorkoutExerciseRequest,
			nil,
		)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		responseParsed := testhelper.ParseSuccessResponseBody[workoutexercise.Model](resp.Body)
		assert.Equal(t, response.StatusSuccess, responseParsed.Status)

		assert.Equal(t, workoutExercise.ID, responseParsed.Data.ID)
		assert.Equal(t, testWorkout.ID, responseParsed.Data.WorkoutID)
		assert.Equal(t, workoutExercise.ExerciseID, responseParsed.Data.ExerciseID)
		assert.Equal(t, updateWorkoutExerciseRequest.Sets, responseParsed.Data.Sets)
		assert.Equal(t, *updateWorkoutExerciseRequest.Repetitions, *responseParsed.Data.Repetitions)
		assert.Equal(t, *updateWorkoutExerciseRequest.Weight, *responseParsed.Data.Weight)
		assert.Equal(t, *updateWorkoutExerciseRequest.Duration, *responseParsed.Data.Duration)
		assert.Equal(t, updateWorkoutExerciseRequest.RestTime, responseParsed.Data.RestTime)
		assert.Equal(t, *updateWorkoutExerciseRequest.Notes, *responseParsed.Data.Notes)
	})
}
