package exercise

import "github.com/google/uuid"

type (
	ListExercisesRequest struct {
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
	}

	CreateExerciseRequest struct {
		Name           string      `json:"name"`
		Description    string      `json:"description"`
		MuscleGroupIDs []uuid.UUID `json:"muscle_group_ids"`
	}
)
