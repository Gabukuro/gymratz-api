package exercise

import "github.com/uptrace/bun"

type (
	ExerciseMuscleGroupModel struct {
		bun.BaseModel `bun:"exercise_muscle_groups"`
		ExerciseID    int `bun:"exercise_id"`
		MuscleGroupID int `bun:"muscle_group_id,pk"`
	}
)
