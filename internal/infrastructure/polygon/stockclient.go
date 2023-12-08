package polygon

import (
	"context"
	"sync"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/polygon-io/client-go/websocket"
	"github.com/polygon-io/client-go/websocket/models"
)


type StocksClient struct {
	conn  *polygonws.Client

	ch chan models.EquityTrade

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}


func (c *StocksClient) Ping(ctx context.Context) error {

	ch := make(chan error)

	go func(ch chan error) {
		ch <- c.conn.Connect()
	}(ch)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ch:
		return err
	}
}

func (c *StocksClient) Close() {
	c.cancel()
	c.conn.Close()
	c.wg.Wait()
	close(c.ch)
}

func NewStocksClient(c *resolver.ConfigMap) (*StocksClient, error) {

	key, err := c.GetStringKey("SECRET_KEY")
	if err != nil {
		return nil, err
	}

	conn, err := polygonws.New(polygonws.Config{
		APIKey: key,
		Feed:   polygonws.Delayed,
		Market: polygonws.Stocks,
	})

	if err != nil {
		return nil, err
	}


	buf, err := c.GetIntKey("BUFFER_SIZE")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &StocksClient{
		conn: conn,
		ch: make(chan models.EquityTrade, buf),

		ctx: ctx,
		cancel: cancel,
	}, nil
}

func (c *StocksClient) Subscribe() (<-chan models.EquityTrade, error) {

	if err := c.conn.Subscribe(polygonws.StocksTrades); err != nil {
		return nil, err
	}

	c.wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-c.conn.Error():
				return
			case out, more := <-c.conn.Output():
				if !more {
					return
				}

				data, ok := out.(models.EquityTrade)
				if !ok {
					return
				}

				c.ch <- data
			}
		}
	}(c.ctx, &c.wg)

	return c.ch, nil
}