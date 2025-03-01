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
		fmt.Fprintln(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	app := application.NewApplication(cfg)
	exitCode := app.Run(context.Background())
	os.Exit(exitCode)
}
