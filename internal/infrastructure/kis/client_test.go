package kis_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/kis"
	"github.com/stretchr/testify/assert"

	_ "github.com/Goboolean/common/pkg/env"
)


var client *kis.Client

func SetupKis() *kis.Client {
	var err error

	c, err := kis.New(&resolver.ConfigMap{
		"APPKEY": os.Getenv("KIS_APPKEY"),
		"SECRET": os.Getenv("KIS_SECRET"),
		"BUFFER_SIZE": 10000,
	})
	if err != nil {
		panic(err)
	}
	return c
}

func TeardownKis(c *kis.Client) {
	c.Close()
}


func TestConstructor(t *testing.T) {

	c := SetupKis()
	defer TeardownKis(c)

	t.Run("Ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		err := c.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("IsMarketOn", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*61)
		defer cancel()

		_, err := c.IsMarketOn(ctx)
		assert.NoError(t, err)
	})
}



func Test_KIS(t *testing.T) {

	c := SetupKis()
	defer TeardownKis(c)

	const symbol = "005930"

	t.Run("Subscribe", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		_, err := c.Subscribe(ctx, symbol)
		assert.NoError(t, err)
	})
}