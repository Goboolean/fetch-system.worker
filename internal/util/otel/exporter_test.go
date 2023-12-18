package otel_test

import (
	"os"
	"testing"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.worker/internal/util/otel"

	_ "github.com/Goboolean/common/pkg/env"
)



func SetupExporter() *otel.Wrapper {
	w, err := otel.New(&resolver.ConfigMap{
		"OTEL_ENDPOINT": os.Getenv("OTEL_ENDPOINT"),
	})

	if err != nil {
		panic(err)
	}
	return w
}


func TestConstructor(t *testing.T) {
	w := SetupExporter()
	defer w.Close()
}