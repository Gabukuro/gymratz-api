package postgres

import (
	"context"

	"github.com/Gabukuro/gymratz-api/internal/infra/ports/repo"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/workout"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/workoutexercise"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type (
	WorkoutRepository struct {
		repo.BaseRepository
	}
)

func NewWorkoutRepository(db *bun.DB) WorkoutRepository {
	repo := WorkoutRepository{}
	repo.SetDB(db)

	return repo
}

func (r *WorkoutRepository) Create(ctx context.Context, model *workout.Model) error {
	_, err := r.GetDB().NewInsert().Model(model).Exec(ctx)
	return err
}

func (r *WorkoutRepository) CreateWorkoutExercises(ctx context.Context, exercises []*workoutexercise.Model) error {
	_, err := r.GetDB().NewInsert().Model(&exercises).Exec(ctx)
	return err
}

func (r *WorkoutRepository) GetByID(ctx context.Context, id uuid.UUID) (*workout.Model, error) {
	model := &workout.Model{}
	err := r.GetDB().NewSelect().Model(model).Relation("WorkoutExercises.Exercise").Where("id = ?", id).Scan(ctx)
	return model, err
}

func (r *WorkoutRepository) GetPaginated(ctx context.Context, userID uuid.UUID, params workout.ListWorkoutsQueryParams) ([]*workout.Model, int, error) {
	var models []*workout.Model
	limit := params.PerPage
	offset := (params.Page - 1) * params.PerPage

	query := r.GetDB().NewSelect().
		Model(&models).
		Relation("WorkoutExercises.Exercise").
		Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset)

	err := query.Scan(ctx)
	if err != nil {
		return nil, 0, nil
	}

	total, err := query.Count(ctx)

	return models, total, err
}
