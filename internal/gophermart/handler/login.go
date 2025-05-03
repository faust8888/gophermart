package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/model"
	"github.com/faust8888/gophermart/internal/gophermart/security"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"github.com/go-playground/validator/v10"
	"net/http"
)

const UserLoginHandlerPath = "/api/user/login"

type Login struct {
	loginService *service.LoginService
}

func (l *Login) LoginUser(res http.ResponseWriter, req *http.Request) {
	requestBody, err := readBody(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var loginUserRequest model.LoginUserRequest
	if err = json.Unmarshal(requestBody.Bytes(), &loginUserRequest); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = validateRequest(loginUserRequest); err != nil {
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)
		errorMessage := "Validation failed:"
		for _, e := range validationErrors {
			errorMessage += fmt.Sprintf(" Field: %s, Error: %s;", e.Field(), e.Tag())
		}
		http.Error(res, errorMessage, http.StatusBadRequest)
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

func NewLoginHandler(srv *service.LoginService) Login {
	return Login{loginService: srv}
}
