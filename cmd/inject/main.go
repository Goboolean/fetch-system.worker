//go:build wireinject
// +build wireinject

package inject

import (
	"fmt"
	"os"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.worker/internal/adapter"
	"github.com/Goboolean/fetch-system.worker/internal/domain/port/out"
	"github.com/Goboolean/fetch-system.worker/internal/domain/service/pipe"
	"github.com/Goboolean/fetch-system.worker/internal/domain/service/task"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/etcd"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/kafka"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/polygon"
	"github.com/google/wire"
)



func ProvideKafkaConfig() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
	}
}

func ProvideMockGenerator() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"MODE":               "BASIC",
		"STANDARD_DEVIATION": 100,
		"PRODUCT_COUNT":      1000,
		"PRODUCTION_RATE":    100,
		"BUFFER_SIZE":        100000,
	}
}


func ProvideKISConfig() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"APPKEY": os.Getenv("KIS_APPKEY"),
		"SECRET": os.Getenv("KIS_SECRET"),
	}
}

func ProvidePolygonConfig() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"API_KEY": os.Getenv("POLYGON_API_KEY"),
	}
}

func ProvideETCDConfig() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"HOST": os.Getenv("ETCD_HOST"),
	}
}

func ProvideWorkerConfig() *vo.Worker {
	return &vo.Worker{
		ID: os.Getenv("WORKER_ID"),
		Platform: vo.Platform(os.Getenv("PLATFORM")),
	}
}



func ProvideKafkaProducer(c *resolver.ConfigMap) (*kafka.Producer, func(), error) {
    producer, err := kafka.NewProducer(c)
    if err != nil {
        return nil, nil, err
    }
    return producer, func() {
        producer.Close()
    }, nil
}

func ProvideETCDClient(c *resolver.ConfigMap) (*etcd.Client, func(), error) {
	client, err := etcd.New(c)
	if err != nil {
		return nil, nil, err
	}
	return client, func() {
		client.Close()
	}, nil
}

func ProvidePolygonStockClient(c *resolver.ConfigMap) (*polygon.StocksClient, func(), error) {
	client, err := polygon.NewStocksClient(c)
	if err != nil {
		return nil, nil, err
	}
	return client, func() {
		client.Close()
	}, nil
}

func ProvidePolygonOptionClient(c *resolver.ConfigMap) (*polygon.OptionClient, func(), error) {
	client, err := polygon.NewOptionClient(c)
	if err != nil {
		return nil, nil, err
	}
	return client, func() {
		client.Close()
	}, nil
}

func ProvidePolygonCryptoClient(c *resolver.ConfigMap) (*polygon.CryptoClient, func(), error) {
	client, err := polygon.NewCryptoClient(c)
	if err != nil {
		return nil, nil, err
	}
	return client, func() {
		client.Close()
	}, nil
}



func InitializeKafkaProducer() (out.DataDispatcher, func(), error) {
	wire.Build(
		ProvideKafkaConfig,
		ProvideKafkaProducer,
		adapter.NewKafkaAdapter,
	)
	return nil, nil, nil
}

func InitializeETCDClient() (out.StorageHandler, func(), error) {
	wire.Build(
		ProvideETCDConfig,
		ProvideETCDClient,
		adapter.NewETCDAdapter,		
	)
	return nil, nil, nil
}

func InitializePolygonStockClient() (out.DataFetcher, func(), error) {
	wire.Build(
		ProvidePolygonConfig,
		ProvidePolygonStockClient,
		adapter.NewStockPolygonAdapter,
	)
	return nil, nil, nil
}

func InitializePolygonOptionClient() (out.DataFetcher, func(), error) {
	wire.Build(
		ProvidePolygonConfig,
		ProvidePolygonOptionClient,
		adapter.NewOptionPolygonAdapter,
	)
	return nil, nil, nil
}

func InitializePolygonCryptoClient() (out.DataFetcher, func(), error) {
	wire.Build(
		ProvidePolygonConfig,
		ProvidePolygonCryptoClient,
		adapter.NewCryptoPolygonAdapter,
	)
	return nil, nil, nil
}




func InitializeFetcher() (out.DataFetcher, func(), error) {
	switch os.Getenv("PLATFORM") {
	case "POLYGON":
		switch os.Getenv("PRODUCT") {
			case "STOCK":
				return InitializePolygonStockClient()
			case "OPTION":
				return InitializePolygonOptionClient()
			case "CRYPTO":
				return InitializePolygonCryptoClient()
			default:
				return nil, nil, fmt.Errorf("invalid product")
			}
	case "KIS":
		return nil, nil, fmt.Errorf("not implemented")
	default:
		return nil, nil, fmt.Errorf("invalid platform")
	}
}

func InitializeETCD() (*etcd.Client, error) {
	wire.Build(
		ProvideETCDConfig,
		etcd.New,
	)
	return nil, nil
}

func InitializePipeManager(out.DataDispatcher, out.DataFetcher) (*pipe.Manager, error) {
	wire.Build(
		pipe.New,
	)
	return nil, nil
}

func InitializeTaskManager(pipe.Handler, out.StorageHandler) (*task.Manager, error) {
	wire.Build(
		ProvideWorkerConfig,
		task.New,
	)
	return nil, nil
}