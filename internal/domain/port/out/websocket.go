package out

import (
	"context"

	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
)

type DataFetcher interface {
	Subscribe(ctx context.Context, symbols ...string) (<-chan vo.Trade, error)
}

type DataDispatcher interface {
	GetPipe() chan<- vo.Trade
}