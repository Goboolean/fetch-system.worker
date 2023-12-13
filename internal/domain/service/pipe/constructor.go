package pipe

import (
	"context"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/domain/port/out"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
)




type Manager struct {
	d out.DataDispatcher
	f out.DataFetcher
}

func New() *Manager {
	return &Manager{}
}

func (m *Manager) Close(ctx context.Context) error {
	return nil
}

func (m *Manager) RunStreamingPipe(ctx context.Context, products []*vo.Product) error {
	return nil
}

func (m *Manager) RunStoringPipe(ctx context.Context, products []*vo.Product) error {
	return nil
}


func (m *Manager) LockUpStoringPipe(ctx context.Context, timestamp time.Time) error {
	return nil
}

func (m *Manager) UpgradeToStreamingPipe(ctx context.Context, timestamp time.Time) error {
	return nil
}