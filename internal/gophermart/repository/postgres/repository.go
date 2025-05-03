package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/config"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rgurov/pgerrors"
	"time"
)

var ErrLoginUniqueViolation = errors.New("user with given login already exists")
var ErrUserNotExist = errors.New("user with given login and password doesn't exist")
var ErrOrderNumberAlreadyExist = errors.New("order number with already exists")

const (
	DatabaseDriverName = "pgx"
)

type OrderRepository struct {
	db *sql.DB
}

func (r *OrderRepository) CreateOrder(userLogin string, orderNumber int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `INSERT INTO "order" (user_login, order_number) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, userLogin, orderNumber)

	if err != nil {
		if errorIs(err, pgerrors.UniqueViolation) {
			return ErrOrderNumberAlreadyExist
		}
		return fmt.Errorf("repository.postgres: couldn't create new order - %w", err)
	}

	return nil
}

func (r *OrderRepository) FindLoginByOrderNumber(orderNumber int64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
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

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) CreateUser(login, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO "user" (login, password) VALUES ($1, $2)`, login, password)

	if err != nil {
		if errorIs(err, pgerrors.UniqueViolation) {
			return ErrLoginUniqueViolation
		}
		return fmt.Errorf("repository.postgres: couldn't save new user - %w", err)
	}

	return nil
}

func (r *UserRepository) CheckUser(login, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := `
        SELECT id
        FROM "user"
        WHERE login = $1
        AND password = $2
    `
	rows, err := r.db.QueryContext(ctx, query, login, password)
	if err != nil {
		return fmt.Errorf("postgres.repository.CheckUser: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return ErrUserNotExist
	}

	var userId int
	if err = rows.Scan(&userId); err != nil {
		return fmt.Errorf("postgres.repository.CheckUser: failed to scan row: %w", err)
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("postgres.repository.CheckUser: error during row iteration: %w", err)
	}

	return nil
}

func NewOrderRepository(cfg *config.Config) *OrderRepository {
	db, err := sql.Open(DatabaseDriverName, cfg.DatabaseURI)
	if err != nil {
		panic(err)
	}
	return &OrderRepository{
		db: db,
	}
}

func NewUserRepository(cfg *config.Config) *UserRepository {
	db, err := sql.Open(DatabaseDriverName, cfg.DatabaseURI)
	if err != nil {
		panic(err)
	}
	return &UserRepository{
		db: db,
	}
}

func errorIs(err error, code string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == code
	}
	return false
}
