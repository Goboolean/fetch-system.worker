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

	duration := time.Millisecond * 100

	t.Run("RunPipe", func(t *testing.T) {
		err := p.RunStreamingPipe(context.Background(), []*vo.Product{})
		assert.NoError(t, err)

		time.Sleep(duration)

		assert.LessOrEqual(t, 85, receiver.Received())
		assert.GreaterOrEqual(t, 105, receiver.Received())
	})

	t.Run("LockupPipe", func(t *testing.T) {
		p.LockupPipe(time.Now().Add(duration))
		time.Sleep(duration * 2)

		got := receiver.Received()
		assert.LessOrEqual(t, 170, got)
		assert.GreaterOrEqual(t, 210, got)

		time.Sleep(duration)
		assert.Equal(t, got, receiver.Received())

		receiver.Initialize()
	})
}



func TestStoringPipe(t *testing.T) {

	fetcher := SetupMockGenerator(1000)
	receiver := NewMockReceiver()

	p := pipe.New(fetcher, receiver)

	duration := time.Millisecond * 100

	t.Run("RunPipe", func(t *testing.T) {
		err := p.RunStoringPipe(context.Background(), []*vo.Product{})
		assert.NoError(t, err)

		time.Sleep(duration)
		assert.Equal(t, 0, receiver.Received())
	})

	t.Run("UpgradePipe", func(t *testing.T) {
		p.UpgradeToStreamingPipe(time.Now().Add(duration))

		time.Sleep(time.Millisecond)
		got := receiver.Received()
		assert.LessOrEqual(t, 85, got)
		assert.GreaterOrEqual(t, 105, got)

		time.Sleep(duration)
		got = receiver.Received()
		assert.LessOrEqual(t, 170, got)
		assert.GreaterOrEqual(t, 210, got)
	})

	t.Run("LockupPipe", func(t *testing.T) {
		p.LockupPipe(time.Now().Add(duration))
		time.Sleep(duration * 2)

		got := receiver.Received()
		assert.LessOrEqual(t, 255, got)
		assert.GreaterOrEqual(t, 315, got)

		time.Sleep(duration)
		assert.Equal(t, got, receiver.Received())

		receiver.Initialize()
	})
}