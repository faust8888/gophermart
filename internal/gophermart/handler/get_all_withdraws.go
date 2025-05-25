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
		logger.Log.Error("Invalid token")
		return
	}

	withdraws, err := r.withdrawHistoryService.FindAllHistoryWithdraws(claims.Login)
	if err != nil {
		logger.Log.Error("Error getting withdraws", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(&withdraws)
	if err != nil {
		logger.Log.Error("Error marshalling withdraws", zap.Error(err))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	writeSuccesfullResponse(resp, res)
}

func NewGetAllHistoryWithdrawsHandler(srv service.WithdrawHistoryService) GetAllWithdraws {
	return GetAllWithdraws{withdrawHistoryService: srv}
}
