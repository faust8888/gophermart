package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/model"
	"time"
)

var ErrNoWithdrawExist = errors.New("withdraws doesn't exist")

type WithdrawRepository struct {
	db *sql.DB
}

func (w *WithdrawRepository) Withdraw(ctx context.Context, login string, order int64, sum float32) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	txn, err := w.db.Begin()
	if err != nil {
		return err
	}

	insertWithdrawHistoryQuery := `INSERT INTO withdraw_history (user_login, order_number, sum) VALUES ($1, $2, $3)`
	_, err = txn.ExecContext(ctx, insertWithdrawHistoryQuery, login, order, sum)
	if err != nil {
		return fmt.Errorf("postgres.repository.createWithdrawHistory: %w", err)
	}

	updateBalanceQuery := `UPDATE balance SET withdrawn_sum = withdrawn_sum + $1, current_sum = current_sum - $2 WHERE user_login = $3 AND current_sum - $4 >= 0;`
	res, err := txn.ExecContext(ctx, updateBalanceQuery, sum, sum, login, sum)
	if err != nil {
		if err = txn.Rollback(); err != nil {
			return fmt.Errorf("postgres.repository.updateWithdrawnSum: %w", err)
		}
		return fmt.Errorf("postgres.repository.updateWithdrawnSum: %w", err)
	}
	updatedRowsCount, err := res.RowsAffected()
	if err != nil {
		if err = txn.Rollback(); err != nil {
			return fmt.Errorf("postgres.repository.updateWithdrawnSum: %w", err)
		}
		return fmt.Errorf("postgres.repository.updateWithdrawnSum: %w", err)
	}
	if updatedRowsCount == 0 {
		if err = txn.Rollback(); err != nil {
			return fmt.Errorf("postgres.repository.updateWithdrawnSum: %w", err)
		}
		return ErrNotEnoughBalance
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (w *WithdrawRepository) FindAllHistoryWithdraws(ctx context.Context, login string) ([]model.WithdrawHistoryItemResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	rows, err := w.db.QueryContext(ctx, `SELECT order_number, sum, created_at FROM withdraw_history WHERE user_login = $1 ORDER BY created_at`, login)
	if err != nil {
		return nil, err
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres.repository.FindAllHistoryWithdraws: error during row iteration: %w", err)
	}

	defer rows.Close()

	var hasRows bool
	var historyWithdraws []model.WithdrawHistoryItemResponse
	for rows.Next() {
		var historyWithdraw model.WithdrawHistoryItemResponse
		if err = rows.Scan(&historyWithdraw.OrderNumber, &historyWithdraw.Sum, &historyWithdraw.ProcessedAt); err != nil {
			return nil, fmt.Errorf("postgres.repository.FindAllHistoryWithdraws: %w", err)
		}
		hasRows = true
		historyWithdraws = append(historyWithdraws, historyWithdraw)
	}
	if !hasRows {
		return nil, ErrNoWithdrawExist
	}
	return historyWithdraws, nil
}

func NewWithdrawRepository(db *sql.DB) *WithdrawRepository {
	return &WithdrawRepository{
		db: db,
	}
}
