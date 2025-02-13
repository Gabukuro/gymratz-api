package workout

import (
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/base"
	"github.com/google/uuid"
)

type (
	CreateWorkoutRequest struct {
		Name      string            `json:"name"`
		UserID    uuid.UUID         `json:"user_id"`
		Exercises []WorkoutExercise `json:"exercises"`
	}

	WorkoutExercise struct {
		ExerciseID  uuid.UUID `json:"exercise_id"`
		Sets        int       `json:"sets"`
		Repetitions *int      `json:"repetitions"`
		Weight      *float64  `json:"weight"`
		Duration    *int      `json:"duration"`
		RestTime    int       `json:"rest_time"`
		Notes       *string   `json:"notes"`
	}

	ListWorkoutsQueryParams struct {
		base.ListQueryParams
		Name string `json:"name"`
	}
)
