package handler

import (
	"context"
	"encoding/json"
	"github.com/faust8888/gophermart/internal/gophermart/model"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"github.com/faust8888/gophermart/internal/middleware/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type GetBalance struct {
	balanceService service.BalanceService
}

func (r *GetBalance) GetUserBalance(res http.ResponseWriter, req *http.Request) {
	isTokenCorrect, claims := validateToken(res, req)
	if !isTokenCorrect {
		logger.Log.Error("Invalid token")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	balance, err := r.balanceService.FindCurrentBalance(ctx, claims.Login)
	if err != nil {
		logger.Log.Error("Failed to get current balance", zap.Error(err), zap.String("login", claims.Login))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(model.GetBalanceResponse{Current: balance.CurrentSum, Withdrawn: balance.WithdrawnSum})
	if err != nil {
		logger.Log.Error("Failed to marshal response to get the balance", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	writeSuccesfullResponse(resp, res)
}

func NewGetBalanceHandler(balanceService service.BalanceService) GetBalance {
	return GetBalance{balanceService: balanceService}
}
