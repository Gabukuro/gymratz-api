package exercise

import (
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/base"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/musclegroup"
	"github.com/uptrace/bun"
)

type (
	Model struct {
		bun.BaseModel `bun:"table:exercises"`
		base.Model
		Name         string              `bun:"name"`
		Description  string              `bun:"description"`
		MuscleGroups []musclegroup.Model `bun:"m2m:exercise_muscle_groups,join:Exercise=MuscleGroup"`
	}
)
