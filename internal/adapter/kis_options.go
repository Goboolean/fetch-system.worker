package adapter

import (
	"context"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/domain/port/out"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/kis"
)


type OptionsKISAdapter struct {
	c *kis.Client
}

func NewOptionsKISAdapter(c *kis.Client) (out.DataFetcher, error) {
	return &OptionsKISAdapter{
		c: c,
	}, nil
}

func (a *OptionsKISAdapter) InputStream(ctx context.Context, symbol ...string) (<-chan *vo.Trade, error) {
	
	ch := make(chan *vo.Trade)

	polyonCh, err := a.c.Subscribe(ctx, symbol...)
	if err != nil {
		return nil, err
	}

	go func() {
		for v := range polyonCh {
			ch <- &vo.Trade{
				ID: v.Symbol,
				TradeDetail: vo.TradeDetail{
					Price: v.Price,
					Size: int64(v.Size),
					Timestamp: time.Unix(v.Timestamp / 1000, v.Timestamp % 1000 * 1000000),
				},
			}
		}
	}()

	return ch, nil
}