package postgres

import (
	"context"

	"github.com/Gabukuro/gymratz-api/internal/pkg/entity/user"
	"github.com/uptrace/bun"
)

type (
	UserRepository struct {
		db *bun.DB
	}
)

func NewUserRepository(db *bun.DB) UserRepository {
	return UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, model user.Model) error {
	_, err := r.db.NewInsert().Model(&model).Exec(ctx)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (user.Model, error) {
	var model user.Model
	err := r.db.NewSelect().Model(&model).Where("email = ?", email).Scan(ctx)
	return model, err
}
