package service

import (
	"context"
	"github.com/faust8888/gophermart/internal/gophermart/model"
)

type BalanceService interface {
	CreateDefaultBalance(ctx context.Context, userLogin string) error
	UpdateBalance(ctx context.Context, userLogin string, sum float32) error
	FindCurrentBalance(ctx context.Context, userLogin string) (*model.BalanceEntity, error)
}

type balanceService struct {
	balanceRepo BalanceRepository
}

func (b *balanceService) CreateDefaultBalance(ctx context.Context, userLogin string) error {
	return b.balanceRepo.CreateDefaultBalance(ctx, userLogin)
}

func (b *balanceService) UpdateBalance(ctx context.Context, userLogin string, sum float32) error {
	return b.balanceRepo.UpdateBalance(ctx, userLogin, sum)
}

func (b *balanceService) FindCurrentBalance(ctx context.Context, userLogin string) (*model.BalanceEntity, error) {
	return b.balanceRepo.FindCurrentBalance(ctx, userLogin)
}

type BalanceRepository interface {
	FindCurrentBalance(ctx context.Context, userLogin string) (*model.BalanceEntity, error)
	UpdateBalance(ctx context.Context, userLogin string, sum float32) error
	CreateDefaultBalance(ctx context.Context, userLogin string) error
}

func NewBalanceService(repo BalanceRepository) BalanceService {
	return &balanceService{balanceRepo: repo}
}
