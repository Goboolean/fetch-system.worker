package polygon

import (
	"context"
	"sync"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/polygon-io/client-go/websocket"
	"github.com/polygon-io/client-go/websocket/models"
)


type CryptoClient struct {
	conn  *polygonws.Client

	ch chan models.CryptoTrade

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}


func (c *CryptoClient) Ping(ctx context.Context) error {

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

func (c *CryptoClient) Close() {
	c.cancel()
	c.conn.Close()
	c.wg.Wait()
	close(c.ch)
}

func NewCryptoClient(c *resolver.ConfigMap) (*CryptoClient, error) {

	key, err := c.GetStringKey("SECRET_KEY")
	if err != nil {
		return nil, err
	}

	conn, err := polygonws.New(polygonws.Config{
		APIKey: key,
		Feed:   polygonws.Delayed,
		Market: polygonws.Crypto,
	})

	if err != nil {
		return nil, err
	}

	buf, err := c.GetIntKey("BUFFER_SIZE")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &CryptoClient{
		conn: conn,
		ch: make(chan models.CryptoTrade, buf),

		ctx: ctx,
		cancel: cancel,
	}, nil
}

func (c *CryptoClient) Subscribe() (chan models.EquityTrade, error) {
	ch := make(chan models.EquityTrade)

	if err := c.conn.Subscribe(polygonws.CryptoTrades); err != nil {
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

				data, ok := out.(models.CryptoTrade)
				if !ok {
					return
				}

				c.ch <- data
			}
		}
	}(c.ctx, &c.wg)

	return ch, nil
}