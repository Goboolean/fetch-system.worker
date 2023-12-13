package etcd_test

import (
	"context"
	"sync"
	"testing"
	"time"

	_ "github.com/Goboolean/common/pkg/env"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/etcd"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/client/v3/concurrency"
)



func Test_Worker(t *testing.T) {

	var workers []*etcd.Worker = []*etcd.Worker{
		{ID: uuid.New().String(), Platform: "polygon", Status: "primary"},
		{ID: uuid.New().String(), Platform: "polygon", Status: "secondary"},
		{ID: uuid.New().String(), Platform: "kis",     Status: "exited"},
		{ID: uuid.New().String(), Platform: "kis",     Status: "onpromotion"},
	}

	t.Run("InsertWorker", func(t *testing.T) {
		for _, w := range workers {
			err := client.InsertWorker(context.Background(), w)
			assert.NoError(t, err)
		}

		for _, w := range workers {
			exists, err := client.WorkerExists(context.Background(), w.ID)
			assert.NoError(t, err)
			assert.True(t, exists)
		}
	})

	t.Run("InsertWorkerAlreadyExists", func(t *testing.T) {
		err := client.InsertWorker(context.Background(), workers[0])
		assert.Error(t, err)
	})

	t.Run("GetWorker", func(t *testing.T) {
		for _, w := range workers {
			worker, err := client.GetWorker(context.Background(), w.ID)
			assert.NoError(t, err)
			assert.Equal(t, w, worker)
		}
	})

	t.Run("UpdateWorkerStatus", func(t *testing.T) {
		err := client.UpdateWorkerStatus(context.Background(), workers[0].ID, "dead")
		assert.NoError(t, err)

		w, err := client.GetWorker(context.Background(), workers[0].ID)
		assert.NoError(t, err)
		assert.Equal(t, "dead", w.Status)
	})

	t.Run("DeleteWorker", func(t *testing.T) {
		err := client.DeleteWorker(context.Background(), workers[0].ID)
		assert.NoError(t, err)

		_, err = client.GetWorker(context.Background(), workers[0].ID)
		assert.Error(t, err)

		workers, err := client.GetAllWorkers(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, len(workers)-1, len(workers))
	})

	t.Run("DeleteAllWorkers", func(t *testing.T) {
		err := client.DeleteAllWorkers(context.Background())
		assert.NoError(t, err)

		workers, err := client.GetAllWorkers(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 0, len(workers))
	})
}



func Test_Product(t *testing.T) {

	var products []*etcd.Product = []*etcd.Product{
		{ID: "test.goboolean.kor", Platform: "kis",      Symbol: "goboolean", Type: "stock" },
		{ID: "test.goboolean.eng", Platform: "polygon",  Symbol: "gofalse",   Type: "crypto"},
		{ID: "test.goboolean.jpn", Platform: "buycycle", Symbol: "gonil",     Type: "option"},
		{ID: "test.goboolean.chi", Platform: "kis",      Symbol: "gotrue",    Type: "future"},
	}

	t.Run("InsertProducts", func(t *testing.T) {
		err := client.InsertProducts(context.Background(), products)
		assert.NoError(t, err)
	})

	t.Run("GetProduct", func(t *testing.T) {
		for _, p := range products {
			product, err := client.GetProduct(context.Background(), p.ID)
			assert.NoError(t, err)
			assert.Equal(t, p, product)
		}
	})

	t.Run("DeleteProduct", func(t *testing.T) {
		err := client.DeleteProduct(context.Background(), products[2].ID)
		assert.NoError(t, err)

		_, err = client.GetProduct(context.Background(), products[2].ID)
		assert.Error(t, err)
	})

	t.Run("InsertOneProduct", func(t *testing.T) {
		err := client.InsertOneProduct(context.Background(), products[2])
		assert.NoError(t, err)

		p, err := client.GetProduct(context.Background(), products[2].ID)
		assert.NoError(t, err)
		assert.Equal(t, products[2], p)
	})

	t.Run("InsertProductAlreadyExists", func(t *testing.T) {
		err := client.InsertOneProduct(context.Background(), products[2])
		assert.Error(t, err)
	})

	t.Run("GetAllProducts", func(t *testing.T) {
		ps, err := client.GetAllProducts(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, len(products), len(ps))
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
