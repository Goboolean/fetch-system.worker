package etcd_test

import (
	"context"
	"os"
	"testing"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/etcd"
	"github.com/stretchr/testify/assert"
)



var client *etcd.Client

func Setup() *etcd.Client {
	c, err := etcd.New(&resolver.ConfigMap{
		"HOST":      os.Getenv("ETCD_HOST"),
		"PEER_HOST": os.Getenv("ETCD_PEER_HOST"),
	})
	if err != nil {
		panic(err)
	}

	return c
}

func Teardown(c *etcd.Client) {
	if err := c.Cleanup(); err != nil {
		panic(err)
	}
	if err := c.Close(); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	client = Setup()

	code := m.Run()
	Teardown(client)

	os.Exit(code)
}

func Test_Constructor(t *testing.T) {
	t.Run("Ping", func(t *testing.T) {
		err := client.Ping(context.Background())
		assert.NoError(t, err)
	})
}