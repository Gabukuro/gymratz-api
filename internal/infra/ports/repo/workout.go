package repo

import (
	"context"

	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/workout"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/workoutexercise"
	"github.com/google/uuid"
)

type (
	WorkoutRepository interface {
		Repository

		Create(ctx context.Context, workout *workout.Model) error
		CreateWorkoutExercises(ctx context.Context, exercises []*workoutexercise.Model) error
		GetByID(ctx context.Context, id uuid.UUID) (*workout.Model, error)
		GetPaginated(ctx context.Context, userID uuid.UUID, params workout.ListWorkoutsQueryParams) ([]*workout.Model, int, error)
	}
)
