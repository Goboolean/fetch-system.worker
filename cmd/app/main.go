package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Goboolean/fetch-system.worker/cmd/inject"

	_ "github.com/Goboolean/common/pkg/env"
)






func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	kafka, cleanup, err := inject.InitializeKafkaProducer()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	etcd, cleanup, err := inject.InitializeETCDClient()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	fetcher, cleanup, err := inject.InitializeFetcher()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	pipeManager, err := inject.InitializePipeManager(kafka, fetcher)
	if err != nil {
		panic(err)
	}
	defer pipeManager.Close()

	taskManager, err := inject.InitializeTaskManager(pipeManager, etcd)
	if err != nil {
		panic(err)
	}

	if err := taskManager.RegisterWorker(ctx); err != nil {
		panic(err)
	}

	fmt.Println("Successfully running")

	select {
	case <-ctx.Done():
		fmt.Println("shutdowning...")
		taskManager.Shutdown()
		return
	case <-taskManager.OnConnectionFailed():
		taskManager.Cease()
		return
	}	
}