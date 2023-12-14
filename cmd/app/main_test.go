package main_test

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/Goboolean/fetch-system.worker/cmd/inject"
	"github.com/Goboolean/fetch-system.worker/internal/domain/port/out"
	"github.com/Goboolean/fetch-system.worker/internal/domain/service/pipe"
	"github.com/Goboolean/fetch-system.worker/internal/domain/service/task"
	"github.com/stretchr/testify/assert"

	_ "github.com/Goboolean/common/pkg/env"
)





func TestMainScenario(t *testing.T) {
	var err error

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var (
		kafka out.DataDispatcher
		etcd out.StorageHandler
		fetcher out.DataFetcher

		pipeManager *pipe.Manager
		taskManager *task.Manager
	)

	t.Cleanup(func() {
		err := etcd.DeleteAllWorkers(ctx)
		assert.NoError(t, err)
	})

	t.Run("Run kafka", func(tt *testing.T) {
		var cleanup func()

		kafka, cleanup, err = inject.InitializeKafkaProducer()
		assert.NoError(tt, err)

		t.Cleanup(cleanup)
	})

	t.Run("Run etcd", func(tt *testing.T) {
		var cleanup func()

		etcd, cleanup, err = inject.InitializeETCDClient()
		assert.NoError(tt, err)

		t.Cleanup(cleanup)
	})

	t.Run("Run fetcher", func(tt *testing.T) {
		var cleanup func()

		fetcher, cleanup, err = inject.InitializeFetcher()
		assert.NoError(tt, err)

		t.Cleanup(cleanup)
	})

	t.Run("Run pipe manager", func(t *testing.T) {
		var err error
		pipeManager, err = inject.InitializePipeManager(kafka, fetcher)
		assert.NoError(t, err)
	})

	t.Run("Run task manager", func(t *testing.T) {
		var err error
		taskManager, err = inject.InitializeTaskManager(pipeManager, etcd)
		assert.NoError(t, err)
	})

	t.Run("Register worker", func(t *testing.T) {
		err := taskManager.RegisterWorker(ctx); 
		assert.NoError(t, err)
	})
}