package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/model"
	"github.com/rgurov/pgerrors"
	"time"
)

type OrderRepository struct {
	db *sql.DB
}

func (r *OrderRepository) CreateOrder(ctx context.Context, userLogin string, orderNumber int64) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	query := `INSERT INTO "order" (user_login, order_number, status) VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, userLogin, orderNumber, model.NewOrderStatus)

	if err != nil {
		if errorIs(err, pgerrors.UniqueViolation) {
			return ErrOrderNumberAlreadyExist
		}
		return fmt.Errorf("repository.postgres: couldn't create new order - %w", err)
	}

	return nil
}

func (r *OrderRepository) FindLoginByOrderNumber(ctx context.Context, orderNumber int64) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	query := `
        SELECT user_login
        FROM "order"
        WHERE order_number = $1
    `

	rows, err := r.db.QueryContext(ctx, query, orderNumber)
	if err != nil {
		return "", fmt.Errorf("postgres.repository.FindLoginByOrderNumber: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return "", nil
	}

	var userLogin string
	if err = rows.Scan(&userLogin); err != nil {
		return "", fmt.Errorf("postgres.repository.FindLoginByOrderNumber: failed to scan row: %w", err)
	}

	if err = rows.Err(); err != nil {
		return "", fmt.Errorf("postgres.repository.FindLoginByOrderNumber: error during row iteration: %w", err)
	}

	return userLogin, nil
}

func (r *OrderRepository) FindAllOrders(ctx context.Context, userLogin string) ([]model.OrderEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_login, order_number, status, accrual, created_at FROM "order" WHERE user_login = $1 ORDER BY created_at`, userLogin)
	if err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres.repository.FindAllOrders: error during row iteration: %w", err)
	}
	defer rows.Close()

	var hasRows bool
	var orders []model.OrderEntity
	for rows.Next() {
		var order model.OrderEntity
		if err = rows.Scan(&order.ID, &order.UserLogin, &order.OrderNumber, &order.Status, &order.Accrual, &order.CreatedAt); err != nil {
			return nil, fmt.Errorf("postgres.repository.FindAllOrders: %w", err)
		}
		hasRows = true
		orders = append(orders, order)
	}
	if !hasRows {
		return nil, ErrOrdersNotExist
	}
	return orders, nil
}

func (r *OrderRepository) FindAllOrdersForAccrualProcessing(ctx context.Context, limit int) ([]model.OrderEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_login, order_number, status, accrual, created_at FROM "order" WHERE status IN ($1, $2) 
                                 ORDER BY created_at LIMIT $3 FOR UPDATE SKIP LOCKED`,
		model.NewOrderStatus, model.ProcessingOrderStatus, limit)
	if err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres.repository.FindAllOrdersForAccrualProcessing: error during row iteration: %w", err)
	}
	defer rows.Close()

	var orders []model.OrderEntity
	for rows.Next() {
		var order model.OrderEntity
		if err = rows.Scan(&order.ID, &order.UserLogin, &order.OrderNumber, &order.Status, &order.Accrual, &order.CreatedAt); err != nil {
			return nil, fmt.Errorf("postgres.repository.FindAllOrdersForAccrualProcessing: %w", err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *OrderRepository) UpdateStatusAndAccrual(ctx context.Context, order model.OrderEntity) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	query := `
       UPDATE "order"
           SET status = $1, accrual = $2
       WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, order.Status, order.Accrual, order.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) BeginTransaction() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *OrderRepository) CommitTransaction(tx *sql.Tx) error {
	return tx.Commit()
}

func (r *OrderRepository) RollbackTransaction(tx *sql.Tx) error {
	return tx.Rollback()
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}
