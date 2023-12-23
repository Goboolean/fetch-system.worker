package kafka_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/pkg/model"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/kafka"
	"github.com/stretchr/testify/assert"
	_ "github.com/Goboolean/fetch-system.worker/internal/util/otel/debug"
	_ "github.com/Goboolean/common/pkg/env"
)

func SetupProducer() *kafka.Producer {
	p, err := kafka.NewProducer(&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
		"TRACER":          "none",
	})
	if err != nil {
		panic(err)
	}
	return p
}

func TeardownProducer(p *kafka.Producer) {
	p.Close()
}

func Test_Producer(t *testing.T) {

	p := SetupProducer()
	defer TeardownProducer(p)

	t.Run("Ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		err := p.Ping(ctx)
		assert.NoError(t, err)
	})
}

func Test_Produce(t *testing.T) {

	p := SetupProducer()
	defer TeardownProducer(p)

	const productId = "test.produce.io"
	const productType = "1s"
	const platform = "mock"
	const market = "stock"

	var now = time.Now().Unix()

	var tradeProtobuf = model.TradeProtobuf{
		Price:     171.55,
		Size:      100,
		Timestamp: now,
	}

	var tradeJson = model.TradeJson{
		Price:     171.55,
		Size:      100,
		Timestamp: now,
	}

	var aggsProtobuf = model.AggregateProtobuf{
		Open:      170.55,
		Closed:    173.55,
		Min:       170.55,
		Max:       171.55,
		Volume:    100,
		Timestamp: time.Now().Unix(),
	}

	var aggsJson = model.AggregateJson{
		Open:      170.55,
		Close:     173.55,
		High:      170.55,
		Low:       171.55,
		Volume:    100,
		Timestamp: time.Now().Unix(),
	}

	t.Run("ProduceProtobufTrade", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()

		err := p.ProduceProtobufTrade(productId, platform, market, &tradeProtobuf)
		assert.NoError(t, err)

		_, err = p.Flush(ctx)
		assert.NoError(t, err)
	})

	t.Run("ProduceProtobufAggs", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()

		err := p.ProduceProtobufAggs(productId, productType, platform, market, &aggsProtobuf)
		assert.NoError(t, err)

		_, err = p.Flush(ctx)
		assert.NoError(t, err)
	})

	t.Run("ProduceJsonTrade", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()

		err := p.ProduceJsonTrade(productId, &tradeJson)
		assert.NoError(t, err)

		_, err = p.Flush(ctx)
		assert.NoError(t, err)
	})

	t.Run("ProduceJsonAggs", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()

		err := p.ProduceJsonAggs(productId, productType, &aggsJson)
		assert.NoError(t, err)

		_, err = p.Flush(ctx)
		assert.NoError(t, err)
	})

}
