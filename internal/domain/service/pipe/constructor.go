package pipe

import (
	"context"
	"sync"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/domain/port/out"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
	"github.com/Goboolean/fetch-system.worker/internal/util/otel"
)



type PipeStatus = int

const (
	StreamingPipe PipeStatus = iota + 1
	StoringPipe
)

const (
	defaultBufferSize = 500000
	defaultBufferKeepTime = 100 * time.Second
	defaultBufferSleepTime = 10 * time.Millisecond
)


type Manager struct {
	f out.DataFetcher
	d out.DataDispatcher

	t time.Time
	status PipeStatus

	input  <-chan *vo.Trade
	output chan *vo.Trade
	buffer chan *vo.Trade
	temp   chan *vo.Trade
	arbiter *chan *vo.Trade

	bufferSignal chan struct{}
	isBufferEmptyRunning bool
	wg sync.WaitGroup
}

type Handler interface {
	Close()
	LockupPipe(timestamp time.Time)
	RunStoringPipe(ctx context.Context, products []*vo.Product) error
	RunStreamingPipe(ctx context.Context, products []*vo.Product) error
	UpgradeToStreamingPipe(timestamp time.Time)
}

func New(f out.DataFetcher, d out.DataDispatcher) *Manager {
	instance := &Manager{
		f: f,
		d: d,
		buffer: make(chan *vo.Trade, defaultBufferSize),
		temp: make(chan *vo.Trade),
		output: make(chan *vo.Trade),
		bufferSignal: make(chan struct{}),
	}
	instance.arbiter = &instance.temp
	return instance
}

func (m *Manager) Close() {
	if m.isBufferEmptyRunning {
		m.cancelBufferEmpty()		
	}
}


func (m *Manager) connectInputPipe(ctx context.Context, products []*vo.Product) (err error) {
	symbols := make([]string, len(products))
	for i, product := range products {
		symbols[i] = product.Symbol
	}

	m.input, err = m.f.InputStream(ctx, symbols...)

	symbolToID := make(map[string]string)
	for _, product := range products {
		symbolToID[product.Symbol] = product.ID
	}

	received := make(map[string]bool)
	for _, product := range products {
		received[product.Symbol] = false
	}

	go func ()  {
		for v := range m.input {
			v.ID = symbolToID[v.Symbol]
			*m.arbiter <- v

			if m.status == StreamingPipe && !received[v.Symbol] {
				received[v.Symbol] = true
				otel.ProductReceivedCount.Add(ctx, 1)
			}
		}
	}()
	return
}

func (m *Manager) runBufferEmpty() {
	m.wg.Add(1)
	m.isBufferEmptyRunning = true

	defer func () {
		m.isBufferEmptyRunning = false
		m.wg.Done()
	}()

	for v := range m.buffer {
		if time.Now().Sub(v.Timestamp) < defaultBufferKeepTime {
			timer := time.NewTimer(defaultBufferKeepTime - time.Now().Sub(v.Timestamp))

			select {
			case <- timer.C:
				continue
			case <- m.bufferSignal:
				return
			}
		}
	}
}

func (m *Manager) cancelBufferEmpty() {
	m.bufferSignal <- struct{}{}
	m.wg.Wait()
}



func (m *Manager) connectOutputPipe(ctx context.Context) error {
	if err := m.d.OutputStream(m.output); err != nil {
		return err
	}

	return nil
}


func (m *Manager) RunStreamingPipe(ctx context.Context, products []*vo.Product) error {
	m.status = StreamingPipe

	if err := m.connectInputPipe(ctx, products); err != nil {
		return err
	}
	if err := m.connectOutputPipe(ctx); err != nil {
		return err
	}
	m.arbiter = &m.output

	otel.ProductSubscribedCount.Add(ctx, int64(len(products)))
	return nil
}

func (m *Manager) RunStoringPipe(ctx context.Context, products []*vo.Product) error {
	m.status = StoringPipe

	if err := m.connectInputPipe(ctx, products); err != nil {
		return err
	}
	if err := m.connectOutputPipe(ctx); err != nil {
		return err
	}
	m.arbiter = &m.buffer
	go m.runBufferEmpty()
	return nil
}

func (m *Manager) UpgradeToStreamingPipe(timestamp time.Time) {
	m.arbiter = &m.temp
	m.cancelBufferEmpty()

	size := len(m.buffer)
	for i := 0; i < size; i++ {
		v := <- m.buffer
		m.output <- v
	}
	m.arbiter = &m.output
}

func (m *Manager) LockupPipe(timestamp time.Time) {
	timer := time.NewTimer(time.Until(timestamp))

	go func() {
		<- timer.C
		m.arbiter = &m.temp
	}()
}