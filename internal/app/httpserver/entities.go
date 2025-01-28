package httpserver

import "math/big"

type TransactionsResponse struct {
	Hash     string `json:"Hash"`
	From     string `json:"From"`
	To       string `json:"To"`
	Value    *big.Float `json:"Value"`
}
