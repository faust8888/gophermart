package handler

import (
	"encoding/json"
	"errors"
	"github.com/faust8888/gophermart/internal/gophermart/model"
	"github.com/faust8888/gophermart/internal/gophermart/repository/postgres"
	"github.com/faust8888/gophermart/internal/gophermart/service"
	"net/http"
	"strconv"
)

const WithdrawHandlerPath = "/api/user/balance/withdraw"

type Withdraw struct {
	withdrawService service.WithdrawService
}

func (r *Withdraw) WithdrawSum(res http.ResponseWriter, req *http.Request) {
	isTokenCorrect, claims := validateToken(res, req)
	if !isTokenCorrect {
		return
	}

	requestBody, err := readRequestBody(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var withdrawRequest model.WithdrawRequest
	if err = json.Unmarshal(requestBody.Bytes(), &withdrawRequest); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if validationErrorMessage := validateRequest(withdrawRequest); validationErrorMessage != "" {
		http.Error(res, validationErrorMessage, http.StatusBadRequest)
		return
	}

	orderNumber, _ := strconv.ParseInt(withdrawRequest.Order, 10, 64)

	err = r.withdrawService.Withdraw(claims.Login, orderNumber, withdrawRequest.Sum)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotExist) {
			http.Error(res, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, postgres.ErrNotEnoughBalance) {
			http.Error(res, err.Error(), http.StatusPaymentRequired)
			return
		}
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func NewWithdrawHandler(withdrawSrv service.WithdrawService) Withdraw {
	return Withdraw{withdrawService: withdrawSrv}
}
