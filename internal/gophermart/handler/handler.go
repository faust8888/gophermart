package handler

import (
	"bytes"
	"github.com/faust8888/gophermart/internal/gophermart/config"
	"github.com/faust8888/gophermart/internal/gophermart/repository/postgres"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"net/http"
)

type Handler struct {
	Register
	Login
	CreateOrder
	//FindOrder
	//GetBalance
	//Withdraw
	//GetAllWithdraws
}

type Option func(*Handler)

func New(cfg *config.Config, options ...Option) *Handler {
	registerService := service.NewRegisterService(postgres.NewUserRepository(cfg))
	loginService := service.NewLoginService(postgres.NewUserRepository(cfg))
	orderService := service.NewOrderService(postgres.NewOrderRepository(cfg))

	h := &Handler{
		Register:    NewRegisterHandler(registerService),
		Login:       NewLoginHandler(loginService),
		CreateOrder: NewCreateOrderHandler(orderService),
	}

	for _, option := range options {
		option(h)
	}

	return h
}

func readBody(req *http.Request) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
