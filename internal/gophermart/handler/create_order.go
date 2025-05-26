package handler

import (
	"context"
	"errors"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"github.com/faust8888/gophermart/internal/middleware/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

const CreateOrderHandlerPath = "/api/user/orders"

type CreateOrder struct {
	orderService service.OrderService
}

func (r *CreateOrder) CreateUserOrder(res http.ResponseWriter, req *http.Request) {
	isTokenCorrect, claims := validateToken(res, req)
	if !isTokenCorrect {
		logger.Log.Error("Invalid token")
		return
	}

	requestBody, err := readRequestBody(req)
	if err != nil {
		logger.Log.Error("Error reading request body to create order", zap.Error(err))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if !isOrderNumberValidByLuhn(requestBody.String()) {
		logger.Log.Error("Invalid order number", zap.String("orderNumber", requestBody.String()))
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	orderNumber, err := strconv.ParseInt(requestBody.String(), 10, 64)
	if err != nil {
		logger.Log.Error("Error parsing order number", zap.Error(err))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = r.orderService.CreateOrder(ctx, claims.Login, orderNumber)
	if err != nil {
		if errors.Is(err, service.ErrOrderWasCreatedBefore) {
			logger.Log.Error("Order was created before with the same user", zap.Error(err))
			res.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, service.ErrOrderWithAnotherUserExist) {
			logger.Log.Error("Order was created before by another user", zap.Error(err))
			http.Error(res, err.Error(), http.StatusConflict)
			return
		}
		logger.Log.Error("Error creating order", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusAccepted)
}

func NewCreateOrderHandler(srv service.OrderService) CreateOrder {
	return CreateOrder{orderService: srv}
}
