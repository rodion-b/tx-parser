package services

import (
	"context"
	"math/big"
	"tx-parser/internal/app/domain"
)

type EthDataSource interface {
	GetBlock(ctx context.Context, blockNumber *big.Int) (*domain.Block, error)
}

type Storage interface {
	Save(ctx context.Context, transactions map[string][]*domain.Transaction) error
	GetTransactions(ctx context.Context, addr string) ([]*domain.Transaction, error)
}
