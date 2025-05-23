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

func (b *BalanceRepository) CreateDefaultBalance(userLogin string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `INSERT INTO balance (user_login, withdrawn_sum, current_sum) VALUES ($1, 0, 0)`

	_, err := b.db.ExecContext(ctx, query, userLogin)

	if err != nil {
		return fmt.Errorf("repository.postgres: couldn't create new order - %w", err)
	}

	return nil
}

func (b *BalanceRepository) UpdateBalance(userLogin string, sum float32) error {
	query := `
       UPDATE "balance"
           SET current_sum = $1 + balance.current_sum
       WHERE user_login = $2`
	_, err := b.db.Exec(query, sum, userLogin)
	if err != nil {
		return err
	}
	return nil
}

func (b *BalanceRepository) FindCurrentBalance(userLogin string) (*model.BalanceEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := `
        SELECT id, user_login, withdrawn_sum, current_sum, created_at
        FROM balance
        WHERE user_login = $1`

	rows, err := b.db.QueryContext(ctx, query, userLogin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	var balanceEntity model.BalanceEntity
	if err = rows.Scan(&balanceEntity.ID, &balanceEntity.UserLogin, &balanceEntity.WithdrawnSum, &balanceEntity.CurrentSum, &balanceEntity.CreatedAt); err != nil {
		return nil, err
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
