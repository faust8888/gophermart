package service

import "github.com/faust8888/gophermart/internal/gophermart/model"

type BalanceService interface {
	CreateDefaultBalance(userLogin string) error
	UpdateBalance(userLogin string, sum float32) error
	FindCurrentBalance(userLogin string) (*model.BalanceEntity, error)
}

type balanceService struct {
	balanceRepo BalanceRepository
}

func (b *balanceService) CreateDefaultBalance(userLogin string) error {
	return b.balanceRepo.CreateDefaultBalance(userLogin)
}

func (b *balanceService) UpdateBalance(userLogin string, sum float32) error {
	return b.balanceRepo.UpdateBalance(userLogin, sum)
}

func (b *balanceService) FindCurrentBalance(userLogin string) (*model.BalanceEntity, error) {
	return b.balanceRepo.FindCurrentBalance(userLogin)
}

type BalanceRepository interface {
	FindCurrentBalance(userLogin string) (*model.BalanceEntity, error)
	UpdateBalance(userLogin string, sum float32) error
	CreateDefaultBalance(userLogin string) error
}

func NewBalanceService(repo BalanceRepository) BalanceService {
	return &balanceService{balanceRepo: repo}
}
