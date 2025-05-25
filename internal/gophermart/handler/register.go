package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/config"
	"github.com/faust8888/gophermart/internal/gophermart/model"
	"github.com/faust8888/gophermart/internal/gophermart/repository/postgres"
	"github.com/faust8888/gophermart/internal/gophermart/security"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"github.com/faust8888/gophermart/internal/middleware/logger"
	"go.uber.org/zap"
	"net/http"
)

const UserRegisterHandlerPath = "/api/user/register"

type Register struct {
	registerService service.RegisterService
}

func (r *Register) RegisterUser(res http.ResponseWriter, req *http.Request) {
	requestBody, err := readRequestBody(req)
	if err != nil {
		logger.Log.Error("Error reading request body to register", zap.Error(err))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var registerUserRequest model.RegisterUserRequest
	if err = json.Unmarshal(requestBody.Bytes(), &registerUserRequest); err != nil {
		logger.Log.Error("Error unmarshalling request body to register", zap.Error(err))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if validationErrorMessage := validateRequest(registerUserRequest); validationErrorMessage != "" {
		logger.Log.Info("Failed to validate request to register", zap.String("validationError", validationErrorMessage))
		http.Error(res, validationErrorMessage, http.StatusBadRequest)
		return
	}

	err = r.registerService.Register(registerUserRequest.Login, registerUserRequest.Password)
	if err != nil {
		if errors.Is(err, postgres.ErrLoginUniqueViolation) {
			logger.Log.Info("Failed to register user. User already exists", zap.Error(err))
			http.Error(res, err.Error(), http.StatusConflict)
			return
		}
		logger.Log.Error("Failed to register user", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := security.BuildToken(config.AuthKey, registerUserRequest.Login)
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

func NewRegisterHandler(registerSrv service.RegisterService) Register {
	return Register{registerService: registerSrv}
}
