package setup

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type (
	Setup struct {
		ApplicationName string
		BRLocation      time.Location
		ShutdownChan    chan os.Signal

		shutdownWaitGroup sync.WaitGroup
	}

	SetupConfig struct {
		ApplicationName string
	}
)

func Init(config SetupConfig) *Setup {
	app := &Setup{
		ApplicationName: config.ApplicationName,
	}

	app.configureBRLocation()
	app.configureGracefulShutdown()

	return app
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

func (s *Setup) Shutdown() {
	s.shutdownWaitGroup.Done()
	fmt.Println("Server stopped!")
}
