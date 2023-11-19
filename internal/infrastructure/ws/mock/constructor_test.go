package mock_test

import (
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/ws"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/ws/mock"
)

var instance ws.Fetcher

var (
	receiver ws.Receiver
	count    int = 0
)

func SetupMock() {
	receiver = mock.NewMockReceiver(func() {
		count++
	})

	instance = mock.New(10*time.Millisecond, receiver)
}

func TeardownMock() {
	instance.Close()
}

func TestMain(m *testing.M) {

	SetupMock()
	code := m.Run()
	TeardownMock()

	os.Exit(code)
}

func Test_Constructor(t *testing.T) {

	t.Run("Ping", func(t *testing.T) {
		if err := instance.Ping(); err != nil {
			t.Errorf("Ping() = %v", err)
			return
		}
	})
}

func Test_SubscribeStockAggs(t *testing.T) {
	// I thought count should be in the interval of [5,20],
	// since the data is generated every 10 ms in average.
	// But it seems that error may occur in some case.

	const (
		symbol        = "TEST"
		invalidSymnol = "INVALID"
	)

	t.Run("Subscribe", func(t *testing.T) {
		if err := instance.SubscribeStockAggs(symbol); err != nil {
			t.Errorf("SubscribeStockAggs() = %v", err)
			return
		}

		time.Sleep(100 * time.Millisecond)

		if !(5 <= count) {
			t.Errorf("count = %v, should be at least 5", count)
			return
		}
	})

	t.Run("SubscribeTwice", func(t *testing.T) {
		if err := instance.SubscribeStockAggs(symbol); err == nil {
			t.Errorf("SubscribeStockAggs() = %v, want %v", err, mock.ErrTopicAlreadyExists)
			return
		}
	})

	t.Run("InvalidUnsubscribe", func(t *testing.T) {
		if err := instance.UnsubscribeStockAggs(invalidSymnol); err == nil {
			t.Errorf("UnsubscribeStockAggs() = %v, want %v", err, mock.ErrTopicNotFound)
			return
		}
	})

	t.Run("Unsubscribe", func(t *testing.T) {

		if err := instance.UnsubscribeStockAggs(symbol); err != nil {
			t.Errorf("UnsubscribeStockAggs() = %v", err)
			return
		}

		countBeforeSubscription := count

		time.Sleep(100 * time.Millisecond)

		countAfterUnsubscription := count
		diff := countAfterUnsubscription - countBeforeSubscription

		if diff != 0 {
			t.Errorf("UnsubscribeStockAggs() received %d, want 0", diff)
			return
		}
	})
}
