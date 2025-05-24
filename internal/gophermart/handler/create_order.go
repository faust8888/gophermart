package handler

import (
	"errors"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"net/http"
	"strconv"
)

const CreateOrderHandlerPath = "/api/user/orders"

type CreateOrder struct {
	orderService service.OrderService
}

func (r *CreateOrder) CreateUserOrder(res http.ResponseWriter, req *http.Request) {
	isTokenCorrect, claims := validateToken(res, req)
	if !isTokenCorrect {
		return
	}

	requestBody, err := readRequestBody(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidByLuhn(requestBody.String()) {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	orderNumber, err := strconv.ParseInt(requestBody.String(), 10, 64)
	if err != nil {
		//4539319503436467
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.orderService.CreateOrder(claims.Login, orderNumber)
	if err != nil {
		if errors.Is(err, service.ErrOrderWasCreatedBefore) {
			res.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, service.ErrOrderWithAnotherUserExist) {
			http.Error(res, err.Error(), http.StatusConflict)
			return
		}
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusAccepted)
}

func NewCreateOrderHandler(srv service.OrderService) CreateOrder {
	return CreateOrder{orderService: srv}
}
