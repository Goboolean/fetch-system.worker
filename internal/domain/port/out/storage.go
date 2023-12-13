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
	GetAllWorker(ctx context.Context) ([]vo.Worker, error)
	RegisterWorker(ctx context.Context, worker vo.Worker) error
	UpdateWorkerStatus(ctx context.Context, workerId string, status vo.WorkerStatus) error
	UpdateWorkerStatusExited(ctx context.Context, workerId string, status vo.WorkerStatus, timestamp time.Time) error
	DeleteWorker(ctx context.Context, workerId string) error
	CreateConnection(ctx context.Context, workerId string) (chan struct{}, error)
	WatchConnectionEnds(ctx context.Context, workerId string) (chan struct{}, error)
	WatchPromotion(ctx context.Context, workerId string) (chan struct{}, error)
	GetAllProducts(ctx context.Context) ([]vo.Product, error)
}