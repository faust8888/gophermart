package handler

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/config"
	"github.com/faust8888/gophermart/internal/gophermart/repository/postgres"
	"github.com/faust8888/gophermart/internal/gophermart/security"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

type Handler struct {
	Register
	Login
	CreateOrder
	FindOrders
	GetBalance
	Withdraw
	GetAllWithdraws
}

type RegisterHandler interface {
	RegisterUser(res http.ResponseWriter, req *http.Request)
}

type Option func(*Handler)

func New(cfg *config.Config, options ...Option) *Handler {
	db, err := sql.Open(postgres.DatabaseDriverName, cfg.DatabaseURI)
	if err != nil {
		panic(err)
	}

	registerService := service.NewRegisterService(postgres.NewUserRepository(db))
	loginService := service.NewLoginService(postgres.NewUserRepository(db))
	balanceService := service.NewBalanceService(postgres.NewBalanceRepository(db))
	orderService := service.NewOrderService(postgres.NewOrderRepository(db), balanceService, cfg)
	withdrawService := service.NewWithdrawService(postgres.NewWithdrawRepository(db), postgres.NewOrderRepository(db))
	withdrawHistoryService := service.NewWithdrawHistoryService(postgres.NewWithdrawRepository(db))

	h := &Handler{
		Register:        NewRegisterHandler(registerService),
		Login:           NewLoginHandler(loginService),
		CreateOrder:     NewCreateOrderHandler(orderService),
		FindOrders:      NewFindOrdersHandler(orderService),
		Withdraw:        NewWithdrawHandler(withdrawService),
		GetAllWithdraws: NewGetAllHistoryWithdrawsHandler(withdrawHistoryService),
		GetBalance:      NewGetBalanceHandler(balanceService),
	}

	for _, option := range options {
		option(h)
	}

	return h
}

func readRequestBody(req *http.Request) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

func writeSuccesfullResponse(response []byte, responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	_, err := responseWriter.Write(response)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}
}

func validateRequest(req interface{}) string {
	err := validator.New().Struct(req)
	if err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)
		errorMessage := "Validation failed:"
		for _, e := range validationErrors {
			errorMessage += fmt.Sprintf(" Field: %s, Error: %s;", e.Field(), e.Tag())
		}
		return errorMessage
	}
	return ""
}

func validateToken(res http.ResponseWriter, req *http.Request) (bool, *security.Claims) {
	token := security.GetToken(req)
	if token == "" {
		res.WriteHeader(http.StatusUnauthorized)
		return false, nil
	}
	claims, err := security.GetClaims(token, config.AuthKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return false, nil
	}

	if !security.CheckUserSession(claims.SessionID) {
		res.WriteHeader(http.StatusUnauthorized)
		return false, nil
	}
	return true, claims
}

func isOrderNumberValidByLuhn(number string) bool {
	var sum int
	double := false
	for i := len(number) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false // если есть нецифровые символы — невалидно
		}

		if double {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}
		sum += digit
		double = !double
	}
	return sum%10 == 0
}
