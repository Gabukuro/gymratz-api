package repo

import (
	"context"

	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/exercise"
	"github.com/google/uuid"
)

type (
	ExerciseRepository interface {
		Repository

		Create(ctx context.Context, model *exercise.Model) error
		CreateExerciseMuscleGroupAssociations(ctx context.Context, associations []*exercise.ExerciseMuscleGroupModel) error
		GetByID(ctx context.Context, id uuid.UUID) (*exercise.Model, error)
		GetPaginated(ctx context.Context, params exercise.ListExercisesQueryParams) ([]*exercise.Model, error)
		Count(ctx context.Context) (int, error)
		Update(ctx context.Context, id uuid.UUID, model *exercise.Model) error
		UpdateExerciseMuscleGroupAssociations(ctx context.Context, id uuid.UUID, associations []*exercise.ExerciseMuscleGroupModel) error
	}
)
