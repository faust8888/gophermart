package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rgurov/pgerrors"
	"time"
)

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

	var userID int
	if err = rows.Scan(&userID); err != nil {
		return fmt.Errorf("postgres.repository.CheckUser: failed to scan row: %w", err)
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("postgres.repository.CheckUser: error during row iteration: %w", err)
	}

	return nil
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
