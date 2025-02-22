package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Vojan-Najov/daec/internal/http/handler"
	"github.com/Vojan-Najov/daec/internal/orchestrator/config"
	"github.com/Vojan-Najov/daec/internal/service"
)

func Run(
	ctx context.Context,
	logger *log.Logger,
	cfg config.Config,
) (func(context.Context) error, error) {
	calcService := service.NewCalcService(cfg)

	muxHandler, err := newMuxHandler(ctx, logger, calcService)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{Addr: ":8081", Handler: muxHandler}
	logger.Printf("START SERVER ON PORT 8081\n")

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Printf("ListenAndServe: %v\n", err)
		}
	}()

	return srv.Shutdown, nil

}

func newMuxHandler(
	ctx context.Context,
	logger *log.Logger,
	calcService *service.CalcService,
) (http.Handler, error) {
	muxHandler, err := handler.NewHandler(ctx, calcService)
	if err != nil {
		return nil, fmt.Errorf("handler initialization error: %w", err)
	}

	// middleware для обработчиков
	muxHandler = handler.Decorate(muxHandler, loggingMiddleware(logger))

	return muxHandler, nil
}

// middleware для логированя запросов
func loggingMiddleware(
	logger *log.Logger,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Пропуск запроса к следующему обработчику
			next.ServeHTTP(w, r)

			// Завершение логирования после выполнения запроса
			if r.URL.Path == "/internal/task" && r.Method == "GET" {
				return
			}
			duration := time.Since(start)
			logger.Printf(
				"HTTP request - method: %s, path: %s, duration: %d\n",
				r.Method,
				r.URL.Path,
				duration,
			)
		})
	}
}
