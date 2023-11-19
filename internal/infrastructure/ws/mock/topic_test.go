package mock

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/ws"
)

var (
	generater *mockGenerater
	ch        chan *ws.StockAggregate
	symbol    = "MOCK"
)

func SetupMockGenerater() {
	ctx := context.Background()
	ch = make(chan *ws.StockAggregate)

	duration := time.Second / 1000 // the data is generated every 0.1 second in average.

	generater = newMockGenerater(symbol, ctx, ch, duration)
}

func TeardownMockGenerater() {
	close(ch)
	generater.Close()
}

// It fails the test when it generates empty data.
func Test_generateRandomStockAggs(t *testing.T) {

	SetupMockGenerater()

	agg := generater.generateRandomStockAggs()
	if same := reflect.DeepEqual(agg, ws.StockAggregate{}); same {
		t.Errorf("generateRandomStockAggs() = %v, want not empty", agg)
		return
	}

	TeardownMockGenerater()
}

// It verdicts the test as success when it generates data 5 times for a second.
func Test_newMockGenerater(t *testing.T) {

	SetupMockGenerater()
	defer TeardownMockGenerater()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()

	t.Run("generateRandomStockAggs", func(t *testing.T) {
		for count := 5; count >= 0; count-- {
			select {
			case <-ctx.Done():
				t.Errorf("newMockGenerater() got timeout")
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
			t.Errorf("newMockGenerater got data after closing")
			return
		case <-time.After(time.Second / 10):
			break
		}
	})
}
