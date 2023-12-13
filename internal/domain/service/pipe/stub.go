package pipe

import (
	"context"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
)




type Stub struct {}

func NewStub() *Manager {
	return &Manager{}
}

func (m *Stub) Close() error {
	return nil
}

func (m *Stub) RunStreamingPipe(ctx context.Context, products []*vo.Product) error {
	return nil
}
func (m *Stub) RunStoringPipe(ctx context.Context, products []*vo.Product) error {
	return nil
}


func (m *Stub) LockUpStoringPipe(ctx context.Context, timestamp time.Time) error {
	return nil
}

func (m *Stub) UpgradeToStreamingPipe(ctx context.Context, timestamp time.Time) error {
	return nil
}