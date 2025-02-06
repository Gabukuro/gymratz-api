package exercise

import (
	"context"

	"github.com/Gabukuro/gymratz-api/internal/infra/ports/repo"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/exercise"
	"github.com/google/uuid"
)

type (
	Service struct {
		exerciseRepo repo.ExerciseRepository
	}

	ServiceParams struct {
		ExerciseRepo repo.ExerciseRepository
	}
)

func NewService(params ServiceParams) *Service {
	return &Service{
		exerciseRepo: params.ExerciseRepo,
	}
}

func (s *Service) CreateExercise(ctx context.Context, params exercise.CreateExerciseRequest) (*exercise.Model, error) {
	exerciseModel := exercise.Model{
		Name:        params.Name,
		Description: params.Description,
	}

	err := s.exerciseRepo.ExecTx(ctx, func(txCtx context.Context) error {
		err := s.exerciseRepo.Create(txCtx, &exerciseModel)
		if err != nil {
			return err
		}

		return s.createExerciseAssociations(txCtx, exerciseModel.ID, params)
	})

	if err != nil {
		return nil, err
	}

	return s.exerciseRepo.GetByID(ctx, exerciseModel.ID)
}

func (s *Service) createExerciseAssociations(
	ctx context.Context,
	exerciseID uuid.UUID,
	params exercise.CreateExerciseRequest,
) error {
	associations := make([]*exercise.ExerciseMuscleGroupModel, len(params.MuscleGroupIDs))

	for index, muscleGroup := range params.MuscleGroupIDs {
		associations[index] = &exercise.ExerciseMuscleGroupModel{
			ExerciseID:    exerciseID,
			MuscleGroupID: muscleGroup,
		}
	}

	return s.exerciseRepo.CreateExerciseMuscleGroupAssociations(ctx, associations)
}

func (s *Service) ListExercises(ctx context.Context, params exercise.ListExercisesQueryParams) ([]*exercise.Model, int, error) {
	exercises, err := s.exerciseRepo.GetPaginated(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.exerciseRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return exercises, total, nil
}

func (s *Service) UpdateExercise(ctx context.Context, id uuid.UUID, exerciseEdited exercise.UpdateExerciseRequest) (*exercise.Model, error) {
	exerciseModel, err := s.exerciseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	exerciseModel.Name = exerciseEdited.Name
	exerciseModel.Description = exerciseEdited.Description

	err = s.exerciseRepo.ExecTx(ctx, func(txCtx context.Context) error {
		err := s.exerciseRepo.Update(txCtx, id, exerciseModel)
		if err != nil {
			return err
		}

		return s.updateExerciseAssociations(txCtx, id, exerciseEdited)
	})

	if err != nil {
		return nil, err
	}

	return s.exerciseRepo.GetByID(ctx, exerciseModel.ID)
}

func (s *Service) updateExerciseAssociations(ctx context.Context, exerciseID uuid.UUID, exerciseEdited exercise.UpdateExerciseRequest) error {
	associations := make([]*exercise.ExerciseMuscleGroupModel, len(exerciseEdited.MuscleGroupIDs))

	for index, muscleGroup := range exerciseEdited.MuscleGroupIDs {
		associations[index] = &exercise.ExerciseMuscleGroupModel{
			ExerciseID:    exerciseID,
			MuscleGroupID: muscleGroup,
		}
	}

	return s.exerciseRepo.UpdateExerciseMuscleGroupAssociations(ctx, exerciseID, associations)
}
