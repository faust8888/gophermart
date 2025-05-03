package service

import (
	"errors"
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/repository/postgres"
)

var ErrOrderWithAnotherUserExist = errors.New("order number with another user already exists")
var ErrOrderWasCreatedBefore = errors.New("order with this number was created before")

type OrderService struct {
	repo OrderRepository
}
type RegisterService struct {
	repo UserRepository
}

type LoginService struct {
	repo UserRepository
}

func (s *OrderService) CreateOrder(userLogin string, orderNumber int64) error {
	err := s.repo.CreateOrder(userLogin, orderNumber)
	if err != nil {
		if errors.Is(err, postgres.ErrOrderNumberAlreadyExist) {
			currentUserLogin, findLoginErr := s.repo.FindLoginByOrderNumber(orderNumber)
			if findLoginErr != nil {
				return fmt.Errorf("postgres.repository.FindLoginByOrderNumber: failed to scan row: %w", findLoginErr)
			}
			if currentUserLogin != userLogin {
				return ErrOrderWithAnotherUserExist
			}
			return ErrOrderWasCreatedBefore
		}
		return err
	}
	return nil
}

func (s *RegisterService) Register(login, password string) error {
	return s.repo.CreateUser(login, password)
}

func (l *LoginService) Login(login, password string) error {
	return l.repo.CheckUser(login, password)
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func NewRegisterService(repo UserRepository) *RegisterService {
	return &RegisterService{repo: repo}
}

func NewLoginService(repo UserRepository) *LoginService {
	return &LoginService{repo: repo}
}

type OrderRepository interface {
	CreateOrder(userLogin string, orderNumber int64) error
	FindLoginByOrderNumber(orderNumber int64) (string, error)
}

type UserRepository interface {
	CreateUser(login, password string) error
	CheckUser(login, password string) error
}
