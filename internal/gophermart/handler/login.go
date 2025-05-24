package handler

import (
	"encoding/json"
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/model"
	"github.com/faust8888/gophermart/internal/gophermart/security"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"net/http"
)

const UserLoginHandlerPath = "/api/user/login"

type Login struct {
	loginService service.LoginService
}

func (l *Login) LoginUser(res http.ResponseWriter, req *http.Request) {
	requestBody, err := readRequestBody(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var loginUserRequest model.LoginUserRequest
	if err = json.Unmarshal(requestBody.Bytes(), &loginUserRequest); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if validationErrorMessage := validateRequest(loginUserRequest); validationErrorMessage != "" {
		http.Error(res, validationErrorMessage, http.StatusBadRequest)
		return
	}

	err = l.loginService.Login(loginUserRequest.Login, loginUserRequest.Password)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := security.BuildToken(security.AuthSecretKey, loginUserRequest.Login)
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

func NewLoginHandler(srv service.LoginService) Login {
	return Login{loginService: srv}
}
