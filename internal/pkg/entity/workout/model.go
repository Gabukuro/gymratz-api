package workout

import (
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/base"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type (
	Model struct {
		bun.BaseModel `bun:"workouts"`
		base.Model
		UserID uuid.UUID `bun:"user_id"`
		Name   string    `bun:"name"`
	}
)
