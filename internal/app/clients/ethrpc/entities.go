package ethrpc

import "encoding/json"

type JsonRPCRequest struct {
	ID      string        `json:"id"`
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type JsonRPCResponse struct {
	ID      string          `json:"id"`
	JsonRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *JsonRPCError   `json:"error,omitempty"`
}

type JsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BlockInformation struct {
	Number       string                   `json:"number"`
	Hash         string                   `json:"hash"`
	Transactions []TransactionInformation `json:"transactions"`
}

type TransactionInformation struct {
	Hash  string `json:"hash"`
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}