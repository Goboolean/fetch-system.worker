package polygon

import (
	"context"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/polygon-io/client-go/websocket"
	"github.com/polygon-io/client-go/websocket/models"
)




type StocksClient struct {
	*client[models.EquityTrade]
}


func (c *StocksClient) Ping(ctx context.Context) error {
	return c.client.Ping(ctx)
}

func (c *StocksClient) Close() {
	c.client.Close()
}

func NewStocksClient(c *resolver.ConfigMap) (*StocksClient, error) {
	
	if err := c.SetKey("FEED", string(polygonws.Delayed)); err != nil {
		return nil, err
	}

	if err := c.SetKey("MARKET", string(polygonws.Stocks)); err != nil {
		return nil, err
	}

	if err := c.SetKey("TOPIC", int(polygonws.StocksTrades)); err != nil {
		return nil, err
	}

	client, err := newClient[models.EquityTrade](c)
	if err != nil {
		return nil, err
	}

	return &StocksClient{
		client: client,
	}, nil
}

func (c *StocksClient) Subscribe() (<-chan models.EquityTrade, error) {
	return c.client.Subscribe()
}


func (c *StocksClient) IsMarketOn(ctx context.Context) (bool, error) {
	
	resp, err := c.rest.GetMarketStatus(ctx)
	if err != nil {
		return false, err
	}

	return (resp.Market == "open"), nil
}