package pipe

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/domain/port/out"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
)



type PipeStatus = int

const (
	StreamingPipe PipeStatus = iota + 1
	StoringPipe
)

const (
	defaultBufferSize = 500000
	defaultBufferKeepTime = 100 * time.Second
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

	bufferSignal bool
	wg sync.WaitGroup
}

func New(f out.DataFetcher, d out.DataDispatcher) *Manager {
	instance := &Manager{
		f: f,
		d: d,
		buffer: make(chan *vo.Trade, defaultBufferSize),
		temp: make(chan *vo.Trade),
		output: make(chan *vo.Trade),
	}
	instance.arbiter = &instance.temp
	fmt.Println("arbiter initial: ", *instance.arbiter)
	return instance
}

func (m *Manager) Close(ctx context.Context) error {
	return nil
}


func (m *Manager) connectInputPipe(ctx context.Context, products []*vo.Product) (err error) {
	symbols := make([]string, len(products))
	for i, product := range products {
		symbols[i] = product.Symbol
	}

	m.input, err = m.f.InputStream(ctx, symbols...)

	go func ()  {
		for v := range m.input {
			*m.arbiter <- v
		}
	}()
	return
}

func (m *Manager) runBufferEmpty() {
	m.wg.Add(1)
	defer m.wg.Done()

	for v := range m.buffer {
		if m.bufferSignal {
			return
		}

		if time.Now().Sub(v.Timestamp) > defaultBufferKeepTime {
			time.Sleep(time.Second)
			continue
		}
	}
}

func (m *Manager) cancelBufferEmpty() {
	m.bufferSignal = true
	m.wg.Done()
}



func (m *Manager) connectOutputPipe(ctx context.Context) error {
	if err := m.d.OutputStream(m.output); err != nil {
		return err
	}

	return nil
}


func (m *Manager) RunStreamingPipe(ctx context.Context, products []*vo.Product) error {
	if err := m.connectInputPipe(ctx, products); err != nil {
		return err
	}
	if err := m.connectOutputPipe(ctx); err != nil {
		return err
	}
	m.arbiter = &m.output
	fmt.Println("arbiter runstreamingpipe: ", *m.arbiter)

	return nil
}

func (m *Manager) RunStoringPipe(ctx context.Context, products []*vo.Product) error {
	if err := m.connectInputPipe(ctx, products); err != nil {
		return err
	}
	if err := m.connectOutputPipe(ctx); err != nil {
		return err
	}
	m.arbiter = &m.buffer
	m.runBufferEmpty()
	return nil
}

func (m *Manager) UpgradeToStreamingPipe(timestamp time.Time) {
	m.arbiter = &m.temp
	m.cancelBufferEmpty()

	for i := 0; i < len(m.buffer); i++ {
		m.output <- <- m.buffer
	}

	fmt.Println(len(m.temp))
	m.arbiter = &m.output
}

func (m *Manager) LockupPipe(timestamp time.Time) {
	timer := time.NewTimer(time.Until(timestamp))

	go func() {
		<- timer.C
		m.arbiter = &m.temp
		fmt.Println("switched??", *m.arbiter)
	}()
}