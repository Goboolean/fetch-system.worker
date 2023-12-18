package pipe_test

import (
	"context"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.worker/internal/adapter"
	"github.com/Goboolean/fetch-system.worker/internal/domain/port/out"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/mock"
)


func SetupMockGenerator(rate int) out.DataFetcher {
	m, err := mock.New(&resolver.ConfigMap{
		"MODE": "BASIC",
		"PRODUCT_COUNT": 1,
		"PRODUCTION_RATE": rate,
		"STANDARD_DEVIATION": 1,
		"BUFFER_SIZE": 1000,
	})
	if err != nil {
		panic(err)
	}

	a, err := adapter.NewMockGeneratorAdapter(m)
	return a
}



type MockReceiver struct {
	ctx context.Context
	cancel context.CancelFunc

	received int
	receivedByID map[string]int
}


func (m *MockReceiver) OutputStream(ch <-chan *vo.Trade) error {
	go func() {
		for {
			select {
			case <-m.ctx.Done():
				return
			case v := <- ch:
				m.received++
				m.receivedByID[v.Symbol]++
				continue
			}
		}
	}()
	return nil	
}

func NewMockReceiver() *MockReceiver {
	ctx, cancel := context.WithCancel(context.Background())
	return &MockReceiver{
		ctx: ctx,
		cancel: cancel,
		receivedByID: make(map[string]int),
	}
}

func (m *MockReceiver) Received() int {
	return m.received
}

func (m *MockReceiver) ReceivedByID(id string) int {
	return m.receivedByID[id]
}

func (m *MockReceiver) Initialize() {
	m.received = 0
}

func (m *MockReceiver) Close() {
	m.cancel()
}
