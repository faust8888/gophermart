package service

import "context"

type RegisterService interface {
	Register(ctx context.Context, login, password string) error
}

type registerService struct {
	repo UserRepository
}

func (s *registerService) Register(ctx context.Context, login, password string) error {
	return s.repo.CreateUser(ctx, login, password)
}

func NewRegisterService(repo UserRepository) RegisterService {
	return &registerService{repo: repo}
}

type UserRepository interface {
	CreateUser(ctx context.Context, login, password string) error
	CheckUser(ctx context.Context, login, password string) error
}
