package polygon

import (
	"context"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/polygon-io/client-go/websocket"
	"github.com/polygon-io/client-go/websocket/models"
)


type CryptoClient struct {
	*client[models.CryptoTrade]
}


func (c *CryptoClient) Ping(ctx context.Context) error {
	return c.client.Ping(ctx)
}

func (c *CryptoClient) Close() {
	c.client.Close()
}

func NewCryptoClient(c *resolver.ConfigMap) (*CryptoClient, error) {

	if err := c.SetKey("FEED", string(polygonws.Delayed)); err != nil {
		return nil, err
	}

	if err := c.SetKey("MARKET", string(polygonws.Crypto)); err != nil {
		return nil, err
	}

	if err := c.SetKey("TOPIC", int(polygonws.CryptoTrades)); err != nil {
		return nil, err
	}

	client, err := newClient[models.CryptoTrade](c)
	if err != nil {
		return nil, err
	}

	return &CryptoClient{
		client: client,
	}, nil
}

func (c *CryptoClient) Subscribe() (<-chan models.CryptoTrade, error) {
	return c.client.Subscribe()
}