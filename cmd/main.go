package main

import (
	"fmt"

	"github.com/Gabukuro/gymratz-api/internal/domain/user"
	"github.com/Gabukuro/gymratz-api/internal/infra/adapters/postgres"
	"github.com/Gabukuro/gymratz-api/internal/infra/database"
	"github.com/Gabukuro/gymratz-api/internal/pkg/jwt"
	"github.com/Gabukuro/gymratz-api/internal/pkg/setup"
	"github.com/gofiber/fiber/v2"
)

func main() {
	setup := setup.Init()

	app := fiber.New(fiber.Config{
		AppName:           setup.ApplicationName,
		EnablePrintRoutes: true,
	})

	db := database.NewDB(setup.EnvVariables.DatabaseURL)
	userRepository := postgres.NewUserRepository(db)

	tokenService := jwt.NewTokenService(jwt.TokenServiceParams{
		JwtSecret: setup.EnvVariables.JWTSecret,
	})

	userService := user.NewService(user.ServiceParams{
		UserRepo:     &userRepository,
		TokenService: tokenService,
	})

	user.NewHTTPHandler(user.HTTPHandlerParams{
		App:     app,
		Service: userService,
	})

	go func() {
		<-setup.ShutdownChan

		database.CloseDB()

		fmt.Println("Shutting down server...")
		app.Shutdown()

		setup.Shutdown()
	}()

	if err := app.Listen(":3000"); err != nil {
		fmt.Println("Error starting server:", err)
	}

	setup.WaitShutdown()
}
