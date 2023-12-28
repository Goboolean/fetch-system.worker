package otel

import (
	"go.opentelemetry.io/otel/metric"
)



var (
	KafkaProducerSuccessCount metric.Int64Counter
	KafkaProducerErrorCount   metric.Int64Counter
	ETCDErrorCount            metric.Int64Counter
	MockGeneratorCount		  metric.Int64Counter
	PolygonStockErrorCount    metric.Int64Counter
	PolygonStockReceivedCount metric.Int64Counter
	KISStockErrorCount        metric.Int64Counter
	KISStockReceivedCount     metric.Int64Counter
	ProductSubscribedCount    metric.Int64Counter
	ProductReceivedCount      metric.Int64Counter
)



func initMetric(meter metric.Meter) (err error) {
	KafkaProducerSuccessCount, err = meter.Int64Counter("kafka.producer.success.count")
	if err != nil {
		return
	}

	KafkaProducerErrorCount, err = meter.Int64Counter("kafka.producer.error.count")
	if err != nil {
		return
	}

	ETCDErrorCount, err = meter.Int64Counter("etcd.error.count")
	if err != nil {
		return
	}

	MockGeneratorCount, err = meter.Int64Counter("mock.generator.count")
	if err != nil {
		return
	}

	PolygonStockErrorCount, err = meter.Int64Counter("polygon.stock.error.count")
	if err != nil {
		return
	}

	PolygonStockReceivedCount, err = meter.Int64Counter("polygon.stock.received.count")
	if err != nil {
		return
	}

	KISStockErrorCount, err = meter.Int64Counter("kis.stock.error.count")
	if err != nil {
		return
	}
	
	KISStockReceivedCount, err = meter.Int64Counter("kis.stock.received.count")
	if err != nil {
		return
	}

	ProductSubscribedCount, err = meter.Int64Counter("product.subscribed.count")
	if err != nil {
		return
	}

	ProductReceivedCount, err = meter.Int64Counter("product.received.count")
	if err != nil {
		return
	}

	return
}
