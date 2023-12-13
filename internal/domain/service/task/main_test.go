package task_test

import (
	"context"
	"os"
	"testing"

	"github.com/Goboolean/fetch-system.worker/internal/adapter"
	"github.com/Goboolean/fetch-system.worker/internal/domain/service/pipe"
	"github.com/Goboolean/fetch-system.worker/internal/domain/service/task"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)


var manager *task.Manager

func SetupTaskManager() *task.Manager {

	w := vo.Worker{
		ID: "test",
		Platform: vo.PlatformPolygon,
	}

	m, err := task.New(w, adapter.NewETCDStub(), pipe.New())
	if err != nil {
		panic(err)
	}

	return m
}

func TeardownTaskManager(m *task.Manager) {
	
}


func TestMain(m *testing.M) {
	manager = SetupTaskManager()
	code := m.Run()
	os.Exit(code)
	TeardownTaskManager(manager)
}




func TestTaskManager(t *testing.T) {

	t.Run("Register", func(t *testing.T) {
		err := manager.RegisterWorker(context.Background())
		assert.NoError(t, err)
	})
}