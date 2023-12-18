package production

import (
	"context"
	"os"

	"github.com/Goboolean/fetch-system.worker/internal/util/otel"
)



func init() {
	var err error

	endpoint := os.Getenv("OTEL_ENDPOINT")

	close, err = otel.InitGRPCMeter(context.Background(), endpoint)
	if err != nil {
		panic(err)
	}
}



var close func(context.Context) error

func Close(ctx context.Context) error {
	return close(ctx)
}