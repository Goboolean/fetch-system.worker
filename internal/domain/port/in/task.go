package in

import "context"



type TaskCommander interface {
	RegisterWorker(ctx context.Context) error
	Shutdown() error
}