package exercise

import "github.com/google/uuid"

type (
	ListExercisesQueryParams struct {
		Page             int      `query:"page"`
		PerPage          int      `query:"per_page"`
		Name             string   `query:"name"`
		MuscleGroupNames []string `query:"muscle_group_names"`
	}

	CreateExerciseRequest struct {
		Name           string      `json:"name"`
		Description    string      `json:"description"`
		MuscleGroupIDs []uuid.UUID `json:"muscle_group_ids"`
	}

	UpdateExerciseRequest struct {
		Name           string      `json:"name"`
		Description    string      `json:"description"`
		MuscleGroupIDs []uuid.UUID `json:"muscle_group_ids"`
	}
)
