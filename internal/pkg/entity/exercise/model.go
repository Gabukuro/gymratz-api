package exercise

import (
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/base"
	"github.com/uptrace/bun"
)

type (
	Model struct {
		bun.BaseModel `bun:"exercises"`
		base.Model
		Name        string `bun:"name"`
		Description string `bun:"description"`
	}
)
