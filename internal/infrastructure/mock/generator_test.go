package mock_test

import (
	"context"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/mock"
	"github.com/stretchr/testify/assert"
)



func TestRunGenerator(t *testing.T) {
	
	const (
		symbol = "MOCK"
		duration = time.Millisecond * 1

		count = 1000
		deadline = time.Second * 1
	)

	ch := make(chan *mock.Trade, count * 2)

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	// runs until the deadline
	go mock.RunGenerator(ctx, symbol, duration, ch)

	time.Sleep(deadline)
	time.Sleep(time.Millisecond * 10)

	// assures 99% accuracy 
	assert.LessOrEqual(t, int(count * 0.90), len(ch))
	assert.GreaterOrEqual(t, int(count * 1), len(ch))
}