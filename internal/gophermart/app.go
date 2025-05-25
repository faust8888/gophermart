package gophermart

import (
	"context"
	"errors"
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/config"
	"github.com/faust8888/gophermart/internal/gophermart/handler"
	"github.com/faust8888/gophermart/internal/gophermart/migration"
	"github.com/faust8888/gophermart/internal/gophermart/route"
	"github.com/faust8888/gophermart/internal/middleware/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	cfg   *config.Config
	Route *chi.Mux
}

func NewApp(handlerOptions ...handler.Option) *App {
	err := logger.Initialize("INFO")
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v", err)
	}

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	err = migration.Run(cfg.DatabaseURI)
	if err != nil {
		panic(err)
	}

	h := handler.New(cfg, handlerOptions...)
	router := route.New(h)

	return &App{
		cfg:   cfg,
		Route: router,
	}
}

func (s *App) Run() {
	server := &http.Server{
		Addr:    s.cfg.RunAddress,
		Handler: s.Route,
	}

	// Handle OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Start server in background
	go func() {
		logger.Log.Info("Starting server", zap.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for signal to shut down
	<-sigChan
	logger.Log.Info("Gracefully shutting down server")

	// Create a timeout for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Fatal("Failed to shutdown server", zap.Error(err))
	}

	logger.Log.Info("Server exited gracefully")
}
