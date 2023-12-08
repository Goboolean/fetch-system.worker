package out

import "github.com/Goboolean/fetch-system.worker/internal/domain/vo"

type DataFetcher interface {
	Subscribe(symbols ...string) (<-chan vo.Trade, error)
}

type DataDispatcher interface {
	GetPipe() chan<- vo.Trade
}