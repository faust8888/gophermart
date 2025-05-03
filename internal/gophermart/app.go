package gophermart

import (
	"github.com/faust8888/gophermart/internal/gophermart/config"
	"github.com/faust8888/gophermart/internal/gophermart/handler"
	"github.com/faust8888/gophermart/internal/gophermart/migration"
	"github.com/faust8888/gophermart/internal/gophermart/route"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type App struct {
	cfg   *config.Config
	Route *chi.Mux
}

func NewApp(handlerOptions ...handler.Option) *App {
	cfg, err := config.New()
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
	err := migration.Run(s.cfg.DatabaseURI)
	if err != nil {
		panic(err)
	}
	if err = http.ListenAndServe(s.cfg.RunAddress, s.Route); err != nil {
		panic(err)
	}
}
