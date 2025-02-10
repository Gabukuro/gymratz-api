package workout_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/Gabukuro/gymratz-api/internal/infra/database"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/user"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/workout"
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

		user := testhelper.CreateUser(ctx, database.DB(), &user.Model{
			Name:     "John Doe",
			Email:    "john@doe.com",
			Password: "password",
		})

		exercise, _ := testhelper.CreateExerciseWithMuscleGroup(ctx, database.DB(), "Barbell Bench Press", "Barbell Bench Press Description", "Chest")

		rep, err := testhelper.RunRequest(
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
		assert.Equal(t, http.StatusCreated, rep.StatusCode)
	})
}
