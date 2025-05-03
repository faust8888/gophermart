package route

import (
	"github.com/faust8888/gophermart/internal/middleware/logger"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type route interface {
	RegisterUser(res http.ResponseWriter, req *http.Request)
	LoginUser(res http.ResponseWriter, req *http.Request)
	CreateUserOrder(res http.ResponseWriter, req *http.Request)
	//FindOrder(res http.ResponseWriter, req *http.Request)
	//GetBalance(res http.ResponseWriter, req *http.Request)
	//Withdraw(res http.ResponseWriter, req *http.Request)
	//GetAllWithdraws(res http.ResponseWriter, req *http.Request)
}

func New(h route) *chi.Mux {
	router := chi.NewRouter()
	router.Use(logger.NewMiddleware)
	//router.Use(security.NewMiddleware)
	router.Post("/api/user/register", h.RegisterUser)
	router.Post("/api/user/login", h.LoginUser)
	router.Post("/api/user/orders", h.CreateUserOrder)
	//router.Get("/api/user/orders", h.FindOrder)
	//router.Get("/api/user/balance", h.GetBalance)
	//router.Post("/api/user/balance/withdraw", h.Withdraw)
	//router.Get("/api/user/withdrawals", h.GetAllWithdraws)
	return router
}
