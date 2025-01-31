package main

import (
	"fmt"

	"github.com/Gabukuro/gymratz-api/internal/infra/database"
	"github.com/Gabukuro/gymratz-api/internal/pkg/setup"
)

func main() {
	setup, ctx := setup.Init()

	go func() {
		<-setup.ShutdownChan

		database.CloseDB()

		fmt.Println("Shutting down server...")
		setup.App.Shutdown()

		setup.Shutdown(ctx)
	}()

	if err := setup.App.Listen(":3000"); err != nil {
		fmt.Println("Error starting server:", err)
	}

	setup.WaitShutdown()
}
