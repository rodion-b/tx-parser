package inmem

import (
	"context"
	"sync"
	"tx-parser/internal/app/domain"
)

type Storage struct {
	mu   sync.RWMutex
	data map[string][]*domain.Transaction
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string][]*domain.Transaction, 1_000),
	}
}

func (s *Storage) Save(_ context.Context, transactions map[string][]*domain.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for addr := range transactions {
		if _, ok := s.data[addr]; !ok {
			s.data[addr] = make([]*domain.Transaction, 0, len(transactions[addr]))
		}
		s.data[addr] = append(s.data[addr], transactions[addr]...)
	}
	return nil
}

func (s *Storage) GetTransactions(_ context.Context, addr string) ([]*domain.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.data[addr]; !ok {
		return nil, domain.ErrNotFound
	}
	return s.data[addr], nil
}
