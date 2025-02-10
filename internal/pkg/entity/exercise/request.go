package exercise

import (
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/base"
	"github.com/google/uuid"
)

type (
	ListExercisesQueryParams struct {
		base.ListQueryParams
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
