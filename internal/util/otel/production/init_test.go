package production_test

import (
	"context"
	"testing"

	"github.com/Goboolean/fetch-system.worker/internal/util/otel/production"
	"github.com/stretchr/testify/assert"

	_ "github.com/Goboolean/common/pkg/env"
)


func TestInitOtelProductionMode(t *testing.T) {
	t.Run("Initialize & Close Otel Debug Mode", func(t *testing.T) {
		err := production.Close(context.Background())
		assert.NoError(t, err)
	})
}