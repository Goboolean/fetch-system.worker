package out

import (
	"context"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
)


type Mutex interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
	TryLock(ctx context.Context) error
}


type MutexKey int

const (
	MutexKeyWorker MutexKey = iota
	MutexKeyProduct
)

func (mu MutexKey) String() string {
	switch mu {
	case MutexKeyWorker:
		return "worker"
	case MutexKeyProduct:
		return "product"
	default:
		return ""
	}
}



type StorageHandler interface {
	Mutex(ctx context.Context, key MutexKey) (Mutex, error)
	GetAllProducts(ctx context.Context) ([]*vo.Product, error)
	GetProducts(ctx context.Context, platform vo.Platform, market vo.Market) ([]*vo.Product, error)
	GetAllWorker(ctx context.Context) ([]vo.Worker, error)
	GetWorker(ctx context.Context, workerID string) (*vo.Worker, error)
	GetWorkerTimestamp(ctx context.Context, workerID string) (time.Time, error)
	GetWorkerStatus(ctx context.Context, workerID string) (vo.WorkerStatus, error)
	DeleteAllWorkers(ctx context.Context) error
	RegisterWorker(ctx context.Context, worker vo.Worker) error
	UpdateWorkerStatus(ctx context.Context, workerId string, status vo.WorkerStatus) error
	UpdateWorkerStatusExited(ctx context.Context, workerId string, status vo.WorkerStatus, timestamp time.Time) error
	CreateConnection(ctx context.Context, workerId string) (chan struct{}, error)
	WatchConnectionEnds(ctx context.Context, workerId string) (chan struct{}, error)
	WatchPromotion(ctx context.Context, workerId string) (chan struct{}, error)
}