package adapter

import (
	"context"

	"github.com/Goboolean/fetch-system.worker/internal/domain/port/out"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/mock"
)




type MockGeneratorAdapter struct {
	g *mock.Client
}

func NewMockGeneratorAdapter(g *mock.Client) (out.DataFetcher, error) {
	return &MockGeneratorAdapter{
		g: g,
	}, nil
}

func (a *MockGeneratorAdapter) InputStream(ctx context.Context, symbols ...string) (<-chan *vo.Trade, error) {

	ch := make(chan *vo.Trade)

	go func() {
		for v := range a.g.Subscribe(symbols...) {
			ch <- &vo.Trade{
				Symbol: v.Symbol,
				TradeDetail: vo.TradeDetail{
					Price: v.Price,
					Size: v.Amount,
					Timestamp: v.Timestamp,
				},
			}
		}
	}()

	return ch, nil
}