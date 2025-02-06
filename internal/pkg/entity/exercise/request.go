package exercise

import "github.com/google/uuid"

type (
	ListExercisesQueryParams struct {
		Page             int      `json:"page"`
		PerPage          int      `json:"per_page"`
		Name             string   `json:"name"`
		MuscleGroupNames []string `json:"muscle_group_names"`
	}

	CreateExerciseRequest struct {
		Name           string      `json:"name"`
		Description    string      `json:"description"`
		MuscleGroupIDs []uuid.UUID `json:"muscle_group_ids"`
	}
)
