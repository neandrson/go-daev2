package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Vojan-Najov/daec/internal/agent/application"
	"github.com/Vojan-Najov/daec/internal/agent/config"
)

func main() {
	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	app := application.NewApplication(cfg)
	exitCode := app.Run(context.Background())
	os.Exit(exitCode)
}
