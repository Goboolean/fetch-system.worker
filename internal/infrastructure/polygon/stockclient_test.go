package polygon_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/polygon"
	"github.com/stretchr/testify/assert"

	_ "github.com/Goboolean/common/pkg/env"
)



func SetupStockClient() *polygon.StocksClient {
	c, err := polygon.NewStocksClient(&resolver.ConfigMap{
		"SECRET_KEY": os.Getenv("POLYGON_SECRET_KEY"),
		"BUFFER_SIZE": 10000,
	})
	if err != nil {
		panic(err)
	}

	return c
}

func TeardownStockClient(c *polygon.StocksClient) {
	c.Close()
}


func TestStockClient(t *testing.T) {

	c := SetupStockClient()
	defer TeardownStockClient(c)

	t.Run("Ping()", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
		defer cancel()

		
		err := c.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("Subscribe()", func(t *testing.T) {
		ch, err := c.Subscribe()
		assert.NoError(t, err)
		
		time.Sleep(time.Second * 1)

		t.Log(len(ch))
	})
}