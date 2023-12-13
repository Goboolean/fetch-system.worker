package adapter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/domain/port/out"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
)



type ETCDStub struct {
	m sync.Mutex

	workerList map[string]vo.Worker
	productList []vo.Product

	workerTimestamp map[string]time.Time

	promCh chan struct{}
	connCh chan struct{}
	ttlCh  chan struct{}
}

func NewETCDStub() out.StorageHandler {
	return &ETCDStub{
		workerList: make(map[string]vo.Worker),
		productList: make([]vo.Product, 0),
		workerTimestamp: make(map[string]time.Time),

		promCh: make(chan struct{}),
		connCh: make(chan struct{}),
		ttlCh:  make(chan struct{}),
	}
}

type MutexStup struct {
	m *sync.Mutex
}

func (s *MutexStup) Lock(ctx context.Context) error {
	s.m.Lock()
	return nil
}

func (s *MutexStup) Unlock(ctx context.Context) error {
	s.m.Unlock()
	return nil
}

func (s *MutexStup) TryLock(ctx context.Context) error {
	if flag := s.m.TryLock(); flag {
		return fmt.Errorf("failed to lock")
	}
	return nil
}



func (s *ETCDStub) Mutex(ctx context.Context, key out.MutexKey) (out.Mutex, error) {
	return &MutexStup{m: &s.m}, nil
}

func (s *ETCDStub) GetAllWorker(ctx context.Context) ([]vo.Worker, error) {
	workers := make([]vo.Worker, 0)
	for _, worker := range s.workerList {
		workers = append(workers, worker)
	}
	return workers, nil
}

func (s *ETCDStub) RegisterWorker(ctx context.Context, worker vo.Worker) error {
	s.workerList[worker.ID] = worker
	return nil
}

func (s *ETCDStub) GetWorker(ctx context.Context, workerID string) (*vo.Worker, error) {
	worker, ok := s.workerList[workerID]
	if !ok {
		return nil, fmt.Errorf("worker not found")
	}
	return &worker, nil
}

func (s *ETCDStub) GetWorkerTimestamp(ctx context.Context, workerID string) (time.Time, error) {
	timestamp, ok := s.workerTimestamp[workerID]
	if !ok {
		return time.Time{}, fmt.Errorf("worker not found%s", workerID)
	}
	return timestamp, nil
}

func (s *ETCDStub) UpdateWorkerStatus(ctx context.Context, workerId string, status vo.WorkerStatus) error {
	worker, ok := s.workerList[workerId]
	if !ok {
		return fmt.Errorf("worker not found%s", workerId)
	}
	worker.Status = status
	s.workerList[workerId] = worker
	return nil	
}

func (s *ETCDStub) UpdateWorkerStatusExited(ctx context.Context, workerId string, status vo.WorkerStatus, timestamp time.Time) error {
	worker, ok := s.workerList[workerId]
	if !ok {
		return fmt.Errorf("worker not found")
	}
	worker.Status = status
	s.workerTimestamp[workerId] = timestamp
	s.workerList[workerId] = worker

	s.promCh <- struct{}{}
	return nil
}

func (s *ETCDStub) DeleteWorker(ctx context.Context, workerId string) error {
	delete(s.workerList, workerId)
	return nil
}

func (s *ETCDStub) CreateConnection(ctx context.Context, workerId string) (chan struct{}, error) {
	return s.connCh, nil
}

func (s *ETCDStub) WatchConnectionEnds(ctx context.Context, workerId string) (chan struct{}, error) {
	return s.connCh, nil
}

func (s *ETCDStub) WatchPromotion(ctx context.Context, workerId string) (chan struct{}, error) {
	return s.promCh, nil
}

func (s *ETCDStub) GetAllProducts(ctx context.Context) ([]vo.Product, error) {
	return s.productList, nil
}

func (s *ETCDStub) Shutdown(ctx context.Context) error {
	s.connCh <- struct{}{}
	s.ttlCh <- struct{}{}
	return nil
}

func (s *ETCDStub) DeleteAllWorkers(ctx context.Context) error {
	s.workerList = make(map[string]vo.Worker)
	return nil
}