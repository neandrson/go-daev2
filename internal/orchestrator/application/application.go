package application

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/Vojan-Najov/daec/internal/http/server"
	"github.com/Vojan-Najov/daec/internal/orchestrator/config"
)

type Application struct {
	cfg config.Config
}

func NewApplication(cfg *config.Config) *Application {
	return &Application{
		cfg: *cfg,
	}
}

func (orch *Application) Run(ctx context.Context) int {
	logger := log.New(
		os.Stderr,
		"Orchestrator: ",
		log.Ldate|log.Ltime|log.Lmsgprefix,
	)

	shutDownFunc, err := server.Run(ctx, logger, orch.cfg)
	if err != nil {
		logger.Printf("Run server error: %v\n", err)
		return 1
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-c
	cancel()
	shutDownFunc(ctx)

	return 0
}
