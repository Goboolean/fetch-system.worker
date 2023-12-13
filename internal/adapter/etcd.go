package adapter

import (
	"context"
	"strconv"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/domain/port/out"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/etcd"
	"go.etcd.io/etcd/client/v3"
)



type ETCDAdapter struct {
	etcd *etcd.Client

	leaseID clientv3.LeaseID
}

func NewETCDAdapter(etcd *etcd.Client) out.StorageHandler {
	return &ETCDAdapter{etcd: etcd}
}



func (a ETCDAdapter) Mutex(ctx context.Context, key out.MutexKey) (out.Mutex, error) {
	mu, err := a.etcd.NewMutex(ctx, key.String())
	return mu, err
}

func (a ETCDAdapter) GetAllWorker(ctx context.Context) ([]vo.Worker, error) {
	workerList, err := a.etcd.GetAllWorkers(ctx)
	if err != nil {
		return nil, err
	}
	
	workers := make([]vo.Worker, len(workerList))

	for worker := range workers {
		workers[worker] = vo.Worker{
			ID: workerList[worker].ID,
			Status: vo.WorkerStatus(workerList[worker].Status),
			Platform: vo.Platform(workerList[worker].Platform),
		}
	}

	return workers, nil
}

func (a ETCDAdapter) RegisterWorker(ctx context.Context, worker vo.Worker) error {
	return a.etcd.InsertWorker(ctx, &etcd.Worker{
		ID: worker.ID,
		Status: string(worker.Status),
		Platform: string(worker.Platform),
		LeaseID: strconv.Itoa(int(a.leaseID)),
	})
}

func (a ETCDAdapter) GetWorker(ctx context.Context, workerID string) (*vo.Worker, error) {
	worker, err := a.etcd.GetWorker(ctx, workerID)
	if err != nil {
		return nil, err
	}

	return &vo.Worker{
		ID: worker.ID,
		Status: vo.WorkerStatus(worker.Status),
		Platform: vo.Platform(worker.Platform),
	}, nil
}

func (a ETCDAdapter) GetWorkerTimestamp(ctx context.Context, workerID string) (time.Time, error) {
	timestampStr, err := a.etcd.GetWorkerTimestamp(ctx, workerID)
	if err != nil {
		return time.Time{}, err
	}

	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		return time.Time{}, err
	}
	return timestamp, nil
}

func (a ETCDAdapter) GetWorkerStatus(ctx context.Context, workerID string) (vo.WorkerStatus, error) {
	status, err := a.etcd.GetWorkerStatus(ctx, workerID)
	if err != nil {
		return "", err
	}
	return vo.WorkerStatus(status), nil
}

func (a ETCDAdapter) UpdateWorkerStatus(ctx context.Context, workerId string, status vo.WorkerStatus) error {
	return a.etcd.UpdateWorkerStatus(ctx, workerId, string(status))
}

func (a ETCDAdapter) UpdateWorkerStatusExited(ctx context.Context, workerId string, status vo.WorkerStatus, timestamp time.Time) error {
	return a.etcd.UpdateWorkerStatusExited(ctx, workerId, string(status), timestamp.String())
}

func (a ETCDAdapter) DeleteWorker(ctx context.Context, workerId string) error {
	return a.etcd.DeleteWorker(ctx, workerId)
}

func (a ETCDAdapter) WatchConnectionEnds(ctx context.Context, workerId string) (chan struct{}, error) {

	ch := make(chan struct{})
	respCh := a.etcd.WatchWorkerStatus(ctx, workerId)

	go func() {
		for range respCh {
			ch <- struct{}{}
		}
	}()
	return ch, nil
}

func (a ETCDAdapter) WatchPromotion(ctx context.Context, workerId string) (chan struct{}, error) {
	
	ch := make(chan struct{})
	respCh := a.etcd.WatchWorkerStatus(ctx, workerId)

	go func() {
		for range respCh {
			ch <- struct{}{}
		}
	}()
	return ch, nil
}

func (a ETCDAdapter) GetAllProducts(ctx context.Context) ([]vo.Product, error) {
	productList, err := a.etcd.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]vo.Product, len(productList))

	for product := range products {
		products[product] = vo.Product{
			Symbol: productList[product].Symbol,
			ID: productList[product].ID,
			Platform: vo.Platform(productList[product].Platform),
		}
	}

	return products, nil
}

func (a ETCDAdapter) CreateConnection(ctx context.Context, workerId string) (chan struct{}, error) {
	ch := make(chan struct{})
	respCh, leaseID, err := a.etcd.CreateKeepAlive(ctx)
	if err != nil {
		return nil, err
	}

	a.leaseID = leaseID

	go func() {
		for range respCh {
			ch <- struct{}{}
		}
	}()
	return ch, nil
}

func (a ETCDAdapter) Shutdown(ctx context.Context) error {
	return a.etcd.ShutdownKeepAlive(ctx, a.leaseID)
}