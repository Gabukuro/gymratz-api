package workouthistory

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type (
	Model struct {
		bun.BaseModel `bun:"workout_history"`
		ID            uuid.UUID `bun:"id,pk"`
		UserID        uuid.UUID `bun:"user_id"`
		WorkoutID     uuid.UUID `bun:"workout_id"`
		CompletedAt   time.Time `bun:"completed_at"`
		Notes         *string   `bun:"notes"`
	}
)

func (m *Model) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.ID = uuid.New()
		m.CompletedAt = time.Now()
	}
	return nil
}
