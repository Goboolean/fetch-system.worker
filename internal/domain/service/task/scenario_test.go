package task_test

import (
	"context"
	"testing"

	"github.com/Goboolean/fetch-system.worker/internal/domain/vo"
	"github.com/stretchr/testify/assert"
)



func TestTwoRegister(t *testing.T) {
	w1 := vo.Worker{
		ID: "first",
		Platform: vo.PlatformKIS,
	}

	w2 := vo.Worker{
		ID: "second",
		Platform: vo.PlatformKIS,
	}

	m1 := SetupTaskManager(&w1)
	m2 := SetupTaskManager(&w2)

	t.Run("RegisterPrimary", func(t *testing.T) {
		err := m1.RegisterWorker(context.Background())
		assert.NoError(t, err)

		assert.Equal(t, vo.WorkerStatusPrimary, w1.Status)
		workers, err := etcdStub.GetAllWorker(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 1, len(workers))
	})

	t.Run("RegisterSecondary", func(t *testing.T) {
		err := m2.RegisterWorker(context.Background())
		assert.NoError(t, err)

		assert.Equal(t, vo.WorkerStatusSecondary, w2.Status)	
	})
}