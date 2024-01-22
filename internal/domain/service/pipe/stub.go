package pipe

import (
	"context"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
)




type Stub struct {}

func NewStub() *Stub {
	return &Stub{}
}

func (m *Stub) Close() {}

func (m *Stub) RunStreamingPipe(ctx context.Context, products []*vo.Product) error {
	return nil
}
func (m *Stub) RunStoringPipe(ctx context.Context, products []*vo.Product) error {
	return nil
}

func (m *Stub) LockupPipe(timestamp time.Time) {}

func (m *Stub) UpgradeToStreamingPipe(timestamp time.Time) {}