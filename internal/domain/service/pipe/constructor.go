package pipe

import (
	"context"

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

func (m Manager) RunStreamingPipe(ctx context.Context, products []*vo.Product) error {
	return nil
}

func (m Manager) RunStoringPipe(ctx context.Context, products []*vo.Product) error {
	return nil
}

func (m Manager) Close(ctx context.Context) error {
	return nil
}

func (m Manager) UpgradeToStreamingPipe(ctx context.Context) error {
	return nil
}