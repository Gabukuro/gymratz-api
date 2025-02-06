package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	db     *bun.DB
	testDB *bun.DB
)

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

func NewTestDB(ctx context.Context) (*bun.DB, context.Context) {
	// TODO: move this to a test helper package
	if testDB != nil {
		return testDB, ctx
	}

	db, ctx := setupTestContainers(ctx)
	time.Sleep(2 * time.Second)
	runMigrations(db)

	testDB = bun.NewDB(db, pgdialect.New())
	if err := testDB.Ping(); err != nil {
		panic(fmt.Errorf("could not connect to database: %w", err))
	}

	return testDB, ctx
}

func CloseTestDB(ctx context.Context) {
	if testDB != nil {
		if err := testDB.Close(); err != nil {
			fmt.Printf("could not close test database: %v\n", err)
		}

		postgresContainer := ctx.Value("postgresContainer").(testcontainers.Container)
		postgresContainer.Terminate(ctx)
	}
}

func DB() *bun.DB {
	if testDB != nil {
		return testDB
	}

	return db
}

func setupTestContainers(ctx context.Context) (*sql.DB, context.Context) {
	// TODO: move this to a test helper
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "gymratz-api-test",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	ctx = context.WithValue(ctx, "postgresContainer", postgresContainer)

	if err != nil {
		panic(fmt.Errorf("failed to start PostgreSQL container: %v", err))
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		panic(fmt.Errorf("failed to get container port: %v", err))
	}

	dsn := fmt.Sprintf("postgres://user:password@localhost:%s/gymratz-api-test?sslmode=disable", port.Port())
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("failed to connect to PostgreSQL: %v", err))
	}

	return db, ctx
}

func runMigrations(db *sql.DB) {
	rootDir, err := getRootDirectory()
	if err != nil {
		panic(fmt.Errorf("could not get root directory: %w", err))
	}

	migrations := &migrate.FileMigrationSource{
		Dir: rootDir + "migrations",
	}

	migrate.SetTable("migrations")
	migrationsCount, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		panic(fmt.Errorf("could not apply migrations: %w", err))
	}

	fmt.Printf("applied %d migrations", migrationsCount)
}

func getRootDirectory() (string, error) {
	dir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	dirs := strings.Split(dir, "/")
	rootDir := ""

	for _, dir := range dirs {
		rootDir += dir + "/"

		info, _ := os.Stat(rootDir + "go.mod")

		if info != nil {
			if info.Name() == "go.mod" {
				break
			}
		}
	}

	return rootDir, nil
}
