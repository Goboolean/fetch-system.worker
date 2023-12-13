package etcd

import (
	"context"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)




type Client struct {
	client *clientv3.Client
}

func New(c *resolver.ConfigMap) (*Client, error) {

	host, err := c.GetStringKey("HOST")
	if err != nil {
		return nil, err
	}

	peer_host, exists, err := c.GetStringKeyOptional("PEER_HOST")
	if err != nil {
		return nil, err
	}

	host_list := []string{host}
	if exists {
		host_list = append(host_list, peer_host)
	}

	config := clientv3.Config{
		Endpoints:   host_list,
		DialTimeout: 5 * time.Second,
	}

	client, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}

	s, err := concurrency.NewSession(client)
	s.Done()

	concurrency.NewLocker(s, "lock")

	return &Client{
		client: client,
	}, nil
}


func (c *Client) Ping(ctx context.Context) error {
	mapi := clientv3.NewMaintenance(c.client)
	_, err := mapi.Status(ctx, c.client.Endpoints()[0])
	return err
}


func (c *Client) Close() error {
	return c.client.Close()
}


func (c *Client) Cleanup() error {
	_, err := c.client.Delete(context.Background(), "", clientv3.WithPrefix())
	return err
}