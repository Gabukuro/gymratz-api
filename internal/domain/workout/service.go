package workout

import (
	"context"

	"github.com/Gabukuro/gymratz-api/internal/infra/ports/repo"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/workout"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/workoutexercise"
	"github.com/google/uuid"
)

type (
	Service struct {
		workoutRepo repo.WorkoutRepository
		userRepo    repo.UserRepository
	}

	ServiceParams struct {
		WorkoutRepo repo.WorkoutRepository
		UserRepo    repo.UserRepository
	}
)

func NewService(params ServiceParams) *Service {
	return &Service{
		workoutRepo: params.WorkoutRepo,
		userRepo:    params.UserRepo,
	}
}

func (s *Service) CreateWorkout(ctx context.Context, params workout.CreateWorkoutRequest) (*workout.Model, error) {
	workoutModel := workout.Model{
		UserID: params.UserID,
		Name:   params.Name,
	}

	err := s.workoutRepo.ExecTx(ctx, func(txCtx context.Context) error {
		err := s.workoutRepo.Create(txCtx, &workoutModel)
		if err != nil {
			return err
		}

		return s.createWorkoutExercises(txCtx, workoutModel.ID, params)
	})

	if err != nil {
		return nil, err
	}

	return s.workoutRepo.GetByIDWithRelations(ctx, workoutModel.ID)

}

func (s *Service) createWorkoutExercises(ctx context.Context, workoutID uuid.UUID, params workout.CreateWorkoutRequest) error {
	workoutExercises := make([]*workoutexercise.Model, len(params.Exercises))

	for index, workoutExercise := range params.Exercises {
		workoutExercises[index] = &workoutexercise.Model{
			WorkoutID:   workoutID,
			ExerciseID:  workoutExercise.ExerciseID,
			Sets:        workoutExercise.Sets,
			Repetitions: workoutExercise.Repetitions,
			Weight:      workoutExercise.Weight,
			Duration:    workoutExercise.Duration,
			RestTime:    workoutExercise.RestTime,
			Notes:       workoutExercise.Notes,
		}
	}

	return s.workoutRepo.CreateWorkoutExercises(ctx, workoutExercises)
}

func (s *Service) ListUserWorkouts(ctx context.Context, userEmail string, params workout.ListWorkoutsQueryParams) ([]*workout.Model, int, error) {
	userModel, err := s.userRepo.FindByEmail(ctx, userEmail)
	if err != nil {
		return nil, 0, err
	}

	return s.workoutRepo.GetPaginated(ctx, userModel.ID, params)
}

func (s *Service) UpdateWorkout(ctx context.Context, workoutID uuid.UUID, params workout.UpdateWorkoutRequest) (*workout.Model, error) {
	workoutModel, err := s.workoutRepo.GetByID(ctx, workoutID)
	if err != nil {
		return nil, err
	}

	workoutModel.Name = params.Name

	err = s.workoutRepo.UpdateWorkout(ctx, workoutID, workoutModel)
	if err != nil {
		return nil, err
	}

	return s.workoutRepo.GetByIDWithRelations(ctx, workoutID)
}

func (s *Service) UpdateWorkoutExercise(ctx context.Context, workoutID uuid.UUID, workoutExerciseID uuid.UUID, params workout.UpdateWorkoutExerciseRequest) (*workoutexercise.Model, error) {
	workoutExercise, err := s.workoutRepo.GetWorkoutExercise(ctx, workoutID, workoutExerciseID)
	if err != nil {
		return nil, err
	}

	workoutExercise.Sets = params.Sets
	workoutExercise.Repetitions = params.Repetitions
	workoutExercise.Weight = params.Weight
	workoutExercise.Duration = params.Duration
	workoutExercise.RestTime = params.RestTime
	workoutExercise.Notes = params.Notes

	err = s.workoutRepo.UpdateWorkoutExercise(ctx, workoutID, workoutExerciseID, workoutExercise)
	if err != nil {
		return nil, err
	}

	return s.workoutRepo.GetWorkoutExercise(ctx, workoutID, workoutExerciseID)
}
