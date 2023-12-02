package etcd

import "errors"



var ErrProductNotExists = errors.New("product not exists")

var ErrWorkerNotExists = errors.New("worker not exists")

var ErrObjectExists = errors.New("object exists")