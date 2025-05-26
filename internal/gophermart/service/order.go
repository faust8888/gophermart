package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/config"
	"github.com/faust8888/gophermart/internal/gophermart/model"
	"github.com/faust8888/gophermart/internal/gophermart/repository/postgres"
	"github.com/faust8888/gophermart/internal/middleware/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

var ErrOrderWithAnotherUserExist = errors.New("order number with another user already exists")
var ErrOrderWasCreatedBefore = errors.New("order with this number was created before")

type OrderService interface {
	CreateOrder(ctx context.Context, userLogin string, orderNumber int64) error
	FindAllOrders(ctx context.Context, userLogin string) ([]model.OrderItemResponse, error)
	FetchOrderStatusFromAccrual(asyncJob FetchOrderAccrualStatusJob)
}

type orderService struct {
	orderRepo      OrderRepository
	balanceService BalanceService
	workerPool     *WorkerPool
	cfg            *config.Config
}

func (s *orderService) CreateOrder(ctx context.Context, userLogin string, orderNumber int64) error {
	tx, err := s.orderRepo.BeginTransaction()
	if err != nil {
		return err
	}

	logger.Log.Info("Creating Order", zap.String("userLogin", userLogin), zap.Int64("orderNumber", orderNumber))
	err = s.orderRepo.CreateOrder(ctx, userLogin, orderNumber)
	if err != nil {
		if errors.Is(err, postgres.ErrOrderNumberAlreadyExist) {
			currentUserLogin, findLoginErr := s.orderRepo.FindLoginByOrderNumber(ctx, orderNumber)
			if findLoginErr != nil {
				return fmt.Errorf("postgres.OrderRepository.FindLoginByOrderNumber: failed to scan row: %w", findLoginErr)
			}
			if currentUserLogin != userLogin {
				return ErrOrderWithAnotherUserExist
			}
			return ErrOrderWasCreatedBefore
		}
		return err
	}

	err = s.balanceService.CreateDefaultBalance(ctx, userLogin)
	if err != nil {
		if err = s.orderRepo.RollbackTransaction(tx); err != nil {
			logger.Log.Error("Error rolling back transaction", zap.Error(err))
		}
		return err
	}

	if err = s.orderRepo.CommitTransaction(tx); err != nil {
		return err
	}
	return nil
}

func (s *orderService) FindAllOrders(ctx context.Context, userLogin string) ([]model.OrderItemResponse, error) {
	orders, err := s.orderRepo.FindAllOrders(ctx, userLogin)
	if err != nil {
		return nil, err
	}

	var response []model.OrderItemResponse
	for _, order := range orders {
		status, _ := toAccrualStatus(order.Status)
		var accrualStr float32
		if order.Accrual != nil {
			accrualStr = *order.Accrual
		} else {
			accrualStr = 0.00
		}
		response = append(response, model.OrderItemResponse{
			OrderNumber: strconv.FormatInt(order.OrderNumber, 10),
			Accrual:     accrualStr,
			Status:      status,
			UploadedAt:  order.CreatedAt.Format(time.RFC3339),
		})
	}

	return response, nil
}

