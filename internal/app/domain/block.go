package domain

import "fmt"

type Block struct {
	hash         string
	number       uint64
	transactions []*Transaction
}

func NewBlock(hash string, number uint64, transactions []*Transaction) (*Block, error) {
	if hash == "" {
		return nil, fmt.Errorf("%w: hash is required", ErrRequired)
	}

	return &Block{
		hash:         hash,
		number:       number,
		transactions: transactions,
	}, nil
}

func (b *Block) GetHash() string {
	return b.hash
}

func (b *Block) GetNumber() uint64 {
	return b.number
}

func (b *Block) GetTransactions() []*Transaction {
	return b.transactions
}
