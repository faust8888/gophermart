package handler

import (
	"encoding/json"
	"errors"
	"github.com/faust8888/gophermart/internal/gophermart/repository/postgres"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"github.com/faust8888/gophermart/internal/middleware/logger"
	"go.uber.org/zap"
	"net/http"
)

const GetOrderHandlerPath = "/api/user/orders"

type FindOrders struct {
	findOrderService service.OrderService
}

func (r *FindOrders) FindAllOrders(res http.ResponseWriter, req *http.Request) {
	isTokenCorrect, claims := validateToken(res, req)
	if !isTokenCorrect {
		return
	}

	orders, errNew := r.findOrderService.FindAllOrders(claims.Login)
	if errNew != nil {
		if errors.Is(errNew, postgres.ErrOrdersNotExist) {
			res.WriteHeader(http.StatusNoContent)
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(&orders)
	if err != nil {
		logger.Log.Error("Failed to marshal orders response", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(resp)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func NewFindOrdersHandler(srv service.OrderService) FindOrders {
	return FindOrders{findOrderService: srv}
}
