package vo



type WorkerStatus string

const (
	WorkerStatusPrimary               WorkerStatus = "Primary"
	WorkerStatusSecondary             WorkerStatus = "Secondary"
	WorkerStatusExited                WorkerStatus = "Exited"

	WorkerStatusExitedShutdownOccured WorkerStatus = "ExitedShutdownOccured"
	WorkerStatusExitedTTlFailed       WorkerStatus = "ExitedTtlFailed"
	WorkerStatusExitedRegisterFailed  WorkerStatus = "ExitedRegisterFailed"
	WorkerStatusExitedUnknownError    WorkerStatus = "ExitedUnknownError"
	WorkerStatusUnknown               WorkerStatus = "Unknown"
)


type Platform string

const (
	PlatformKIS      Platform = "KIS"
	PlatformPolygon  Platform = "Polygon"
)

func (wp Platform) String() string {
	switch wp {
	case PlatformKIS:
		return "KIS"
	case PlatformPolygon:
		return "Polygon"
	default:
		return "Unknown"
	}
}