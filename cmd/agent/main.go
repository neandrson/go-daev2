package main

import (
	"context"
	"fmt"
	"os"

	"github.com/neandrson/go-daev2/internal/agent/application"
	"github.com/neandrson/go-daev2/internal/agent/config"
)

func main() {
	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx := context.Background()
	app := application.NewApplication(cfg)
	exitCode := app.Run(ctx)
	os.Exit(exitCode)
}
