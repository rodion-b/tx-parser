package domain

import (
	"fmt"
	"math/big"
)

type Transaction struct {
	Hash  string
	From  string
	To    string
	Value *big.Float
}

func NewTransaction(hash, from, to string, value *big.Float) (*Transaction, error) {
	if hash == "" {
		return nil, fmt.Errorf("%w: hash is required", ErrRequired)
	}

	if from == "" {
		return nil, fmt.Errorf("%w: from is required", ErrRequired)
	}

	return &Transaction{
		Hash:  hash,
		From:  from,
		To:    to,
		Value: value,
	}, nil
}

func (t *Transaction) GetHash() string {
	return t.Hash
}

func (t *Transaction) GetFrom() string {
	return t.From
}

func (t *Transaction) GetTo() string {
	return t.To
}

func (t *Transaction) GetValue() *big.Float {
	return t.Value
}
