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
	PlatformPolygon  Platform = "POLYGON"
	PlatformMock	 Platform = "MOCK"
)


type Market string

const (
	MarketStock  Market = "STOCK"
	MarketCrypto Market = "CRYPTO"
)


type Locale string

const (
	LocaleUSA Locale = "USA"
	LocaleKOR Locale = "KOR"
)