func (s *orderService) FetchOrderStatusFromAccrual(asyncJob FetchOrderAccrualStatusJob) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	orders, err := s.orderRepo.FindAllOrdersForAccrualProcessing(ctx, asyncJob.selectLimit)
	if err != nil {
		logger.Log.Error("Couldn't select orders to fetch statuses", zap.Error(err))
		return
	}

	ordersToUpdate := make([]model.OrderEntity, 0, len(orders))
	for _, order := range orders {
		res, sendingErr := createPostRequest(order.OrderNumber, s.cfg.AccrualSystemAddress).Send()
		if sendingErr != nil {
			logger.Log.Error("Couldn't create post request to fetch status", zap.Error(sendingErr))
			return
		}

		if res.StatusCode() == http.StatusOK {
			var orderAccrualResponse model.OrderAccrualResponse
			if err = json.Unmarshal(res.Body(), &orderAccrualResponse); err != nil {
				logger.Log.Error("Error unmarshalling order accrual response", zap.Error(err))
				continue
			}
			order.Status, _ = toOrderStatus(orderAccrualResponse.Status)
			if orderAccrualResponse.Status == "PROCESSED" {
				a := orderAccrualResponse.Accrual
				order.Accrual = &a
			}
			ordersToUpdate = append(ordersToUpdate, order)
			logger.Log.Info("Job result of fetching order accrual status",
				zap.Int64("orderNumber", order.OrderNumber),
				zap.String("status", orderAccrualResponse.Status),
				zap.Int("task", asyncJob.taskNumber),
			)
		} else {
			logger.Log.Info("Job result of fetching order accrual status",
				zap.Int("response code", res.StatusCode()),
				zap.Int("task", asyncJob.taskNumber),
			)
		}
	}

	for _, orderToUpdate := range ordersToUpdate {
		tx, beginTransactionErr := s.orderRepo.BeginTransaction()
		if beginTransactionErr != nil {
			logger.Log.Error("Couldn't start transaction", zap.Error(err))
			continue
		}

		err = s.orderRepo.UpdateStatusAndAccrual(ctx, orderToUpdate)
		if err != nil {
			logger.Log.Error("Error updating orders", zap.Error(err))
		}

		err = s.balanceService.UpdateBalance(ctx, orderToUpdate.UserLogin, *orderToUpdate.Accrual)
		if err != nil {
			logger.Log.Error("Error updating orders", zap.Error(err))
			if err = s.orderRepo.RollbackTransaction(tx); err != nil {
				logger.Log.Error("Error rolling back transaction", zap.Error(err))
			}
			continue
		}
		err = s.orderRepo.CommitTransaction(tx)
		if err != nil {
			logger.Log.Error("Couldn't commit transaction", zap.Error(err))
		}
	}
}

func NewOrderService(orderRepo OrderRepository, balanceService BalanceService, cfg *config.Config) OrderService {
	wp := NewWorkerPool(1, 1, 5)
	srv := &orderService{orderRepo: orderRepo, balanceService: balanceService, workerPool: wp, cfg: cfg}
	go wp.Run(srv.FetchOrderStatusFromAccrual)
	return srv
}

func createPostRequest(orderNumber int64, host string) *resty.Request {
	req := resty.New().R()
	req.Method = http.MethodGet
	req.URL = host + "/api/orders/{orderNumber}"
	req.PathParams = map[string]string{
		"orderNumber": strconv.FormatInt(orderNumber, 10),
	}
	return req
}

func toOrderStatus(accrualStatus string) (string, error) {
	switch accrualStatus {
	case "REGISTERED":
		return model.NewOrderStatus, nil
	case "INVALID":
		return model.InvalidOrderStatus, nil
	case "PROCESSING":
		return model.ProcessingOrderStatus, nil
	case "PROCESSED":
		return model.ProcessedOrderStatus, nil
	default:
		return model.InvalidOrderStatus, nil
	}
}

func toAccrualStatus(orderStatus string) (string, error) {
	switch orderStatus {
	case model.NewOrderStatus:
		return model.NewOrderStatus, nil
	case model.InvalidOrderStatus:
		return "INVALID", nil
	case model.ProcessingOrderStatus:
		return "PROCESSING", nil
	case model.ProcessedOrderStatus:
		return "PROCESSED", nil
	default:
		return model.InvalidOrderStatus, nil
	}
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, userLogin string, orderNumber int64) error
	FindLoginByOrderNumber(ctx context.Context, orderNumber int64) (string, error)
	FindAllOrdersForAccrualProcessing(ctx context.Context, selectLimit int) ([]model.OrderEntity, error)
	FindAllOrders(ctx context.Context, userLogin string) ([]model.OrderEntity, error)
	UpdateStatusAndAccrual(ctx context.Context, order model.OrderEntity) error
	BeginTransaction() (*sql.Tx, error)
	CommitTransaction(tx *sql.Tx) error
	RollbackTransaction(tx *sql.Tx) error
}
