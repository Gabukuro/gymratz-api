package setup_test

import (
	"os"
	"testing"

	setup "github.com/Gabukuro/gymratz-api/internal/pkg/setup"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	t.Parallel()

	t.Run("should init default configurations", func(t *testing.T) {
		t.Parallel()

		os.Setenv("APPLICATION_NAME", "TestApp")
		os.Setenv("DATABASE_URL", "postgres://localhost:5432")
		os.Setenv("JWT_SECRET", "secret")

		app := setup.Init()

		assert.NotNil(t, app.EnvVariables)
		assert.NotNil(t, app.BRLocation)
		assert.NotNil(t, app.ShutdownChan)
		assert.Equal(t, "TestApp", app.EnvVariables.ApplicationName)
		assert.Equal(t, "postgres://localhost:5432", app.EnvVariables.DatabaseURL)
		assert.Equal(t, "secret", app.EnvVariables.JWTSecret)
	})

	t.Run("should shutdown using the interrupt chan", func(t *testing.T) {
		t.Parallel()

		setup := setup.Init()

		hasShutdown := false

		go func() {
			<-setup.ShutdownChan

			hasShutdown = true

			setup.Shutdown()
		}()

		setup.ShutdownChan <- os.Interrupt
		setup.WaitShutdown()

		assert.NotNil(t, setup.ShutdownChan)
		assert.True(t, hasShutdown)
	})
}
