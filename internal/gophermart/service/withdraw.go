package service

import (
	"errors"
	"github.com/faust8888/gophermart/internal/gophermart/model"
)

var ErrOrderNotExist = errors.New("order doesn't exist")

type WithdrawService interface {
	Withdraw(login string, order int64, sum float32) error
}

type WithdrawHistoryService interface {
	FindAllHistoryWithdraws(login string) ([]model.WithdrawHistoryItemResponse, error)
}

type withdrawService struct {
	withdrawRepo WithdrawRepository
	orderRepo    OrderRepository
}

type withdrawHistoryService struct {
	withdrawRepo WithdrawRepository
}

func (s *withdrawService) Withdraw(login string, orderNumber int64, sum float32) error {
	return s.withdrawRepo.Withdraw(login, orderNumber, sum)
}

func (s *withdrawHistoryService) FindAllHistoryWithdraws(login string) ([]model.WithdrawHistoryItemResponse, error) {
	return s.withdrawRepo.FindAllHistoryWithdraws(login)
}

func NewWithdrawService(withdrawRepo WithdrawRepository, orderRepo OrderRepository) WithdrawService {
	return &withdrawService{withdrawRepo: withdrawRepo, orderRepo: orderRepo}
}

func NewWithdrawHistoryService(withdrawRepo WithdrawRepository) WithdrawHistoryService {
	return &withdrawHistoryService{withdrawRepo: withdrawRepo}
}

type WithdrawRepository interface {
	Withdraw(login string, order int64, sum float32) error
	FindAllHistoryWithdraws(login string) ([]model.WithdrawHistoryItemResponse, error)
}
