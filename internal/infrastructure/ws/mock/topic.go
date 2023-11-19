package mock

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/ws"
)

type mockGenerater struct {
	symbol string

	ctx    context.Context
	cancel context.CancelFunc
	ch     chan<- *ws.StockAggregate
	d      time.Duration

	curTime  time.Time
	curPrice float64
}

func (m *mockGenerater) generateRandomStockAggs() *ws.StockAggregate {

	lastTime := m.curTime
	lastPrice := m.curPrice

	curTime := time.Now()
	curPrice := m.curPrice * (rand.Float64()*0.2 + 0.9)

	stockAggs := &ws.StockAggregate{
		StockAggsMeta: ws.StockAggsMeta{
			Platform: platformName,
			Symbol:   m.symbol,
		},
		StockAggsDetail: ws.StockAggsDetail{
			StartTime: lastTime.UnixNano(),
			EndTime:   curTime.UnixNano(),
			Average:   (lastPrice + curPrice) / 2,
			Min:       math.Min(lastPrice, curPrice),
			Max:       math.Max(lastPrice, curPrice),
			Start:     lastPrice,
			End:       curPrice,
		},
	}

	m.curTime = curTime
	m.curPrice = curPrice

	return stockAggs
}

func newMockGenerater(symbol string, ctx context.Context, ch chan<- *ws.StockAggregate, d time.Duration) *mockGenerater {
	ctx, cancel := context.WithCancel(ctx)

	instance := &mockGenerater{
		ctx:    ctx,
		cancel: cancel,
		ch:     ch,
		d:      d,
		symbol: symbol,
	}

	instance.curTime = time.Now()
	instance.curPrice = 1000

	go func(ctx context.Context) {

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(instance.d):
				agg := instance.generateRandomStockAggs()
				ch <- agg
				break
			}
		}
	}(ctx)

	return instance
}

// GC will not immediately release the memory mockGenerater occupies,
// therefore be sure to call Close() when you are done with mockGenerater.
func (m *mockGenerater) Close() {
	m.cancel()
}
