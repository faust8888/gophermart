package handler

import (
	"errors"
	"github.com/faust8888/gophermart/internal/gophermart/security"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"net/http"
	"strconv"
)

type CreateOrder struct {
	orderService *service.OrderService
}

func (r *CreateOrder) CreateUserOrder(res http.ResponseWriter, req *http.Request) {
	token := security.GetToken(req)
	if token == "" {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims, err := security.GetClaims(token, security.AuthSecretKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if !security.CheckUserSession(claims.SessionId) {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	requestBody, err := readBody(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	orderNumber, err := strconv.ParseInt(requestBody.String(), 10, 64)
	if err != nil {
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

func NewCreateOrderHandler(srv *service.OrderService) CreateOrder {
	return CreateOrder{orderService: srv}
}
