package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/faust8888/gophermart/internal/gophermart/model"
	"github.com/faust8888/gophermart/internal/gophermart/repository/postgres"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"github.com/faust8888/gophermart/internal/middleware/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

const WithdrawHandlerPath = "/api/user/balance/withdraw"

type Withdraw struct {
	withdrawService service.WithdrawService
}

func (r *Withdraw) WithdrawSum(res http.ResponseWriter, req *http.Request) {
	isTokenCorrect, claims := validateToken(res, req)
	if !isTokenCorrect {
		logger.Log.Error("Invalid token")
		return
	}

	requestBody, err := readRequestBody(req)
	if err != nil {
		logger.Log.Error("Failed to read request body to withdraw", zap.Error(err))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var withdrawRequest model.WithdrawRequest
	if err = json.Unmarshal(requestBody.Bytes(), &withdrawRequest); err != nil {
		logger.Log.Error("Failed to unmarshal request body to withdraw", zap.Error(err))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if validationErrorMessage := validateRequest(withdrawRequest); validationErrorMessage != "" {
		logger.Log.Info("Failed to validate request to withdraw", zap.String("validationError", validationErrorMessage))
		http.Error(res, validationErrorMessage, http.StatusBadRequest)
		return
	}

	orderNumber, _ := strconv.ParseInt(withdrawRequest.Order, 10, 64)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = r.withdrawService.Withdraw(ctx, claims.Login, orderNumber, withdrawRequest.Sum)
	if err != nil {
		if errors.Is(err, postgres.ErrNotEnoughBalance) {
			logger.Log.Error("Not enough balance to withdraw", zap.Error(err))
			http.Error(res, err.Error(), http.StatusPaymentRequired)
			return
		}
		logger.Log.Error("Failed to withdraw", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func NewWithdrawHandler(withdrawSrv service.WithdrawService) Withdraw {
	return Withdraw{withdrawService: withdrawSrv}
}
