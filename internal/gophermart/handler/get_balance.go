package handler

import (
	"encoding/json"
	"github.com/faust8888/gophermart/internal/gophermart/model"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"github.com/faust8888/gophermart/internal/middleware/logger"
	"go.uber.org/zap"
	"net/http"
)

type GetBalance struct {
	balanceService service.BalanceService
}

func (r *GetBalance) GetUserBalance(res http.ResponseWriter, req *http.Request) {
	isTokenCorrect, claims := validateToken(res, req)
	if !isTokenCorrect {
		return
	}

	balance, err := r.balanceService.FindCurrentBalance(claims.Login)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Log.Info("GetUserBalance", zap.Any("login", balance.UserLogin), zap.Float32("sum", balance.CurrentSum), zap.Float32("withdrawn", balance.WithdrawnSum))

	resp, err := json.Marshal(model.GetBalanceResponse{Current: balance.CurrentSum, Withdrawn: balance.WithdrawnSum})
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(resp)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func NewGetBalanceHandler(balanceService service.BalanceService) GetBalance {
	return GetBalance{balanceService: balanceService}
}
