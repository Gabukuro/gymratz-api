package repo

import (
	"context"

	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/musclegroup"
	"github.com/google/uuid"
)

type (
	MuscleGroupRepository interface {
		Repository

		Create(ctx context.Context, model *musclegroup.Model) error
		GetPaginated(ctx context.Context, params musclegroup.ListMuscleGroupsQueryParams) ([]*musclegroup.Model, error)
		Count(ctx context.Context) (int, error)
		Update(ctx context.Context, id uuid.UUID, model *musclegroup.Model) error
		Delete(ctx context.Context, id uuid.UUID) error
	}
)
