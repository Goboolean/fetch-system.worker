package adapter

import (
	"github.com/Goboolean/fetch-system.worker/internal/domain/port/out"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/kafka"
	"github.com/Goboolean/fetch-system.IaC/pkg/model"
)



type KafkaAdapter struct {
	p *kafka.Producer
}

func NewKafkaAdapter(p *kafka.Producer) out.DataDispatcher {
	return &KafkaAdapter{p: p}
}

func (a *KafkaAdapter) OutputStream(ch <-chan *vo.Trade) error {
	go func() {
		for trade := range ch {
			a.p.ProduceTrade(trade.ID, &model.Trade{
				Price:     trade.Price,
				Size:      trade.Size,
				Timestamp: trade.Timestamp.Unix(),
			})
		}
	}()
	return nil
}