package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gabukuro/gymratz-api/internal/domain/user"
	"github.com/gofiber/fiber/v2"
)

func main() {
	println("Starting server...")

	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		fmt.Println("Error loading location", err)
	}

	time.Local = loc

	app := fiber.New(fiber.Config{
		AppName:           "gymratz-api",
		EnablePrintRoutes: true,
	})

	user.NewHTTPHandler(user.HTTPHandlerParams{
		App: app,
	})

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":3000"); err != nil {
			fmt.Println("Error starting server", err)
		}
	}()

	<-shutdownChan
	fmt.Println("Gracefully shutting down server...")
	if err := app.Shutdown(); err != nil {
		fmt.Println("Error shutting down server", err)
	}
	fmt.Println("Server shutdown")
}
