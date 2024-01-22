package vo

import "time"


type Product struct {
	Symbol   string
	ID       string
	Platform Platform
	Market   Market
	Locale   Locale
}


type Trade struct {
	TradeDetail
	ID     string
	Symbol string
}

type TradeDetail struct {
	Price     float64
	Size      int64
	Timestamp time.Time
}

type Worker struct {
	ID       string
	Status   WorkerStatus
	Platform Platform
	Market   Market
}