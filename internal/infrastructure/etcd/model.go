package etcd



type Worker struct {
	ID       string `etcd:"id"`        // uuid format
	Platform string `etcd:"platform"`  // kis, polygon, buycycle, ...
	Status   string `etcd:"status"`    // active, waiting, dead
	LeaseID  string  `etcd:"lease_id"`  // lease id
}

func (w *Worker) Name() string {
	return "worker"
}



type Product struct {
	ID       string `etcd:"id"`       // product_type.name.region
	Platform string `etcd:"platform"` // kis, polygon, buycycle, ...
	Symbol   string `etcd:"symbol"`   // identifier inside platform
	Type     string `etcd:"type"`     // stock, future, option, ...
}

func (p *Product) Name() string {
	return "product"
}