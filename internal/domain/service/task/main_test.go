package task_test

import (
	"os"
	"testing"

	"github.com/Goboolean/fetch-system.worker/internal/adapter"
	"github.com/Goboolean/fetch-system.worker/internal/domain/port/in"
	"github.com/Goboolean/fetch-system.worker/internal/domain/service/pipe"
	"github.com/Goboolean/fetch-system.worker/internal/domain/service/task"
	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
)


var etcdStub = adapter.NewETCDStub()

func SetupTaskManager(workerConfig *vo.Worker) in.TaskCommander {
	m, err := task.New(workerConfig, etcdStub, pipe.NewStub())
	if err != nil {
		panic(err)
	}
	return m
}

func TeardownTaskManager(m *task.Manager) {
	if err := m.Shutdown(); err != nil {
		panic(err)
	}	
}



func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}