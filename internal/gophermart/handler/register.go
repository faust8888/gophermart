package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/model"
	"github.com/faust8888/gophermart/internal/gophermart/repository/postgres"
	"github.com/faust8888/gophermart/internal/gophermart/security"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"net/http"
)

const UserRegisterHandlerPath = "/api/user/register"

type Register struct {
	registerService service.RegisterService
}

func (r *Register) RegisterUser(res http.ResponseWriter, req *http.Request) {
	requestBody, err := readRequestBody(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var registerUserRequest model.RegisterUserRequest
	if err = json.Unmarshal(requestBody.Bytes(), &registerUserRequest); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if validationErrorMessage := validateRequest(registerUserRequest); validationErrorMessage != "" {
		http.Error(res, validationErrorMessage, http.StatusBadRequest)
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

	token, err := security.BuildToken(security.AuthSecretKey, registerUserRequest.Login)
	if err != nil {
		http.Error(res, fmt.Sprintf("build token: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	http.SetCookie(res, &http.Cookie{
		Name:  security.AuthorizationTokenName,
		Value: token,
	})

	res.WriteHeader(http.StatusOK)
}

func NewRegisterHandler(registerSrv service.RegisterService) Register {
	return Register{registerService: registerSrv}
}
