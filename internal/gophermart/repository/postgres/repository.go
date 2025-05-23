package postgres

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var ErrLoginUniqueViolation = errors.New("user with given login already exists")
var ErrUserNotExist = errors.New("user with given login and password doesn't exist")
var ErrOrderNumberAlreadyExist = errors.New("order number with already exists")
var ErrOrdersNotExist = errors.New("orders don't exist")
var ErrNotEnoughBalance = errors.New("not enough balance")

const (
	DatabaseDriverName = "pgx"
)

func errorIs(err error, code string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == code
	}
	return false
}
