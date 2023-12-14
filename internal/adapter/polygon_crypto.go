package adapter

import (
	"context"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/domain/port/out"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/polygon"
)


type CryptoPolygonAdapter struct {
	c *polygon.CryptoClient
}

func NewCryptoPolygonAdapter(c *polygon.CryptoClient) (out.DataFetcher, error) {
	return &CryptoPolygonAdapter{
		c: c,
	}, nil
}

func (a *CryptoPolygonAdapter) InputStream(ctx context.Context, symbol ...string) (<-chan *vo.Trade, error) {
	
	ch := make(chan *vo.Trade)

	polyonCh, err := a.c.Subscribe()
	if err != nil {
		return nil, err
	}

	go func() {
		for v := range polyonCh {
			ch <- &vo.Trade{
				ID: v.Pair,
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