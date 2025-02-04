package msuclegroup

import (
	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/base"
	"github.com/uptrace/bun"
)

type (
	Model struct {
		bun.BaseModel `bun:"muscle_groups"`
		base.Model
		Name string `bun:"name"`
	}
)
