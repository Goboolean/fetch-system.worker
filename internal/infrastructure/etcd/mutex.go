package etcd

import (
	"context"

	etcdutil "github.com/Goboolean/fetch-system.worker/internal/infrastructure/etcd/util"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func (c *Client) NewMutex(ctx context.Context, key string) (*concurrency.Mutex, error) {
	s, err := concurrency.NewSession(c.client)
	if err != nil {
		return nil, err
	}

	key = etcdutil.Semaphore(key)
	return concurrency.NewMutex(s, key), nil
}
