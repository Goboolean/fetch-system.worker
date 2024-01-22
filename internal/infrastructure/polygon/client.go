package polygon

import (
	"context"
	"fmt"
	"sync"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.worker/internal/util/otel"
	polygonrest "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/websocket"
	"github.com/polygon-io/client-go/websocket/models"
	log "github.com/sirupsen/logrus"
)



type client[T models.EquityTrade | models.CryptoTrade] struct {
	conn  *polygonws.Client
	rest  *polygonrest.Client

	ch chan T

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	topic polygonws.Topic
}


func (c *client[T]) Ping(ctx context.Context) error {

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

func (c *client[T]) Close() {
	c.cancel()
	c.conn.Close()
	c.wg.Wait()
	close(c.ch)
}

func newClient[T models.EquityTrade | models.CryptoTrade](c *resolver.ConfigMap) (*client[T], error) {

	key, err := c.GetStringKey("SECRET_KEY")
	if err != nil {
		return nil, err
	}

	feed, err := c.GetStringKey("FEED")
	if err != nil {
		return nil, err
	}

	market, err := c.GetStringKey("MARKET")
	if err != nil {
		return nil, err
	}

	conn, err := polygonws.New(polygonws.Config{
		APIKey: key,
		Feed:   polygonws.Feed(feed),
		Market: polygonws.Market(market),
	})
	if err != nil {
		return nil, err
	}

	topic, err := c.GetIntKey("TOPIC")
	if err != nil {
		return nil, err
	}

	buf, err := c.GetIntKey("BUFFER_SIZE")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	rest := polygonrest.New(key)

	return &client[T]{
		conn: conn,
		ch: make(chan T, buf),
		rest: rest,

		ctx: ctx,
		cancel: cancel,

		topic: polygonws.Topic(topic),
	}, nil
}

func (c *client[T]) Subscribe() (<-chan T, error) {

	if err := c.conn.Subscribe(c.topic); err != nil {
		return nil, err
	}

	c.wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case err := <-c.conn.Error():
				log.WithField("error", err).
				Panic("Failed to receive data from polygon")
				return
			case out, more := <-c.conn.Output():
				if !more {
					log.WithField("error", fmt.Errorf("there is no more output from websocket")).
					Panic("Failed to receive data from polygon")
					return
				}

				data, ok := out.(T)
				if !ok {
					log.WithField("data", out).
					Error("Failed to cast data. There is inconsistency between the infrastructure type and data type")
					otel.PolygonStockErrorCount.Add(ctx, 1)
					continue
				}

				otel.PolygonStockReceivedCount.Add(ctx, 1)
				c.ch <- data
			}
		}
	}(c.ctx, &c.wg)

	return c.ch, nil
}

