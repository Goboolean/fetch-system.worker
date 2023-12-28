package kis_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	_ "github.com/Goboolean/common/pkg/env"
)





func TestConstructor(t *testing.T) {

	t.Run("Ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		err := client.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("IsMarketOn", func(t *testing.T) {
		if !client.AccessTokenAvailable() {
			t.Skip("Cannot test IsMarketOn without access token")
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*61)
		defer cancel()

		_, err := client.IsMarketOn(ctx)
		assert.NoError(t, err)
	})
}



func TestWebsocket(t *testing.T) {

	const symbol = "005930"

	t.Run("Subscribe", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		ch, err := client.Subscribe(ctx, symbol)
		assert.NoError(t, err)

		//if !client.AccessTokenAvailable() {
		//	t.Skip("Cannot test IsMarketOn without access token")
		//}

		//on, err := client.IsMarketOn(ctx)
		//assert.NoError(t, err)

		for {
			select {
			case <-ctx.Done():
				return
			case v := <-ch:
				assert.Equal(t, symbol, v.Symbol)
			}
		}
	})
}