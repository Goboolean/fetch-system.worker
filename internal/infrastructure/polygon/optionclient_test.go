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



func SetupOptionClient() *polygon.OptionClient {
	c, err := polygon.NewOptionClient(&resolver.ConfigMap{
		"SECRET_KEY": os.Getenv("POLYGON_SECRET_KEY"),
		"BUFFER_SIZE": 100000,
	})
	if err != nil {
		panic(err)
	}

	return c
}

func TeardownOptionClient(c *polygon.OptionClient) {
	c.Close()
}


func TestOptionClient(t *testing.T) {

	t.Skip("Skipping the test, since it does not have privileges to access option api")

	c := SetupOptionClient()
	defer TeardownOptionClient(c)

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

		on, err := c.IsMarketOn(context.Background())
		assert.NoError(t, err)
		if on {
			assert.NotZero(t, len(ch))
			t.Log("ch size:", len(ch))
		}
	})
}