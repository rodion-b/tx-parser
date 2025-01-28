package repository

import (
	"context"
	"tx-parser/internal/app/domain"
)

type Storage interface {
	Save(ctx context.Context, transactions map[string][]*domain.Transaction) error
	GetTransactions(ctx context.Context, addr string) ([]*domain.Transaction, error)
}
