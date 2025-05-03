package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/model"
	"github.com/faust8888/gophermart/internal/gophermart/repository/postgres"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"github.com/go-playground/validator/v10"
	"net/http"
)

const UserRegisterHandlerPath = "/api/user/register"

type Register struct {
	registerService *service.RegisterService
}

func (r *Register) RegisterUser(res http.ResponseWriter, req *http.Request) {
	requestBody, err := readBody(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var registerUserRequest model.RegisterUserRequest
	if err = json.Unmarshal(requestBody.Bytes(), &registerUserRequest); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = validateRequest(registerUserRequest); err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)
		errorMessage := "Validation failed:"
		for _, e := range validationErrors {
			errorMessage += fmt.Sprintf(" Field: %s, Error: %s;", e.Field(), e.Tag())
		}
		http.Error(res, errorMessage, http.StatusBadRequest)
		return
	}

	err = r.registerService.Register(registerUserRequest.Login, registerUserRequest.Password)
	if err != nil {
		if errors.Is(err, postgres.ErrLoginUniqueViolation) {
			http.Error(res, err.Error(), http.StatusConflict)
			return
		}
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func NewRegisterHandler(srv *service.RegisterService) Register {
	return Register{registerService: srv}
}

func validateRequest(req interface{}) error {
	validate := validator.New()
	return validate.Struct(req)
}
