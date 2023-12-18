package otel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)



func newProductionResource() (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName("fetch-system.worker"),
			semconv.ServiceVersion("0.1.0"),
		))
}

func newDebugResource() (*resource.Resource, error) {
	return resource.Default(), nil
}


func newSTDExporter() (metric.Exporter, error) {
	return stdoutmetric.New()
}

func newGRPCExporter(ctx context.Context, endpoint string) (metric.Exporter, error) {
	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	return exporter, nil
}


func newMeterProvider(res *resource.Resource, exporter metric.Exporter) *metric.MeterProvider {
	return metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter,
			metric.WithInterval(1 * time.Second)),
		),
	)
}



func InitSTDMeter() (func(ctx context.Context) error, error) {
	res, err := newDebugResource()
	if err != nil {
		return nil, err
	}

	exporter, err := newSTDExporter()
	if err != nil {
		return nil, err
	}

	meterProvider := newMeterProvider(res, exporter)
	otel.SetMeterProvider(meterProvider)

	if err := initMetric(meterProvider.Meter("fetch-system.worker")); err != nil {
		return nil, err
	}

	return func(ctx context.Context) error {
		if err := meterProvider.ForceFlush(ctx); err != nil {
			return err
		}
		if err := meterProvider.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	}, nil
}



func InitGRPCMeter(ctx context.Context, endpoint string) (func(ctx context.Context) error, error) {
	res, err := newProductionResource()
	if err != nil {
		return nil, err
	}

	exporter, err := newGRPCExporter(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	meterProvider := newMeterProvider(res, exporter)
	otel.SetMeterProvider(meterProvider)

	if err := initMetric(meterProvider.Meter("fetch-system.worker")); err != nil {
		return nil, err
	}

	return func(ctx context.Context) error {
		if err := meterProvider.ForceFlush(ctx); err != nil {
			return err
		}
		if err := meterProvider.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	}, nil
}