package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tx-parser/internal/app/clients/ethrpc"
	"tx-parser/internal/app/config"
	"tx-parser/internal/app/httpserver"
	"tx-parser/internal/app/repository/inmem"
	"tx-parser/internal/app/services"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := run(); err != nil {
		log.Err(err).Msg("Error in run")
	}
	os.Exit(0)
}

func run() error {
	ctx := context.Background()

	// read config from env
	cfg, err := config.Read()
	if err != nil {
		return fmt.Errorf("error reading config: %v", err)
	}
	//creating new ETH RPC client
	ethRpcClient := ethrpc.NewClient()

	//Init inmem storage
	repo := inmem.NewStorage()

	//Init ParserService
	txParserService := services.NewTxParserService(ctx, ethRpcClient, repo)

	//Starting obesrivng new blocks
	go txParserService.Start()

	//Init Http server
	httpServer := httpserver.NewHttpServer(txParserService)

	//Setting up routes
	router := mux.NewRouter()
	router.HandleFunc("/currentblock", httpServer.GetCurrentBlock).Methods("GET")
	router.HandleFunc("/subscribe", httpServer.Subscribe).Methods("GET")
	router.HandleFunc("/transactions", httpServer.GetTransactions).Methods("GET")

	srv := &http.Server{
		Addr:    cfg.HTTP_ADDR,
		Handler: router,
	}

	// listen to OS signals and gracefully shutdown HTTP server
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Err(err).Msg("HTTP Server Shutdown")
		}
		close(stopped)
	}()

	// start HTTP server
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped
	log.Info().Msg("Server stopped")

	return nil
}
