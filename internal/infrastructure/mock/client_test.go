package mock_test

import (
	"context"
	"testing"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/mock"
	"github.com/stretchr/testify/assert"
)




func Test_Generate1mPerSec(t *testing.T) {

	var g *mock.Client

	const (
		productionRate = 100
		productNum = 1000
		total = productionRate * productNum
	)

	var received int = 0
	var ch = make(<-chan *mock.Trade, total)

	t.Run("Initialize Generator", func(t *testing.T) {
		var err error
		g, err = mock.New(&resolver.ConfigMap{
			"MODE": "BASIC",
			"STANDARD_DEVIATION": 100,
			"PRODUCT_COUNT": productNum,
			"PRODUCTION_RATE": productionRate,
			"BUFFER_SIZE": total,
		})
		assert.NoError(t, err)
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 1)
	defer cancel()

	t.Run("Run Generator", func(t *testing.T) {
		ch = g.Subscribe()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ch:
				received++
			}
		}
	})

	t.Run("Stop Generator", func(t *testing.T) {
		g.Close()

		time.Sleep(time.Millisecond * 10)
		received += len(ch)
	})

	assert.LessOrEqual(t, int(total * 0.90), received)
	assert.GreaterOrEqual(t, int(total * 1.1), received)
}

func TeardownMockGenerater(m *mock.Client) {
	m.Close()
}