package mock

import (
	"context"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/ws"
	"github.com/stretchr/testify/assert"
)

var (
	generater *mockGenerater
	ch        chan *ws.StockAggregate
	symbol    = "MOCK"
)

func SetupMockGenerater() {
	ctx := context.Background()
	ch = make(chan *ws.StockAggregate)

	duration := time.Second / 1000 // the data is generated every 0.1 milisecond in average.

	generater = newMockGenerater(symbol, ctx, ch, duration)
}

func TeardownMockGenerater() {
	close(ch)
	generater.Close()
}

// It fails the test when it generates empty data.
func Test_generateRandomStockAggs(t *testing.T) {

	SetupMockGenerater()
	defer TeardownMockGenerater()

	agg := generater.generateRandomStockAggs()
	assert.NotEqual(t, agg, ws.StockAggregate{})
}

// It verdicts the test as success when it generates data 5 times for a second.
func Test_newMockGenerater(t *testing.T) {

	const count = 5000

	SetupMockGenerater()
	defer TeardownMockGenerater()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()

	t.Run("generateRandomStockAggs", func(t *testing.T) {
		for i := 5; i >= 0; i-- {
			select {
			case <-ctx.Done():
				assert.Fail(t, "newMockGenerater() got timeout")
				return
			case <-ch:
				continue
			}
		}
	})

	t.Run("afterStoptoGenerate", func(t *testing.T) {
		generater.Close()

		for len(ch) > 0 {
			<-ch
		}

		select {
		case <-ch:
			assert.Fail(t, "newMockGenerater got data after closing")
			return
		case <-time.After(time.Second / 10):
			break
		}
	})
}
