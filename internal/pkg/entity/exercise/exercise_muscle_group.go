package exercise

import (
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/musclegroup"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type (
	ExerciseMuscleGroupModel struct {
		bun.BaseModel `bun:"exercise_muscle_groups"`
		ExerciseID    uuid.UUID          `bun:"exercise_id,pk"`
		Exercise      *Model             `bun:"rel:belongs-to,join:exercise_id=id"`
		MuscleGroupID uuid.UUID          `bun:"muscle_group_id,pk"`
		MuscleGroup   *musclegroup.Model `bun:"rel:belongs-to,join:muscle_group_id=id"`
	}
)
