package httpserver

import "tx-parser/internal/app/domain"

type ParserService interface {
	GetCurrentBlock() uint64
	Subscribe(address string) bool
	GetTransactions(address string) ([]*domain.Transaction, error)
}
