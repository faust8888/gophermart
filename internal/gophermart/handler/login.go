package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/config"
	"github.com/faust8888/gophermart/internal/gophermart/model"
	"github.com/faust8888/gophermart/internal/gophermart/security"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"github.com/faust8888/gophermart/internal/middleware/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const UserLoginHandlerPath = "/api/user/login"

type Login struct {
	loginService service.LoginService
}

func (l *Login) LoginUser(res http.ResponseWriter, req *http.Request) {
	requestBody, err := readRequestBody(req)
	if err != nil {
		logger.Log.Error("Error reading request body to login", zap.Error(err))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var loginUserRequest model.LoginUserRequest
	if err = json.Unmarshal(requestBody.Bytes(), &loginUserRequest); err != nil {
		logger.Log.Error("Error unmarshalling request body to login", zap.Error(err))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if validationErrorMessage := validateRequest(loginUserRequest); validationErrorMessage != "" {
		logger.Log.Info("Failed to validate request to login", zap.String("validationError", validationErrorMessage))
		http.Error(res, validationErrorMessage, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = l.loginService.Login(ctx, loginUserRequest.Login, loginUserRequest.Password)
	if err != nil {
		logger.Log.Error("Failed to login", zap.Error(err))
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := security.BuildToken(config.AuthKey, loginUserRequest.Login)
	if err != nil {
		logger.Log.Error("Failed to build token", zap.Error(err))
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
