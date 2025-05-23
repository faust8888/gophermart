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
	FindAllOrders(res http.ResponseWriter, req *http.Request)
	GetUserBalance(res http.ResponseWriter, req *http.Request)
	WithdrawSum(res http.ResponseWriter, req *http.Request)
	GetAllHistoryWithdraws(res http.ResponseWriter, req *http.Request)
}

func New(h route) *chi.Mux {
	router := chi.NewRouter()
	router.Use(logger.NewMiddleware)
	//router.Use(security.NewMiddleware)
	router.Post("/api/user/register", h.RegisterUser)
	router.Post("/api/user/login", h.LoginUser)
	router.Post("/api/user/orders", h.CreateUserOrder)
	router.Get("/api/user/orders", h.FindAllOrders)
	router.Get("/api/user/balance", h.GetUserBalance)
	router.Post("/api/user/balance/withdraw", h.WithdrawSum)
	router.Get("/api/user/withdrawals", h.GetAllHistoryWithdraws)
	return router
}
