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
	"github.com/Goboolean/fetch-system.worker/internal/util/otel/production"
	_ "github.com/Goboolean/fetch-system.worker/internal/util/otel/production"
)






func main() {
	log.Info("Application Running...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

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

	ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	defer production.Close(ctx)

	if err := taskManager.RegisterWorker(ctx); err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("Panic recovered: ", r)
			if err := taskManager.Shutdown(); err != nil {
				log.WithField("error", err).Error("Failed to shutdown task manager")
			}
			cancel()
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("Application Shutdown...")
		if err := taskManager.Shutdown(); err != nil {
			log.WithField("error", err).Error("Failed to shutdown task manager")
		}
		return
	case <-taskManager.OnConnectionFailed():
		log.Error("Connection failed...")
		if err := taskManager.Cease(); err != nil {
			log.WithField("error", err).Error("Failed to cease task manager")
		}
		return
	}	
}