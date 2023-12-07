package mock

import "github.com/Goboolean/fetch-system.worker/internal/infrastructure/ws"

// This class is not for product, but for testing.
// any test package can use this class to test the function of receiver.
type MockReceiver struct {
	f func()
}

func NewMockReceiver(f func()) *MockReceiver {
	return &MockReceiver{f: f}
}

func (m *MockReceiver) OnReceiveStockAggs(stock *ws.StockAggregate) error {
	m.f()
	return nil
}
