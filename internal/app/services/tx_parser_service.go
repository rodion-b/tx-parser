package services

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"tx-parser/internal/app/domain"
	"tx-parser/internal/app/repository"

	"github.com/rs/zerolog/log"
)

type TxParserService struct {
	ctx              context.Context
	currentBlock     uint64
	ethDataSource    EthDataSource
	observeAddresses map[string]struct{}
	mu               sync.RWMutex
	repo             repository.Storage
}

func NewTxParserService(ctx context.Context, ethDataSource EthDataSource, storage Storage) *TxParserService {
	return &TxParserService{
		ctx:              ctx,
		ethDataSource:    ethDataSource,
		observeAddresses: make(map[string]struct{}),
		repo:             storage,
	}
}

func (s *TxParserService) GetCurrentBlock() uint64 {
	return atomic.LoadUint64(&s.currentBlock)
}

func (s *TxParserService) GetTransactions(address string) ([]*domain.Transaction, error) {
	address = strings.ToLower(address)
	return s.repo.GetTransactions(s.ctx, address)
}

func (s *TxParserService) Subscribe(address string) bool {
	address = strings.ToLower(address)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.observeAddresses[address]; ok {
		return false
	}
	s.observeAddresses[address] = struct{}{}
	return true
}

func (s *TxParserService) Start() {
	//processing latest block first
	if s.currentBlock == 0 {
		block, err := s.ethDataSource.GetBlock(s.ctx, nil)
		if err != nil {
			log.Err(err).Msg("error fetching current block")
			return
		}
		err = s.processBlockTransactions(block)
		if err != nil {
			log.Err(err).Msg("error processing block")
			return
		}
		atomic.StoreUint64(&s.currentBlock, uint64(block.GetNumber()))
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			block, err := s.ethDataSource.GetBlock(s.ctx, nil)
			if err != nil {
				log.Err(err).Msg("error fetching block")
				return
			}

			currentBlock := atomic.LoadUint64(&s.currentBlock)
			if block.GetNumber() == currentBlock {
				//skipping the rest as the retrieved block is already processed
				continue
			}

			//processing blocks including missed ones if any
			//Can be parallelize with waitgroups
			for i := currentBlock + 1; i <= block.GetNumber(); i++ {
				blockNumber := new(big.Int).SetUint64(i)
				block, err := s.ethDataSource.GetBlock(s.ctx, blockNumber)
				if err != nil {
					log.Err(err).Msg("error fetching block")
					return
				}
				err = s.processBlockTransactions(block)
				if err != nil {
					log.Err(err).Msg("error processing transaction")
					return
				}

			}
			atomic.StoreUint64(&s.currentBlock, block.GetNumber())
		}
	}
}

func (s *TxParserService) processBlockTransactions(block *domain.Block) error {
	transactions := block.GetTransactions()
	saveData := make(map[string][]*domain.Transaction, len(transactions))
	s.mu.Lock()
	for i := range transactions {
		from := strings.ToLower(transactions[i].From)
		to := strings.ToLower(transactions[i].To)
		_, validFrom := s.observeAddresses[from]
		_, validTo := s.observeAddresses[to]
		if !(validFrom || validTo) {
			continue
		}
		if validFrom {
			if _, ok := saveData[from]; !ok {
				saveData[from] = make([]*domain.Transaction, 0, len(transactions))
			}
			saveData[from] = append(saveData[from], transactions[i])
		}
		if validTo {
			if _, ok := saveData[to]; !ok {
				saveData[to] = make([]*domain.Transaction, 0, len(transactions))
			}
			saveData[to] = append(saveData[to], transactions[i])
		}
	}
	s.mu.Unlock()
	if len(saveData) == 0 {
		return nil
	}
	if err := s.repo.Save(s.ctx, saveData); err != nil {
		return fmt.Errorf("error saving block transactions: %w", err)
	}
	return nil
}
