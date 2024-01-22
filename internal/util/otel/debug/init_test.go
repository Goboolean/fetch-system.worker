package debug_test

import (
	"context"
	"testing"

	"github.com/Goboolean/fetch-system.worker/internal/util/otel/debug"
	"github.com/stretchr/testify/assert"
)


func TestInitOtelDebugMode(t *testing.T) {
	t.Run("Initialize & Close Otel Debug Mode", func(t *testing.T) {
		err := debug.Close(context.Background())
		assert.NoError(t, err)
	})
}