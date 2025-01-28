package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"tx-parser/internal/app/domain"
	"tx-parser/internal/app/utils"
)

type HttpServer struct {
	service ParserService
}

func NewHttpServer(service ParserService) HttpServer {
	return HttpServer{
		service: service,
	}
}

func (h *HttpServer) GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	currentBlock := h.service.GetCurrentBlock()
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, currentBlock)
}

func (h *HttpServer) Subscribe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	address := r.URL.Query().Get("address")

	isValid := utils.IsValidEthereumAddress(address)
	if !isValid {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Subscription failed: Invalid Eth Address")
		return
	}

	success := h.service.Subscribe(address)
	if success {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Subscription successful")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Subscription failed: Internal Server Error")
	}
}

func (h *HttpServer) GetTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	address := r.URL.Query().Get("address")
	isValid := utils.IsValidEthereumAddress(address)
	if !isValid {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error: Invalid Eth Address")
		return
	}

	tx, err := h.service.GetTransactions(address)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error: No Transactions found for address")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error: Internal Server Error")
		}
		return
	}

	//if all good
	var transactionsResponse []TransactionsResponse
	for _, t := range tx {
		transactionResponse := TransactionsResponse{
			Hash:  t.Hash,
			From:  t.From,
			To:    t.To,
			Value: t.Value,
		}
		transactionsResponse = append(transactionsResponse, transactionResponse)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(transactionsResponse)
}
