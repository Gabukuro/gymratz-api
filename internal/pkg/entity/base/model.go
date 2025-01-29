package base

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type (
	Model struct {
		ID        uuid.UUID  `json:"id" bun:"id,pk"`
		CreatedAt time.Time  `json:"created_at" bun:"created_at"`
		UpdatedAt time.Time  `json:"updated_at" bun:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	}
)

func (m *Model) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.ID = uuid.New()
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
