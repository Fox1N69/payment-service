package service

import (
	"context"

	"github.com/Fox1N69/iq-testtask/internal/domain/entity"
	"github.com/Fox1N69/iq-testtask/internal/repository"
)

type TransactionService interface {
	Replenish(ctx context.Context, userID int64, amount int64) error
	Transfer(ctx context.Context, fromUserID, toUserID int64, amount int64) error
	LastTransactions(ctx context.Context, userID int64, limit int8) ([]entity.Transaction, error)
}

type transactionService struct {
	repository repository.TransactionRepository
}

func NewTransactionService(repository repository.TransactionRepository) TransactionService {
	return &transactionService{
		repository: repository,
	}
}

func (s *transactionService) Replenish(ctx context.Context, userID int64, amount int64) error {
	return s.repository.Replenish(ctx, userID, amount)
}

func (s *transactionService) Transfer(ctx context.Context, fromUserID, toUserID int64, amount int64) error {
	return s.repository.Transfer(ctx, fromUserID, toUserID, amount)
}

func (s *transactionService) LastTransactions(ctx context.Context, userID int64, limit int8) ([]entity.Transaction, error) {
	return s.repository.LastTransactions(ctx, userID, limit)
}
