package postgres

import (
	"context"

	"github.com/Gabukuro/gymratz-api/internal/infra/ports/repo"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/exercise"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type (
	ExerciseRepository struct {
		repo.BaseRepository
	}
)

func NewExerciseRepository(db *bun.DB) ExerciseRepository {
	repo := ExerciseRepository{}
	repo.SetDB(db)

	repo.GetDB().RegisterModel((*exercise.ExerciseMuscleGroupModel)(nil))

	return repo
}

func (r *ExerciseRepository) Create(ctx context.Context, model *exercise.Model) error {
	_, err := r.GetDB().NewInsert().Model(model).Exec(ctx)
	return err
}

func (r *ExerciseRepository) CreateExerciseMuscleGroupAssociations(ctx context.Context, associations []*exercise.ExerciseMuscleGroupModel) error {
	_, err := r.GetDB().NewInsert().Model(&associations).Exec(ctx)
	return err
}

func (r *ExerciseRepository) GetByID(ctx context.Context, id uuid.UUID) (*exercise.Model, error) {
	model := &exercise.Model{}
	err := r.GetDB().NewSelect().Model(model).Relation("MuscleGroups").Where("id = ?", id).Scan(ctx)
	return model, err
}

func (r *ExerciseRepository) GetPaginated(ctx context.Context, limit, offset int) ([]*exercise.Model, error) {
	var models []*exercise.Model
	err := r.GetDB().NewSelect().Model(&models).Limit(limit).Offset(offset).Scan(ctx)
	return models, err
}

func (r *ExerciseRepository) Count(ctx context.Context) (int, error) {
	return r.GetDB().NewSelect().Model((*exercise.Model)(nil)).Count(ctx)
}
