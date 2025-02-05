package repo

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type (
	Repository interface {
		GetDB() *bun.DB
		SetDB(*bun.DB)
		ExecTx(ctx context.Context, txFn func(txCtx context.Context) error) error
	}

	BaseRepository struct {
		db *bun.DB
	}

	contextKey string
)

const txKey = contextKey("tx")

func (r *BaseRepository) SetDB(db *bun.DB) {
	r.db = db
}

func (r *BaseRepository) GetDB() *bun.DB {
	return r.db
}

func (r *BaseRepository) ExecTx(
	ctx context.Context,
	txFn func(txCtx context.Context) error,
) (err error) {
	err = r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		ctxTx := context.WithValue(ctx, txKey, &tx)

		return txFn(ctxTx)
	})

	if err != nil {
		return err
	}

	return nil
}
