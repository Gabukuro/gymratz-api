package setup

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Gabukuro/gymratz-api/internal/domain/exercise"
	"github.com/Gabukuro/gymratz-api/internal/domain/user"
	"github.com/Gabukuro/gymratz-api/internal/infra/adapters/postgres"
	"github.com/Gabukuro/gymratz-api/internal/infra/database"
	"github.com/Gabukuro/gymratz-api/internal/pkg/jwt"
	"github.com/Gabukuro/gymratz-api/internal/pkg/middleware"
	"github.com/caarlos0/env/v11"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
)

type (
	EnvVariables struct {
		GoEnv           string `env:"GO_ENV" envDefault:"development"`
		ApplicationName string `env:"APPLICATION_NAME"`
		DatabaseURL     string `env:"DATABASE_URL"`
		JWTSecret       string `env:"JWT_SECRET"`
	}

	Setup struct {
		App *fiber.App
		DB  *bun.DB

		ApplicationName string
		BRLocation      time.Location
		EnvVariables    EnvVariables
		ShutdownChan    chan os.Signal

		shutdownWaitGroup sync.WaitGroup
	}
)

func Init() (*Setup, context.Context) {
	ctx := context.Background()

	var app Setup

	app.configureBRLocation()
	app.configureGracefulShutdown()
	app.configureDevelopmentEnvironment()
	app.configureEnvironmentVariables()

	ctx = app.configureDatabase(ctx)
	app.configureApp()

	return &app, ctx
}

func (s *Setup) configureDatabase(ctx context.Context) context.Context {
	if s.EnvVariables.GoEnv == "test" {
		s.DB, ctx = database.NewTestDB(ctx)
	} else {
		s.DB = database.NewDB(s.EnvVariables.DatabaseURL)
	}

	s.DB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return ctx
}

func (s *Setup) configureApp() {
	s.App = fiber.New(fiber.Config{
		AppName:           s.EnvVariables.ApplicationName,
		EnablePrintRoutes: true,
	})

	s.App.Use(middleware.TraceMiddleware())

	userRepository := postgres.NewUserRepository(s.DB)
	exerciseRepository := postgres.NewExerciseRepository(s.DB)

	tokenService := jwt.NewTokenService(jwt.TokenServiceParams{
		JwtSecret: s.EnvVariables.JWTSecret,
	})

	userService := user.NewService(user.ServiceParams{
		UserRepo:     &userRepository,
		TokenService: tokenService,
	})

	exerciseService := exercise.NewService(exercise.ServiceParams{
		ExerciseRepo: &exerciseRepository,
	})

	user.NewHTTPHandler(user.HTTPHandlerParams{
		App:       s.App,
		Service:   userService,
		JWTSecret: s.EnvVariables.JWTSecret,
	})

	exercise.NewHTTPHandler(exercise.HTTPHandlerParams{
		App:       s.App,
		Service:   exerciseService,
		JWTSecret: s.EnvVariables.JWTSecret,
	})
}

func (s *Setup) configureBRLocation() {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic("Error loading location")
	}

	s.BRLocation = *loc
}

func (s *Setup) configureGracefulShutdown() {
	s.ShutdownChan = make(chan os.Signal, 1)
	signal.Notify(s.ShutdownChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	s.shutdownWaitGroup = sync.WaitGroup{}
	s.shutdownWaitGroup.Add(1)
}

func (s *Setup) WaitShutdown() {
	s.shutdownWaitGroup.Wait()
}

func (s *Setup) Shutdown(ctx context.Context) {
	s.shutdownWaitGroup.Done()
	ctx.Done()

	fmt.Println("Server stopped!")
}

func (s *Setup) configureDevelopmentEnvironment() {
	err := godotenv.Load(".env.local")
	if err != nil && !errors.Is(err, fs.ErrNotExist) && s.EnvVariables.GoEnv == "development" {
		panic("Failed to setup local environment variables")
	}
}

func (s *Setup) configureEnvironmentVariables() {
	if err := env.Parse(&s.EnvVariables); err != nil {
		panic("Error parsing environment variables")
	}
}
