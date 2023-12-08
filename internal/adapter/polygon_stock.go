package adapter

import (
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/domain/port/out"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/polygon"
)


type StockPolygonAdapter struct {
	c *polygon.StocksClient
}

func NewStockPolygonAdapter(c *polygon.StocksClient) (out.DataFetcher, error) {
	return &StockPolygonAdapter{
		c: c,
	}, nil
}

func (a *StockPolygonAdapter) Subscribe(symbol ...string) (<-chan vo.Trade, error) {
	
	ch := make(chan vo.Trade)

	polyonCh, err := a.c.Subscribe()
	if err != nil {
		return nil, err
	}

	go func() {
		for v := range polyonCh {
			ch <- vo.Trade{
				ID: v.Symbol,
				TradeDetail: vo.TradeDetail{
					Price: v.Price,
					Size: v.Size,
					Timestamp: time.Unix(v.Timestamp / 1000, v.Timestamp % 1000 * 1000000),
				},
			}
		}
	}()

	return ch, nil
}