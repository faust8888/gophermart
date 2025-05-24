package handler

import (
	"encoding/json"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"github.com/faust8888/gophermart/internal/middleware/logger"
	"go.uber.org/zap"
	"net/http"
)

type GetAllWithdraws struct {
	withdrawHistoryService service.WithdrawHistoryService
}

func (r *GetAllWithdraws) GetAllHistoryWithdraws(res http.ResponseWriter, req *http.Request) {
	isTokenCorrect, claims := validateToken(res, req)
	if !isTokenCorrect {
		return
	}

	withdraws, err := r.withdrawHistoryService.FindAllHistoryWithdraws(claims.Login)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Log.Info("GetAllHistoryWithdraws", zap.String("login", claims.Login), zap.Int("size", len(withdraws)))
	resp, err := json.Marshal(&withdraws)
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

func NewGetAllHistoryWithdrawsHandler(srv service.WithdrawHistoryService) GetAllWithdraws {
	return GetAllWithdraws{withdrawHistoryService: srv}
}
