package polygon

import (
	"context"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/polygon-io/client-go/websocket"
	"github.com/polygon-io/client-go/websocket/models"
)


type OptionClient struct {
	*client[models.EquityTrade]
}

func (c *OptionClient) Ping(ctx context.Context) error {
	return c.client.Ping(ctx)
}


func (c *OptionClient) Close() {
	c.client.Close()
}

func NewOptionClient(c *resolver.ConfigMap) (*OptionClient, error) {
	
	if err := c.SetKey("FEED", string(polygonws.Delayed)); err != nil {
		return nil, err
	}

	if err := c.SetKey("MARKET", string(polygonws.Options)); err != nil {
		return nil, err
	}

	if err := c.SetKey("TOPIC", int(polygonws.OptionsTrades)); err != nil {
		return nil, err
	}

	client, err := newClient[models.EquityTrade](c)
	if err != nil {
		return nil, err
	}

	return &OptionClient{
		client: client,
	}, nil
}


func (c *OptionClient) Subscribe() (<-chan models.EquityTrade, error) {
	return c.client.Subscribe()
}