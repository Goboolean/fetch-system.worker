package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Goboolean/fetch-system.worker/cmd/wire"

	_ "github.com/Goboolean/common/pkg/env"
)






func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	kafka, cleanup, err := wire.InitializeKafkaProducer()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	etcd, cleanup, err := wire.InitializeETCDClient()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	fetcher, cleanup, err := wire.InitializeFetcher()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	pipeManager, err := wire.InitializePipeManager(kafka, fetcher)
	if err != nil {
		panic(err)
	}
	defer pipeManager.Close()

	taskManager, err := wire.InitializeTaskManager(pipeManager, etcd)
	if err != nil {
		panic(err)
	}

	if err := taskManager.RegisterWorker(ctx); err != nil {
		panic(err)
	}

	select {
	case <-ctx.Done():
		taskManager.Shutdown()
		return
	case <-taskManager.OnConnectionFailed():
		taskManager.Cease()
		return
	}	
}