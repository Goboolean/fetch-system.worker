package kis_test

import (
	"os"
	"testing"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.worker/internal/infrastructure/kis"

	_ "github.com/Goboolean/common/pkg/env"
)



var client *kis.Client

func SetupKis() *kis.Client {
	var err error

	c, err := kis.New(&resolver.ConfigMap{
		"APPKEY": os.Getenv("KIS_APPKEY"),
		"SECRET": os.Getenv("KIS_SECRET"),
		"BUFFER_SIZE": 10000,
	})
	if err != nil {
		panic(err)
	}
	return c
}

func TeardownKis(c *kis.Client) {
	c.Close()
}


func TestMain(m *testing.M) {
	client = SetupKis()
	code := m.Run()
	TeardownKis(client)
	os.Exit(code)
}