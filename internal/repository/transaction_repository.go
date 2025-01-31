package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Fox1N69/iq-testtask/internal/domain/entity"
	"github.com/Fox1N69/iq-testtask/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	Replenish(ctx context.Context, userID int64, amount int64) error
	Transfer(ctx context.Context, fromUserID, toUserID int64, amount int64) error
	LastTransactions(ctx context.Context, userID int64, limit int8) ([]entity.Transaction, error)
}

type transactionRepositoy struct {
	db  *pgxpool.Pool
	log logger.Logger
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepositoy{
		db:  db,
		log: logger.GetLogger(),
	}
}

func (r *transactionRepositoy) Replenish(ctx context.Context, userID int64, amount int64) error {
	r.log.Debugf("Replenish: starting transaction for userID %d with amount %d", userID, amount)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.log.Errorf("Replenish: failed to start transaction for userID %d: %v", userID, err)
		return err
	}
	defer tx.Rollback(ctx)

	const updateBalanceQuery = `UPDATE users SET balance = balance + $1 WHERE id = $2`

	_, err = tx.Exec(ctx, updateBalanceQuery, amount, userID)
	if err != nil {
		r.log.Errorf("Replenish: failed to update balance for userID %d: %v", userID, err)
		tx.Rollback(ctx)
		return err
	}

	const insertTransactionQuery = `
        INSERT INTO transactions (user_id, amount, type, description, timestamp)
        VALUES ($1, $2, 'replenish', 'Replenishment', $3)
    `

	_, err = tx.Exec(ctx, insertTransactionQuery, userID, amount, time.Now())
	if err != nil {
		r.log.Errorf("Replenish: failed to insert transaction for userID %d: %v", userID, err)
		tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Errorf("Replenish: failed to commit transaction for userID %d: %v", userID, err)
		return err
	}

	r.log.Infof("Replenish: successfully replenished userID %d with amount %d", userID, amount)
	return nil
}

func (r *transactionRepositoy) Transfer(
	ctx context.Context,
	fromUserID, toUserID int64,
	amount int64,
) error {
	r.log.Debugf("Transfer: starting transaction from userID %d to userID %d with amount %d", fromUserID, toUserID, amount)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.log.Errorf("Transfer: failed to start transaction from userID %d to userID %d: %v", fromUserID, toUserID, err)
		return err
	}
	defer tx.Rollback(ctx)

	const debitQuery = `UPDATE users SET balance = balance - $1 WHERE id = $2`

	_, err = tx.Exec(ctx, debitQuery, amount, fromUserID)
	if err != nil {
		r.log.Errorf("Transfer: failed to debit balance for userID %d: %v", fromUserID, err)
		return err
	}

	const creditQuery = `UPDATE users SET balance = balance + $1 WHERE id = $2`

	_, err = tx.Exec(ctx, creditQuery, amount, toUserID)
	if err != nil {
		r.log.Errorf("Transfer: failed to credit balance for userID %d: %v", toUserID, err)
		return err
	}

	const insertSenderTransaction = `
	INSERT INTO transactions (user_id, amount, type, description, timestamp)
	VALUES ($1, $2, 'transfer', $3, $4)
`
	descriptionSender := fmt.Sprintf("Transfer to user %d", toUserID)
	_, err = tx.Exec(ctx, insertSenderTransaction, fromUserID, amount, descriptionSender, time.Now())
	if err != nil {
		r.log.Errorf("Transfer: failed to record sender transaction for userID %d: %v", fromUserID, err)
		return fmt.Errorf("failed to record sender transaction: %w", err)
	}

	const insertRecipientTransaction = `
	INSERT INTO transactions (user_id, amount, type, description, timestamp)
	VALUES ($1, $2, 'transfer', $3, $4)
`
	descriptionRecipient := fmt.Sprintf("Transfer from user %d", fromUserID)
	_, err = tx.Exec(ctx, insertRecipientTransaction, toUserID, amount, descriptionRecipient, time.Now())
	if err != nil {
		r.log.Errorf("Transfer: failed to record recipient transaction for userID %d: %v", toUserID, err)
		return fmt.Errorf("failed to record recipient transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Errorf("Transfer: failed to commit transaction from userID %d to userID %d: %v", fromUserID, toUserID, err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.log.Infof("Transfer: successfully transferred %d from userID %d to userID %d", amount, fromUserID, toUserID)
	return nil
}

func (r *transactionRepositoy) LastTransactions(
	ctx context.Context,
	userID int64,
	limit int8,
) ([]entity.Transaction, error) {
	r.log.Debugf("LastTransactions: fetching last %d transactions for userID %d", limit, userID)

	const query = `
		SELECT id, user_id, amount, type, description, timestamp
		FROM transactions
		WHERE user_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, userID, limit)
	if err != nil {
		r.log.Errorf("LastTransactions: failed to fetch transactions for userID %d: %v", userID, err)
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	var transactions []entity.Transaction
	for rows.Next() {
		var t entity.Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Type, &t.Description, &t.Timestamp); err != nil {
			r.log.Errorf("LastTransactions: failed to scan transaction for userID %d: %v", userID, err)
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		r.log.Errorf("LastTransactions: failed to iterate over transactions for userID %d: %v", userID, err)
		return nil, fmt.Errorf("failed to iterate over transactions: %w", err)
	}

	r.log.Infof("LastTransactions: successfully fetched %d transactions for userID %d", len(transactions), userID)
	return transactions, nil
}
