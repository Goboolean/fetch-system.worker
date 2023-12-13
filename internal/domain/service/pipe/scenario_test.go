package pipe_test

import (
	"context"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/domain/service/pipe"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)



func TestStreamingPipe(t *testing.T) {

	fetcher := SetupMockGenerator(1000)
	receiver := NewMockReceiver()

	p := pipe.New(fetcher, receiver)

	t.Run("RunPipe", func(t *testing.T) {
		err := p.RunStreamingPipe(context.Background(), []*vo.Product{})
		assert.NoError(t, err)

		time.Sleep(time.Millisecond * 100)

		assert.LessOrEqual(t, 90, receiver.Received())
		assert.GreaterOrEqual(t, 110, receiver.Received())

		receiver.Initialize()
	})

	t.Run("LockupPipe", func(t *testing.T) {
		duration := time.Millisecond * 100

		p.LockupPipe(time.Now().Add(duration))
		time.Sleep(duration)

		got := receiver.Received()
		assert.LessOrEqual(t, 90, got)
		assert.GreaterOrEqual(t, 110, got)

		time.Sleep(duration)
		assert.Equal(t, got, receiver.Received())

		receiver.Initialize()
	})
}