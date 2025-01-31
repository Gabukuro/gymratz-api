package repo

import (
	"context"

	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/user"
)

type (
	UserRepository interface {
		Create(ctx context.Context, model user.Model) error
		FindByEmail(ctx context.Context, email string) (*user.Model, error)
	}
)
