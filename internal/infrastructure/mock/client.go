package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
)



type Client struct {
	ctx context.Context
	cancel context.CancelFunc

	config struct {
		mode string
		sd int
		cnt int
		rate int
	}

	ch chan *Trade
}


func New(c *resolver.ConfigMap) (*Client, error) {

	var (
		mode string
		sd int
		cnt int
		rate int
	)

	mode, err := c.GetStringKey("MODE")
	if err != nil {
		return nil, err
	}

	switch mode {

	case "BASIC":
		sd, err = c.GetIntKey("STANDARD_DEVIATION")
		if err != nil {
			return nil, err
		}

		cnt, err = c.GetIntKey("PRODUCT_COUNT")
		if err != nil {
			return nil, err
		}

		rate, err = c.GetIntKey("PRODUCTION_RATE")
		if err != nil {
			return nil, err
		}

	case "CUSTOM":
		return &Client{}, fmt.Errorf("not implemented yet")
	}

	size, err := c.GetIntKey("BUFFER_SIZE")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		ctx: ctx,
		cancel: cancel,

		config: struct {
			mode string
			sd int
			cnt int
			rate int
		}{mode, sd, cnt, rate},

		ch: make(chan *Trade, size),
	}, nil
}

func (c *Client) Subscribe(symbols ...string) <-chan *Trade {

	if len(symbols) == 0 {
		for i := 0; i < c.config.cnt; i++ {
			dur := time.Second / time.Duration(c.config.rate)
			go RunGenerator(c.ctx, Symbol(), dur, c.ch)
		}		
	}

	if len(symbols) > 0 {
		for _, symbol := range symbols {
			dur := time.Second / time.Duration(c.config.rate)
			go RunGenerator(c.ctx, symbol, dur, c.ch)
		}
	}

	return c.ch
}

func (c *Client) Close() {
	c.cancel()
}