package main

import (
	"context"

	_ "github.com/Goboolean/common/pkg/env"
)

func main() {
	ctx := context.Background()
	<- ctx.Done()
}