package main

import (
	"context"
	"fmt"
	"os"

	"github.com/neandrson/go-daev2/internal/orchestrator/application"
	"github.com/neandrson/go-daev2/internal/orchestrator/config"
)

func main() {

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	app := application.NewApplication(cfg)
	exitCode := app.Run(ctx)
	os.Exit(exitCode)
}
