package postgres

import (
	"context"

	"github.com/Gabukuro/gymratz-api/internal/infra/ports/repo"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/musclegroup"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type (
	MuscleGroupRepository struct {
		repo.BaseRepository
	}
)

func NewMuscleGroupRepository(db *bun.DB) MuscleGroupRepository {
	repo := MuscleGroupRepository{}
	repo.SetDB(db)

	return repo
}

func (r *MuscleGroupRepository) Create(ctx context.Context, model *musclegroup.Model) error {
	_, err := r.GetDB().NewInsert().Model(model).Exec(ctx)
	return err
}

func (r *MuscleGroupRepository) GetPaginated(ctx context.Context, params musclegroup.ListMuscleGroupsQueryParams) ([]*musclegroup.Model, error) {
	var models []*musclegroup.Model
	limit := params.PerPage
	offset := (params.Page - 1) * params.PerPage

	query := r.GetDB().NewSelect().
		Model(&models).
		Limit(limit).
		Offset(offset)

	if params.Name != "" {
		query.Where("name ILIKE ?", "%"+params.Name+"%")
	}

	err := query.Scan(ctx)

	return models, err
}

func (r *MuscleGroupRepository) Count(ctx context.Context) (int, error) {
	return r.GetDB().NewSelect().Model(&musclegroup.Model{}).Count(ctx)
}

func (r *MuscleGroupRepository) Update(ctx context.Context, id uuid.UUID, model *musclegroup.Model) error {
	_, err := r.GetDB().NewUpdate().Model(model).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *MuscleGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.GetDB().NewDelete().Model(&musclegroup.Model{}).Where("id = ?", id).Exec(ctx)
	return err
}
