package adapter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/domain/port/out"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
)



// DO NOT USE UNTIL ETCD FAILS BY FATAL ERROR

type ETCDStub struct {
	m sync.Mutex

	workerList []vo.Worker
	productList []vo.Product

	workerTimestamp map[string]time.Time
}

func NewETCDStub() out.StorageHandler {
	return &ETCDStub{}
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
	return s.workerList, nil
}

func (s *ETCDStub) RegisterWorker(ctx context.Context, worker vo.Worker) error {
	s.workerList = append(s.workerList, worker)
	return nil
}

func (s *ETCDStub) UpdateWorkerStatus(ctx context.Context, workerId string, status vo.WorkerStatus) error {
	for worker := range s.workerList {
		if s.workerList[worker].ID == workerId {
			s.workerList[worker].Status = status
			return nil
		}
	}
	return fmt.Errorf("worker not found")
}

func (s *ETCDStub) UpdateWorkerStatusExited(ctx context.Context, workerId string, status vo.WorkerStatus, timestamp time.Time) error {
	for worker := range s.workerList {
		if s.workerList[worker].ID == workerId {
			s.workerList[worker].Status = status
			s.workerTimestamp[workerId] = timestamp
			return nil
		}
	}
	return fmt.Errorf("worker not found")
}

func (s *ETCDStub) DeleteWorker(ctx context.Context, workerId string) error {
	for worker := range s.workerList {
		if s.workerList[worker].ID == workerId {
			s.workerList = append(s.workerList[:worker], s.workerList[worker+1:]...)
			return nil
		}
	}
	return fmt.Errorf("worker not found")
}

func (s *ETCDStub) CreateConnection(ctx context.Context, workerId string) (chan struct{}, error) {
	return make(chan struct{}), nil
}

func (s *ETCDStub) WatchConnectionEnds(ctx context.Context, workerId string) (chan struct{}, error) {
	return make(chan struct{}), nil
}

func (s *ETCDStub) WatchPromotion(ctx context.Context, workerId string) (chan struct{}, error) {
	return make(chan struct{}), nil
}

func (s *ETCDStub) GetAllProducts(ctx context.Context) ([]vo.Product, error) {
	return s.productList, nil
}