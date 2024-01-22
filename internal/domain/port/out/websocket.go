package out

import (
	"context"

	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
)

type DataFetcher interface {
	InputStream(ctx context.Context, symbols ...string) (<-chan *vo.Trade, error)
}

type DataDispatcher interface {
	OutputStream(<-chan *vo.Trade) error
}