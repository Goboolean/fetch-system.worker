package vo



type WorkerStatus int

const (
	WorkerStatusActive WorkerStatus = iota
	WorkerStatusPending
)


type Worker struct {
	ID    string
	Status WorkerStatus
}