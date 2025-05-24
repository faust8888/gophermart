package gophermart

import (
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/config"
	"github.com/faust8888/gophermart/internal/gophermart/handler"
	"github.com/faust8888/gophermart/internal/gophermart/migration"
	"github.com/faust8888/gophermart/internal/gophermart/route"
	"github.com/faust8888/gophermart/internal/middleware/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
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
	if err := http.ListenAndServe(s.cfg.RunAddress, s.Route); err != nil {
		panic(err)
	}
}
