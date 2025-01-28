package services

import (
	"context"
	"math/big"
	"testing"
	"tx-parser/internal/app/domain"
	"tx-parser/internal/app/repository/inmem"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEthDataSource mocks the EthDataSource interface
type MockEthDataSource struct {
	mock.Mock
}

func (m *MockEthDataSource) GetBlock(ctx context.Context, blockNumber *big.Int) (*domain.Block, error) {
	args := m.Called(ctx, blockNumber)
	return args.Get(0).(*domain.Block), args.Error(1)
}

func TestTxParserService_Subscribe(t *testing.T) {
	ctx := context.Background()
	mockEthDataSource := new(MockEthDataSource)
	storage := inmem.NewStorage() // creating new inmem storage for tests

	service := NewTxParserService(ctx, mockEthDataSource, storage)

	// Test subscription of an address
	address := "0xabc123"
	subscribed := service.Subscribe(address)
	assert.True(t, subscribed, "Address should be subscribed successfully")

	// Test duplicate subscription
	subscribed = service.Subscribe(address)
	assert.False(t, subscribed, "Duplicate subscription should return false")
}

func TestTxParserService_GetTransactions(t *testing.T) {
	ctx := context.Background()
	mockEthDataSource := new(MockEthDataSource)
	storage := inmem.NewStorage() // creating new inmem storage for tests

	service := NewTxParserService(ctx, mockEthDataSource, storage)

	// Set up the observed addresses
	service.Subscribe("0xabc123")
	service.Subscribe("0xdef456")

	_, err := storage.GetTransactions(ctx, "0xabc123")
	assert.ErrorIs(t, err, domain.ErrNotFound) //assert no transactions found

	//Mocking transactions
	var transactions []*domain.Transaction

	tx1, err := domain.NewTransaction(
		"0x88a4738035f9e42dfaa7f20fe0bfa2eed808c119bcf395ea0c35310b12ac7190",
		"0xabc123",
		"0xdef456",
		big.NewFloat(100),
	)
	assert.NoError(t, err)

	tx2, err := domain.NewTransaction(
		"0x19c2cb177fede4fba2163925a6255bde6178d26b0b879c6795438daafcb810ef",
		"0xghi789",
		"0xjkl012",
		big.NewFloat(200),
	)
	assert.NoError(t, err)

	tx3, err := domain.NewTransaction(
		"0xb2ae7ae76fa2bcdd0acb0e497254fc62bca021ccff5f8944eb319ee8df8d2b36",
		"0xae68f8",
		"0xb8420s",
		big.NewFloat(300),
	)
	assert.NoError(t, err)

	transactions = append(transactions, tx1, tx2, tx3)

	//Mocking block
	block, err := domain.NewBlock(
		"0x6e24db6a2d2a4ff8e30ff47c90ad0df5b171b326504319305017c9922337bd59",
		1,
		transactions,
	)
	assert.NoError(t, err)

	// Test processBlockTransactions
	err = service.processBlockTransactions(block)
	assert.NoError(t, err)

	addressTransactions, err := service.GetTransactions("0xabc123")
	assert.NoError(t, err)

	assert.Len(t, addressTransactions, 1) //no transaction found
	assert.Equal(t, "0x88a4738035f9e42dfaa7f20fe0bfa2eed808c119bcf395ea0c35310b12ac7190", addressTransactions[0].Hash)

}
