package otel

import "go.opentelemetry.io/otel/metric"



var (
	KafkaProducerSuccessCount metric.Int64Counter
	KafkaProducerErrorCount metric.Int64Counter
)
