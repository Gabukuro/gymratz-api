package main

import (
	"fmt"

	"github.com/Gabukuro/gymratz-api/internal/domain/user"
	"github.com/Gabukuro/gymratz-api/internal/pkg/setup"
	"github.com/gofiber/fiber/v2"
)

func main() {
	setup := setup.Init(setup.SetupConfig{
		ApplicationName: "gymratz-api",
	})

	app := fiber.New(fiber.Config{
		AppName:           setup.ApplicationName,
		EnablePrintRoutes: true,
	})

	user.NewHTTPHandler(user.HTTPHandlerParams{
		App: app,
	})

	go func() {
		<-setup.ShutdownChan

		fmt.Println("Shutting down server...")
		app.Shutdown()

		setup.Shutdown()
	}()

	if err := app.Listen(":3000"); err != nil {
		fmt.Println("Error starting server:", err)
	}

	setup.WaitShutdown()
}
