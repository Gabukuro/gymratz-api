package setup_test

import (
	"os"
	"testing"

	setup "github.com/Gabukuro/gymratz-api/internal/pkg/setup"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	t.Parallel()

	t.Run("should shutdown using the interrupt chan", func(t *testing.T) {
		t.Parallel()

		setup := setup.Init(setup.SetupConfig{
			ApplicationName: "TestApp",
		})

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
