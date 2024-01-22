package main_test

import (
	"context"
	"os"
	"testing"

	"github.com/Goboolean/fetch-system.IaC/pkg/kafka"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/etcd"

	_ "github.com/Goboolean/common/pkg/env"
	"github.com/Goboolean/common/pkg/resolver"
)



func SetupKafkaConsumer() *kafka.Consumer {
	c, err := kafka.NewConsumer(&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
		"GROUP_ID":       "TEST",
	})
	if err != nil {
		panic(err)
	}

	return c
}

func SetupETCD() *etcd.Client {
	c, err := etcd.New(&resolver.ConfigMap{
		"HOST": os.Getenv("ETCD_HOST"),
	})
	if err != nil {
		panic(err)
	}

	return c
}

var e *etcd.Client

func SetupTestEnvironment() {
	e = SetupETCD()

	var products = []*etcd.Product{
		{ID: "stock.TEST1.usa", Platform: "MOCK", Symbol: "TEST1", Market: "STOCK", Locale: "USA"},
		{ID: "stock.TEST2.usa", Platform: "MOCK", Symbol: "TEST2", Market: "STOCK", Locale: "USA"},
		{ID: "stock.TEST3.usa", Platform: "MOCK", Symbol: "TEST3", Market: "STOCK", Locale: "USA"},
	}

	if err := e.UpsertProducts(context.Background(), products); err != nil {
		panic(err)
	}
}

func TeardownTestEnvironment() {
	if err := e.DeleteAllProducts(context.Background()); err != nil {
		panic(err)
	}

	if err := e.DeleteAllWorkers(context.Background()); err != nil {
		panic(err)
	}

	e.Close()
}


func TestMain(m *testing.M) {
	SetupTestEnvironment()
	code := m.Run()
	TeardownTestEnvironment()
	os.Exit(code)
}