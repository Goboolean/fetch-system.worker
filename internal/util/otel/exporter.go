package otel

import (
	"context"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type Wrapper struct {
	close func(context.Context) error
}

func New(c *resolver.ConfigMap) (*Wrapper, error) {
	host, err := c.GetStringKey("OTEL_ENDPOINT")
	if err != nil {
		return nil, err
	}

	close, err := initProvider(context.Background(), host)
	if err != nil {
		return nil, err
	}

	w := &Wrapper{close: close}
	return w, nil
}



func (w *Wrapper) Close() {
	w.close(context.Background())
}



func initProvider(ctx context.Context, host string) (func(context.Context) error, error) {

	res, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("fetch-system.worker"),
			semconv.ServiceVersion("0.1.0"),
		))
	if err != nil {
		return nil, err
	}

	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithEndpoint(host))
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter,
			metric.WithInterval(1 * time.Second))))

	meterProvider.Meter("example")

	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown, nil
}