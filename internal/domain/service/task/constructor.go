package task

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/domain/port/out"
	"github.com/Goboolean/fetch-system.worker/internal/domain/service/pipe"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
	"github.com/ahmetb/go-linq/v3"
	"github.com/pkg/errors"
)



var ConfigRetry = 3
var ConfigTimeout = 3 * time.Second



type Manager struct {
	s out.StorageHandler
	p *pipe.Manager

	worker *vo.Worker
	primaryID string

	connCh chan struct{}
	promCh chan struct{}
	ttlCh  chan struct{}

	ctx context.Context
	cancel context.CancelFunc
	wg sync.WaitGroup
}



func New(worker *vo.Worker, s out.StorageHandler, p *pipe.Manager) (*Manager, error) {
	if worker.ID == "" {
		return nil, fmt.Errorf("worker id is empty")
	}

	if worker.Platform == "" {
		return nil, fmt.Errorf("worker platform is empty")
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		s: s,
		p: p,
		worker: worker,

		ctx: ctx,
		cancel: cancel,
	}, nil
}



func (m *Manager) RegisterWorker(ctx context.Context) error {

	mu, err := m.s.Mutex(ctx, out.MutexKeyWorker)
	if err != nil {
		return err
	}

	if err := mu.Lock(ctx); err != nil {
		return err
	}
	defer mu.Unlock(ctx)

	workers, err := m.s.GetAllWorker(ctx)
	if err != nil {
		return err
	}

	var isSecondary = false

	for _, worker := range workers {
		if worker.Platform == m.worker.Platform && worker.Status == vo.WorkerStatusPrimary {
			isSecondary = true
			m.primaryID = worker.ID
			break
		}
	}

	if isSecondary {
		m.worker.Status = vo.WorkerStatusSecondary
	} else {
		m.worker.Status = vo.WorkerStatusPrimary
	}

	connCh, err := m.s.CreateConnection(ctx, m.worker.ID)
	if err != nil {
		return err
	}
	m.connCh = connCh

	if err := m.s.RegisterWorker(ctx, *m.worker); err != nil {
		return err
	}

	products, err := m.s.GetAllProducts(ctx)
	if err != nil {
		return err
	}

	var productsFiltered []*vo.Product
	linq.From(products).WhereT(func(product vo.Product) bool {
		return product.Platform == m.worker.Platform
	}).ToSlice(&productsFiltered)

	m.s.WatchConnectionEnds(ctx, m.worker.ID)

	var pipeErr error
	if !isSecondary {
		pipeErr = m.p.RunStreamingPipe(ctx, productsFiltered)
	} else {
		pipeErr = m.p.RunStoringPipe(ctx, productsFiltered)
	}

	if pipeErr != nil {
		if err := m.s.UpdateWorkerStatus(context.Background(), m.worker.ID, vo.WorkerStatusExitedRegisterFailed); err != nil {
			return errors.Wrap(err, pipeErr.Error())
		}
		return pipeErr
	}

	if isSecondary {

		go m.trackChan()

		m.promCh, err = m.s.WatchPromotion(ctx, m.primaryID)
		if err != nil {
			return err
		}
		m.ttlCh, err = m.s.WatchConnectionEnds(ctx, m.primaryID)
		if err != nil {
			return err
		}		
	}

	return nil
}

func (m *Manager) Promote(ctx context.Context, _type PromotionType) error {

	mu, err := m.s.Mutex(ctx, out.MutexKeyWorker)
	if err != nil {
		return err
	}

	if err := mu.Lock(ctx); err != nil {
		return err
	}
	defer mu.Unlock(ctx)

	var timestamp time.Time

	switch _type {
		case TTLFailed:
			timestamp = time.Now().Add(- time.Second * 5)
			if err := m.s.UpdateWorkerStatusExited(ctx, m.primaryID, vo.WorkerStatusExitedTTlFailed, timestamp); err != nil {
				return err
			}
		case Shutdown:
			timestamp, err = m.s.GetWorkerTimestamp(ctx, m.primaryID)
			if err != nil {
				return errors.Wrap(err, "failed to get primary worker timestamp")
			}
	}

	if err := m.p.UpgradeToStreamingPipe(ctx, timestamp); err != nil {
		return err
	}

	m.worker.Status = vo.WorkerStatusPrimary
	if err := m.s.UpdateWorkerStatus(ctx, m.worker.ID, vo.WorkerStatusPrimary); err != nil {
		return errors.Wrap(err, "failed to update worker status: ")
	}

	close(m.connCh)
	close(m.promCh)
	return nil
}



func (m *Manager) Shutdown() error {

	var timestamp = time.Now().Truncate(time.Second).Add(time.Second)

	if m.worker.Status == vo.WorkerStatusPrimary {
		m.p.LockUpStoringPipe(m.ctx, timestamp)
	}

	if err := m.s.UpdateWorkerStatusExited(m.ctx, m.worker.ID, vo.WorkerStatusExitedShutdownOccured, timestamp); err != nil {
		return err
	}

	m.cancel()
	return nil
}



type PromotionType int

const (
	TTLFailed PromotionType = iota + 1
	Shutdown
)


func (m *Manager) trackChan() {
	m.wg.Add(1)
	defer m.wg.Done()

	for {
		select {
			case <-m.ctx.Done():
				return
			case _, ok := <- m.promCh:
				if ok {
					if err  := m.Promote(m.ctx, Shutdown); err != nil {
						panic(err)
					}
				}
			case _, ok := <-m.ttlCh:
				if ok {
					if err := m.Promote(m.ctx, TTLFailed); err != nil {
						panic(err)
					}
				}
				continue
		}
	}
}