package etcd_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/etcd"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)



func TestKeepAlive(t *testing.T) {
	var keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	var clientStatusChan <-chan clientv3.WatchResponse
	var clientKeepAliveChan <-chan clientv3.WatchResponse

	var leaseID clientv3.LeaseID

	var worker = etcd.Worker{
		ID: uuid.New().String(),
		Platform: "polygon",
		Status: "primary",
	}

	t.Run("SetupWorker", func(t *testing.T) {
		err := client.InsertWorker(context.Background(), &worker)
		assert.NoError(t, err)
	})

	t.Cleanup(func() {
		err := client.DeleteWorker(context.Background(), worker.ID)
		assert.NoError(t, err)
	})

	t.Run("CreateKeepAlive", func(t *testing.T) {
		var err error

		keepAliveChan, leaseID, err = client.CreateKeepAlive(context.Background())
		assert.NoError(t, err)

		err = client.SetKeepAliveValue(context.Background(), worker.ID, leaseID)
		assert.NoError(t, err)
	})

	t.Run("UpdateStatus", func(t *testing.T) {
		clientStatusChan = client.WatchWorkerStatus(context.Background(), worker.ID)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := client.UpdateWorkerStatus(context.Background(), worker.ID, "dead")
		assert.NoError(t, err)

		select {
		case <-ctx.Done():
			assert.Fail(t, "timeout")
		case resp := <-clientStatusChan:
			assert.Equal(t, "dead", string(resp.Events[0].Kv.Value))
			break
		}
	})

	t.Run("ShutdownKeepAlive", func(t *testing.T) {
		clientKeepAliveChan = client.WatchWorkerKeepAlive(context.Background(), worker.ID)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := client.ShutdownKeepAlive(context.Background(), leaseID)
		assert.NoError(t, err)

		select {
		case <-ctx.Done():
			assert.Fail(t, "timeout")
		case <-keepAliveChan:
			break
		}

		select {
		case <-ctx.Done():
			assert.Fail(t, "timeout")
		case resp := <-clientKeepAliveChan:
			assert.Equal(t, "", string(resp.Events[0].Kv.Value))
			break
		}
	})
}



func Test_Mutex(t *testing.T) {

	var mu1 *concurrency.Mutex
	var mu2 *concurrency.Mutex

	var key string = "test"

	t.Run("AquireMutex", func(t *testing.T) {
		var err error
		mu1, err = client.NewMutex(context.Background(), key)
		assert.NoError(t, err)
		mu2, err = client.NewMutex(context.Background(), key)
		assert.NoError(t, err)
	})

	t.Run("Lock", func(t *testing.T) {
		err := mu1.Lock(context.Background())
		assert.NoError(t, err)

		err = mu2.TryLock(context.Background())
		assert.Error(t, err)
	})

	t.Run("WaitLock", func(t *testing.T) {

		var wg sync.WaitGroup
		go func() {
			wg.Add(1)
			defer wg.Done()

			time.Sleep(1 * time.Second)
			err := mu1.Unlock(context.Background())
			assert.NoError(t, err)
		}()

		err := mu2.Lock(context.Background())
		assert.NoError(t, err)
		wg.Wait()
	})

	t.Run("Unlock", func(t *testing.T) {
		err := mu1.Unlock(context.Background())
		assert.NoError(t, err)

		err = mu2.TryLock(context.Background())
		assert.NoError(t, err)

		err = mu2.Unlock(context.Background())
		assert.NoError(t, err)
	})
}
