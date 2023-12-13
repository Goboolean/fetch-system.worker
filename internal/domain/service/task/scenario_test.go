package task_test

import (
	"context"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/adapter"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)



func TestTTLFailed(t *testing.T) {
	w1 := vo.Worker{
		ID: "first",
		Platform: vo.PlatformKIS,
	}

	w2 := vo.Worker{
		ID: "second",
		Platform: vo.PlatformKIS,
	}

	m1 := SetupTaskManager(&w1)
	m2 := SetupTaskManager(&w2)

	t.Cleanup(func() {
		etcdStub.(*adapter.ETCDStub).DeleteAllWorkers(context.Background())
	})

	t.Run("RegisterPrimary", func(t *testing.T) {
		err := m1.RegisterWorker(context.Background())
		assert.NoError(t, err)

		w, err := etcdStub.GetWorker(context.Background(), w1.ID)
		assert.NoError(t, err)
		assert.Equal(t, w.Status, vo.WorkerStatusPrimary)
	})

	t.Run("RegisterSecondary", func(t *testing.T) {
		err := m2.RegisterWorker(context.Background())
		assert.NoError(t, err)

		w, err := etcdStub.GetWorker(context.Background(), w2.ID)
		assert.NoError(t, err)
		assert.Equal(t, w.Status, vo.WorkerStatusSecondary)
	})

	t.Run("ShutdownOccuredOnPrimary", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
		defer cancel()

		err := m1.Shutdown()
		assert.NoError(t, err)

		time.Sleep(1 * time.Millisecond)

		_time, err := etcdStub.(*adapter.ETCDStub).GetWorkerTimestamp(context.Background(), w1.ID)
		assert.NoError(t, err)
		assert.Greater(t, _time.Unix(), time.Now().Unix())

		w1, err := etcdStub.GetWorker(ctx, w1.ID)
		assert.NoError(t, err)
		assert.Equal(t, vo.WorkerStatusExitedShutdownOccured, w1.Status)

		w2, err := etcdStub.GetWorker(ctx, w2.ID)
		assert.NoError(t, err)
		assert.Equal(t, vo.WorkerStatusPrimary, w2.Status)
	})
}
