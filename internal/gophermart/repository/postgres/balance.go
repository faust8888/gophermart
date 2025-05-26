package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/model"
	"time"
)

type BalanceRepository struct {
	db *sql.DB
}

func (b *BalanceRepository) CreateDefaultBalance(ctx context.Context, userLogin string) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	query := `INSERT INTO balance (user_login, withdrawn_sum, current_sum) VALUES ($1, 0, 0)`
	_, err := b.db.ExecContext(ctx, query, userLogin)

	if err != nil {
		return fmt.Errorf("postgres.BalanceRepository.CreateDefaultBalance: %w", err)
	}

	return nil
}

func (b *BalanceRepository) UpdateBalance(ctx context.Context, userLogin string, sum float32) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	query := `UPDATE "balance" SET current_sum = $1 + balance.current_sum WHERE user_login = $2`
	_, err := b.db.ExecContext(ctx, query, sum, userLogin)

	if err != nil {
		return fmt.Errorf("postgres.BalanceRepository.UpdateBalance: %w", err)
	}

	return nil
}

func (b *BalanceRepository) FindCurrentBalance(ctx context.Context, userLogin string) (*model.BalanceEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	query := `SELECT id, user_login, withdrawn_sum, current_sum, created_at FROM balance WHERE user_login = $1`
	rows, err := b.db.QueryContext(ctx, query, userLogin)
	if err != nil {
		return nil, fmt.Errorf("postgres.BalanceRepository.FindCurrentBalance: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	var balanceEntity model.BalanceEntity
	if err = rows.Scan(&balanceEntity.ID, &balanceEntity.UserLogin, &balanceEntity.WithdrawnSum, &balanceEntity.CurrentSum, &balanceEntity.CreatedAt); err != nil {
		return nil, fmt.Errorf("postgres.BalanceRepository.FindCurrentBalance: couldn't scan rows - %w", err)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &balanceEntity, nil
}

func NewBalanceRepository(db *sql.DB) *BalanceRepository {
	return &BalanceRepository{
		db: db,
	}
}
