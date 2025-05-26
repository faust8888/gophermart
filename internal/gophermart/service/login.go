package service

import "context"

type LoginService interface {
	Login(ctx context.Context, login, password string) error
}

type loginService struct {
	repo UserRepository
}

func (l *loginService) Login(ctx context.Context, login, password string) error {
	return l.repo.CheckUser(ctx, login, password)
}

func NewLoginService(repo UserRepository) LoginService {
	return &loginService{repo: repo}
}
