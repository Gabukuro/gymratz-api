package database

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var db *bun.DB

func NewDB(dsn string) *bun.DB {
	if dsn == "" {
		panic("database dsn is empty")
	}

	if db != nil {
		return db
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("could not connect to database: %w", err))
	}

	return db
}

func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			fmt.Println("could not close database: %w", err)
		}
	}
}
