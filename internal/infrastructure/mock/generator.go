package mock

import (
	"context"
	"math/rand"
	"sync"
	"time"
)


const defaultPrice = 1000

type generator struct {
	symbol string

	curTime  time.Time
	curPrice float64
}

func newGenerator(symbol string, d time.Duration) *generator {
	return &generator{
		symbol: symbol,
		curTime:  time.Now(),
		curPrice: defaultPrice,
	}
}

func (m *generator) Next() *Trade {

	curTime := time.Now()
	curPrice := m.curPrice * (rand.Float64()*0.2 + 0.9)

	trade := &Trade{
		Symbol:   m.symbol,
		Price:    curPrice,
		Amount:   1,
		Timestamp: curTime,
	}

	m.curTime = curTime
	m.curPrice = curPrice

	return trade
}

func RunGenerator(ctx context.Context, wg *sync.WaitGroup, id string, d time.Duration, ch chan<- *Trade) {
	defer wg.Done()

	g := newGenerator(id, d)
	ch <- g.Next()

	ticker := time.NewTicker(d)

	for range ticker.C {
		if ctx.Err() != nil {
			return
		}
		ch <- g.Next()
	}
}