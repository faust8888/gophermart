package service

type LoginService interface {
	Login(login, password string) error
}

type loginService struct {
	repo UserRepository
}

func (l *loginService) Login(login, password string) error {
	return l.repo.CheckUser(login, password)
}

func NewLoginService(repo UserRepository) LoginService {
	return &loginService{repo: repo}
}
