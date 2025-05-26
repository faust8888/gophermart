package service

import (
	"context"
	"github.com/faust8888/gophermart/internal/gophermart/model"
)

type WithdrawService interface {
	Withdraw(ctx context.Context, login string, order int64, sum float32) error
}

type WithdrawHistoryService interface {
	FindAllHistoryWithdraws(ctx context.Context, login string) ([]model.WithdrawHistoryItemResponse, error)
}

type withdrawService struct {
	withdrawRepo WithdrawRepository
	orderRepo    OrderRepository
}

type withdrawHistoryService struct {
	withdrawRepo WithdrawRepository
}

func (s *withdrawService) Withdraw(ctx context.Context, login string, orderNumber int64, sum float32) error {
	return s.withdrawRepo.Withdraw(ctx, login, orderNumber, sum)
}

func (s *withdrawHistoryService) FindAllHistoryWithdraws(ctx context.Context, login string) ([]model.WithdrawHistoryItemResponse, error) {
	return s.withdrawRepo.FindAllHistoryWithdraws(ctx, login)
}

func NewWithdrawService(withdrawRepo WithdrawRepository, orderRepo OrderRepository) WithdrawService {
	return &withdrawService{withdrawRepo: withdrawRepo, orderRepo: orderRepo}
}

func NewWithdrawHistoryService(withdrawRepo WithdrawRepository) WithdrawHistoryService {
	return &withdrawHistoryService{withdrawRepo: withdrawRepo}
}

type WithdrawRepository interface {
	Withdraw(ctx context.Context, login string, order int64, sum float32) error
	FindAllHistoryWithdraws(ctx context.Context, login string) ([]model.WithdrawHistoryItemResponse, error)
}
