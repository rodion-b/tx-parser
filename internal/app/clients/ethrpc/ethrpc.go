package ethrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tx-parser/internal/app/domain"

	"github.com/google/uuid"
)

const ethRpcUrl = "https://ethereum-rpc.publicnode.com"

type Client struct{}

var HTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetBlock(ctx context.Context, blockNumber *big.Int) (*domain.Block, error) {

	// Determine the block number to query
	var blockNumberHex string
	if blockNumber == nil {
		blockNumberHex = "latest" // Fetch the latest block if blockNumber is nil
	} else {
		blockNumberHex = fmt.Sprintf("0x%x", blockNumber) // Convert the block number to hex
	}
	//generate rpc requestId
	requestId := uuid.New()

	// Prepare the JSON-RPC request payload
	requestPayload := JsonRPCRequest{
		JsonRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{blockNumberHex, true}, // Full transaction objects
		ID:      requestId.String(),
	}

	// Serialize the request to JSON
	requestBody, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize JSON: %w", err)
	}

	// Make the HTTP POST request
	resp, err := http.Post(ethRpcUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Decode the response
	var response JsonRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	//verify the response error
	if response.Error != nil {
		return nil, fmt.Errorf("JSON-RPC error: %s", response.Error.Message)
	}

	//verify the response id match supplied id
	if response.ID != requestId.String() {
		return nil, fmt.Errorf("JSON-RPC error - wrong respinse id. Expected: %s, Actual: %s", requestId.String(), response.ID)
	}

	// Unmarshal the RawMessage into BlockInformation
	var blockInfo BlockInformation
	err = json.Unmarshal(response.Result, &blockInfo)
	if err != nil {
		return nil, fmt.Errorf("JSON-RPC error: %w", err)
	}

	//getting the blocknumber int
	blockNumberInt, err := strconv.ParseUint(strings.TrimPrefix(blockInfo.Number, "0x"), 16, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert block number to int: %w", err)
	}

	//parsing int domain transactions
	var transactions []*domain.Transaction
	for _, tx := range blockInfo.Transactions {
		value := new(big.Int)
		value.SetString(tx.Value[2:], 16) // Convert value from hex to integer
		ethValue := new(big.Float).Quo(new(big.Float).SetInt(value), big.NewFloat(1e18))
		transaction, err := domain.NewTransaction(tx.Hash, tx.From, tx.To, ethValue)
		if err != nil {
			return nil, fmt.Errorf("error creating transaction:: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	//creating new domain block
	block, err := domain.NewBlock(blockInfo.Hash, blockNumberInt, transactions)
	if err != nil {
		return nil, fmt.Errorf("failed to create new block %w", err)
	}
	return block, nil
}
