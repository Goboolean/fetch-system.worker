package etcd

import (
	"context"

	etcdutil "github.com/Goboolean/fetch-system.worker/internal/infrastructure/etcd/util"
	"go.etcd.io/etcd/client/v3"
)

func (c *Client) InsertWorker(ctx context.Context, w *Worker) error {
	payload, err := etcdutil.Serialize(w)

	var ops []clientv3.Op
	for k, v := range payload {
		ops = append(ops, clientv3.OpPut(k, v))
	}

	var conditions []clientv3.Cmp
	for k := range payload {
		conditions = append(conditions, clientv3.Compare(clientv3.Version(k), "=", 0))
	}

	resp, err := c.client.Txn(ctx).
		If(conditions...).
		Then(ops...).
		Commit()

	if err != nil {
		return err
	}
	if flag := resp.Succeeded; !flag {
		return ErrObjectExists
	}
	return err
}

func (c *Client) GetWorker(ctx context.Context, id string) (*Worker, error) {

	resp, err := c.client.Get(context.Background(), etcdutil.Identifier("worker", id), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, ErrWorkerNotExists
	}

	m := etcdutil.PayloadToMap(resp)

	var w Worker
	if err := etcdutil.Deserialize(m, &w); err != nil {
		return nil, err
	}
	return &w, nil
}

func (c *Client) GetWorkerTimestamp(ctx context.Context, id string) (string, error) {
	resp, err := c.client.Get(ctx, etcdutil.Field("worker", id, "timestamp"))
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", ErrWorkerNotExists
	}
	return string(resp.Kvs[0].Value), nil
}

func (c *Client) GetWorkerStatus(ctx context.Context, id string) (string, error) {
	resp, err := c.client.Get(ctx, etcdutil.Field("worker", id, "status"))
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", ErrWorkerNotExists
	}
	return string(resp.Kvs[0].Value), nil
}

func (c *Client) DeleteWorker(ctx context.Context, id string) error {
	_, err := c.client.Delete(context.Background(), etcdutil.Identifier("worker", id), clientv3.WithPrefix())
	return err
}

func (c *Client) DeleteAllWorkers(ctx context.Context) error {
	_, err := c.client.Delete(context.Background(), etcdutil.Group("worker"), clientv3.WithPrefix())
	return err
}

func (c *Client) WorkerExists(ctx context.Context, id string) (bool, error) {
	resp, err := c.client.Get(ctx, etcdutil.Identifier("worker", id), clientv3.WithPrefix())
	if err != nil {
		return false, err
	}
	return len(resp.Kvs) != 0, nil
}

func (c *Client) GetAllWorkers(ctx context.Context) ([]*Worker, error) {

	resp, err := c.client.Get(context.Background(), etcdutil.Group("worker"), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	m := etcdutil.PayloadToMap(resp)

	list, err := etcdutil.GroupByPrefix(m)
	if err != nil {
		return nil, err
	}

	var w []*Worker = make([]*Worker, len(list))
	for i, v := range list {
		var worker Worker
		if err := etcdutil.Deserialize(v, &worker); err != nil {
			return nil, err
		}
		w[i] = &worker
	}
	return w, nil
}

func (c *Client) UpdateWorkerStatus(ctx context.Context, id string, status string) error {
	_, err := c.client.Put(context.Background(), etcdutil.Field("worker", id, "status"), status)
	return err
}

func (c *Client) UpdateWorkerStatusExited(ctx context.Context, id string, status string, timestamp string) error {
	_, err := c.client.Txn(ctx).
		Then([]clientv3.Op{
			clientv3.OpPut(etcdutil.Field("worker", id, "status"),    status),
			clientv3.OpPut(etcdutil.Field("worker", id, "timestamp"), timestamp),
		}...).
		Commit()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) InsertOneProduct(ctx context.Context, p *Product) error {
	payload, err := etcdutil.Serialize(p)

	var ops []clientv3.Op
	for k, v := range payload {
		ops = append(ops, clientv3.OpPut(k, v))
	}

	var conditions []clientv3.Cmp
	for k := range payload {
		conditions = append(conditions, clientv3.Compare(clientv3.Version(k), "=", 0))
	}

	resp, err := c.client.Txn(ctx).
		If(conditions...).
		Then(ops...).
		Commit()

	if err != nil {
		return err
	}
	if flag := resp.Succeeded; !flag {
		return ErrObjectExists
	}
	return nil
}

func (c *Client) InsertProducts(ctx context.Context, p []*Product) error {

	var conditions []clientv3.Cmp
	var ops []clientv3.Op

	for _, v := range p {
		payload, err := etcdutil.Serialize(v)
		if err != nil {
			return err
		}
		for k, v := range payload {
			ops = append(ops, clientv3.OpPut(k, v))
			conditions = append(conditions, clientv3.Compare(clientv3.Version(k), "=", 0))
		}
	}

	resp, err := c.client.Txn(ctx).
		If(conditions...).
		Then(ops...).
		Commit()

	if err != nil {
		return err
	}
	if flag := resp.Succeeded; !flag {
		return ErrObjectExists
	}
	return nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {

	resp, err := c.client.Get(context.Background(), etcdutil.Identifier("product", id), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, ErrProductNotExists
	}

	m := etcdutil.PayloadToMap(resp)

	var p Product
	if err := etcdutil.Deserialize(m, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *Client) GetAllProducts(ctx context.Context) ([]*Product, error) {

	resp, err := c.client.Get(context.Background(), etcdutil.Group("product"), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	m := etcdutil.PayloadToMap(resp)

	list, err := etcdutil.GroupByPrefix(m)
	if err != nil {
		return nil, err
	}

	var p []*Product = make([]*Product, len(list))
	for i, v := range list {
		var product Product
		if err := etcdutil.Deserialize(v, &product); err != nil {
			return nil, err
		}
		p[i] = &product
	}
	return p, nil
}

func (c *Client) GetProductsWithCondition(ctx context.Context, platform string, market string, locale string) ([]*Product, error) {

	resp, err := c.client.Get(context.Background(), etcdutil.Group("product"), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	m := etcdutil.PayloadToMap(resp)

	list, err := etcdutil.GroupByPrefix(m)
	if err != nil {
		return nil, err
	}

	var p []*Product = make([]*Product, 0)
	for i, v := range list {
		var product Product
		if err := etcdutil.Deserialize(v, &product); err != nil {
			return nil, err
		}

		if product.Platform != "" && product.Platform != platform {
			continue
		}
		if product.Market != "" && product.Market != market {
			continue
		}
		if product.Locale != "" && product.Locale != locale {
			continue
		}

		p[i] = &product
	}
	return p, nil
}


func (c *Client) DeleteProduct(ctx context.Context, id string) error {
	_, err := c.client.Delete(context.Background(), etcdutil.Identifier("product", id), clientv3.WithPrefix())
	return err
}

func (c *Client) DeleteAllProducts(ctx context.Context) error {
	_, err := c.client.Delete(context.Background(), etcdutil.Group("product"), clientv3.WithPrefix())
	return err
}