package task_test

import (
	"context"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/adapter"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)



func TestShutdownError(t *testing.T) {
	w1 := vo.Worker{
		ID: "first",
		Platform: vo.PlatformKIS,
	}

	w2 := vo.Worker{
		ID: "second",
		Platform: vo.PlatformKIS,
	}

	w3 := vo.Worker{
		ID: "third",
		Platform: vo.PlatformKIS,
	}

	m1 := SetupTaskManager(&w1)
	m2 := SetupTaskManager(&w2)
	m3 := SetupTaskManager(&w3)

	t.Cleanup(func() {
		etcdStub.(*adapter.ETCDStub).Cleanup()
		workers, err :=etcdStub.GetAllWorker(context.Background())
		assert.NoError(t, err)
		assert.Empty(t, workers)
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

		err = m3.RegisterWorker(context.Background())
		assert.NoError(t, err)

		w, err = etcdStub.GetWorker(context.Background(), w3.ID)
		assert.NoError(t, err)
		assert.Equal(t, w.Status, vo.WorkerStatusSecondary)
	})

	t.Run("ShutdownOccuredOnPrimary", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
		defer cancel()

		err := m1.Shutdown()
		assert.NoError(t, err)

		etcdStub.(*adapter.ETCDStub).CreateShutdownEvent(ctx)
		time.Sleep(1 * time.Millisecond)

		_time, err := etcdStub.GetWorkerTimestamp(ctx, w1.ID)
		assert.NoError(t, err)
		assert.Greater(t, _time, time.Now())

		w1, err := etcdStub.GetWorker(ctx, w1.ID)
		assert.NoError(t, err)
		assert.Equal(t, vo.WorkerStatusExitedShutdownOccured, w1.Status)

		w2, err := etcdStub.GetWorker(ctx, w2.ID)
		assert.NoError(t, err)
		w3, err := etcdStub.GetWorker(ctx, w3.ID)
		assert.NoError(t, err)

		assert.True(t, w2.Status == vo.WorkerStatusPrimary   || w3.Status == vo.WorkerStatusPrimary)
		assert.True(t, w2.Status == vo.WorkerStatusSecondary || w3.Status == vo.WorkerStatusSecondary)
	})
}



func TestTTLFailed(t *testing.T) {
	w1 := vo.Worker{
		ID: "first",
		Platform: vo.PlatformKIS,
	}

	w2 := vo.Worker{
		ID: "second",
		Platform: vo.PlatformKIS,
	}

	w3 := vo.Worker{
		ID: "third",
		Platform: vo.PlatformKIS,
	}

	m1 := SetupTaskManager(&w1)
	m2 := SetupTaskManager(&w2)
	m3 := SetupTaskManager(&w3)

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

		err = m3.RegisterWorker(context.Background())
		assert.NoError(t, err)

		w, err = etcdStub.GetWorker(context.Background(), w3.ID)
		assert.NoError(t, err)
		assert.Equal(t, w.Status, vo.WorkerStatusSecondary)
	})

	t.Run("TTLFailed", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
		defer cancel()

		etcdStub.(*adapter.ETCDStub).CreateTTLFailedEvent(ctx)
		time.Sleep(1 * time.Millisecond)
		etcdStub.(*adapter.ETCDStub).CreateTTLFailedEvent(ctx)
		time.Sleep(1 * time.Millisecond)

		_time, err := etcdStub.GetWorkerTimestamp(ctx, w1.ID)
		assert.NoError(t, err)
		assert.LessOrEqual(t, _time, time.Now())

		w1, err := etcdStub.GetWorker(ctx, w1.ID)
		assert.NoError(t, err)
		assert.Equal(t, vo.WorkerStatusExitedTTlFailed, w1.Status)

		w2, err := etcdStub.GetWorker(ctx, w2.ID)
		assert.NoError(t, err)
		w3, err := etcdStub.GetWorker(ctx, w3.ID)
		assert.NoError(t, err)

		t.Log(w2.Status, w3.Status)

		assert.True(t, w2.Status == vo.WorkerStatusPrimary   || w3.Status == vo.WorkerStatusPrimary)
		assert.True(t, w2.Status == vo.WorkerStatusSecondary || w3.Status == vo.WorkerStatusSecondary)
	})
}