package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/neandrson/go-daev2/internal/agent/application"
	"github.com/neandrson/go-daev2/internal/agent/config"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

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
