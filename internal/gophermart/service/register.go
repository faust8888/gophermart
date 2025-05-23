package service

type RegisterService interface {
	Register(login, password string) error
}

type registerService struct {
	repo UserRepository
}

func (s *registerService) Register(login, password string) error {
	return s.repo.CreateUser(login, password)
}

func NewRegisterService(repo UserRepository) RegisterService {
	return &registerService{repo: repo}
}

type UserRepository interface {
	CreateUser(login, password string) error
	CheckUser(login, password string) error
}
