package etcd

import (
	"context"
	"strconv"

	etcdutil "github.com/Goboolean/fetch-system.worker/internal/infrastructure/etcd/util"
	"go.etcd.io/etcd/client/v3"
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


const defaultTTL = 5
func (c *Client) CreateKeepAlive(ctx context.Context) (<-chan *clientv3.LeaseKeepAliveResponse, clientv3.LeaseID, error) {
	resp, err := c.client.Grant(ctx, defaultTTL)
	if err != nil {
		return nil, 0, err
	}

	ch, err := c.client.KeepAlive(ctx, resp.ID)
	return ch, resp.ID, err
}

func (c *Client) SetKeepAliveValue(ctx context.Context, id string, leaseId clientv3.LeaseID) error {
	_, err := c.client.Put(ctx, etcdutil.Field("worker", id, "lease_id"), strconv.FormatInt(int64(leaseId), 10), clientv3.WithLease(leaseId))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ShutdownKeepAlive(ctx context.Context, id clientv3.LeaseID) error {
	_, err := c.client.Revoke(ctx, id)
	return err
}

func (c *Client) WatchWorkerKeepAlive(ctx context.Context, id string) <-chan clientv3.WatchResponse {
	w := c.client.Watcher
	ch := w.Watch(ctx, etcdutil.Field("worker", id, "lease_id"))
	return ch
}

func (c *Client) WatchWorkerStatus(ctx context.Context, id string) <-chan clientv3.WatchResponse {
	w := c.client.Watcher
	ch := w.Watch(ctx, etcdutil.Field("worker", id, "status"))
	return ch
}