package workoutexercise

import (
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/base"
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/exercise"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type (
	Model struct {
		bun.BaseModel `bun:"workout_exercises"`
		base.Model
		WorkoutID   uuid.UUID       `bun:"workout_id"`
		ExerciseID  uuid.UUID       `bun:"exercise_id"`
		Sets        int             `bun:"sets"`
		Repetitions *int            `bun:"repetitions"`
		Weight      *float64        `bun:"weight"`
		Duration    *int            `bun:"duration"`
		RestTime    int             `bun:"rest_time"`
		Notes       *string         `bun:"notes"`
		Exercise    *exercise.Model `bun:"rel:belongs-to,join:exercise_id=id"`
	}
)
