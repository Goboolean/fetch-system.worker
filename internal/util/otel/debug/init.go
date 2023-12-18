package debug

import (
	"context"

	"github.com/Goboolean/fetch-system.worker/internal/util/otel"
)



func init() {
	var err error
	close, err = otel.InitSTDMeter()
	if err != nil {
		panic(err)
	}
	
}



var close func(context.Context) error

func Close(ctx context.Context) error {
	return close(ctx)
}