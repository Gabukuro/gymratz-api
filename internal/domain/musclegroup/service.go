package musclegroup

import (
	"context"

	"github.com/Gabukuro/gymratz-api/internal/infra/ports/repo"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/musclegroup"
	"github.com/google/uuid"
)

type (
	Service struct {
		muscleGroupRepo repo.MuscleGroupRepository
	}

	ServiceParams struct {
		MuscleGroupRepo repo.MuscleGroupRepository
	}
)

func NewService(params ServiceParams) *Service {
	return &Service{
		muscleGroupRepo: params.MuscleGroupRepo,
	}
}

func (s *Service) CreateMuscleGroup(ctx context.Context, bodyRequest musclegroup.CreateMuscleGroupRequest) (*musclegroup.Model, error) {
	model := musclegroup.Model{
		Name: bodyRequest.Name,
	}

	if err := s.muscleGroupRepo.Create(ctx, &model); err != nil {
		return nil, err
	}

	return &model, nil
}

func (s *Service) ListMuscleGroups(ctx context.Context, params musclegroup.ListMuscleGroupsQueryParams) ([]*musclegroup.Model, int, error) {
	muscleGroups, err := s.muscleGroupRepo.GetPaginated(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.muscleGroupRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return muscleGroups, total, nil
}

func (s *Service) UpdateMuscleGroup(ctx context.Context, id uuid.UUID, bodyRequest musclegroup.UpdateMuscleGroupRequest) (*musclegroup.Model, error) {
	model := musclegroup.Model{
		Name: bodyRequest.Name,
	}

	if err := s.muscleGroupRepo.Update(ctx, id, &model); err != nil {
		return nil, err
	}

	return &model, nil
}

func (s *Service) DeleteMuscleGroup(ctx context.Context, id uuid.UUID) error {
	return s.muscleGroupRepo.Delete(ctx, id)
}
