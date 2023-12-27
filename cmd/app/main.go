package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Goboolean/fetch-system.worker/cmd/wire"
	log "github.com/sirupsen/logrus"

	_ "github.com/Goboolean/common/pkg/env"
)






func main() {
	log.Info("Application Running...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	kafka, cleanup, err := wire.InitializeKafkaProducer(ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	etcd, cleanup, err := wire.InitializeETCDClient(ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	fetcher, cleanup, err := wire.InitializeFetcher(ctx)
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

	ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := taskManager.RegisterWorker(ctx); err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("Panic recovered: ", r)
			cancel()
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("Application Shutdown...")
		taskManager.Shutdown()
		return
	case <-taskManager.OnConnectionFailed():
		log.Error("Connection failed...")
		taskManager.Cease()
		return
	}	
}