package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kucinghitam/ethereum-watcher/config"
	handlerParser "github.com/kucinghitam/ethereum-watcher/delivery/parser/http"
	repoEth "github.com/kucinghitam/ethereum-watcher/repository/ethereum"
	repoSt "github.com/kucinghitam/ethereum-watcher/repository/storage"
	usecaseW "github.com/kucinghitam/ethereum-watcher/usecase/watcher"
)

func main() {
	config := config.GetConfig()
	repoEthereum := repoEth.New(&repoEth.Config{
		Host: config.Ethereum.Host,
	})
	repoStorage := repoSt.New(&repoSt.Config{})

	watcher := usecaseW.New(&usecaseW.Config{
		WatchBlockchainInterval:     time.Second * time.Duration(config.Watcher.WatchBlockchainIntervalSeconds),
		WatchNewSubscribersInterval: time.Hour * time.Duration(config.Watcher.WatchNewSubscribersIntervalHours),

		BufferBlock:              config.Watcher.BufferBlock,
		BufferTransaction:        config.Watcher.BufferTransaction,
		BufferBlockHistory:       config.Watcher.BufferBlockHistory,
		BufferTransactionHistory: config.Watcher.BufferTransactionHistory,

		WorkerPoolBlock:              config.Watcher.WorkerPoolBlock,
		WorkerPoolTransaction:        config.Watcher.WorkerPoolTransaction,
		WorkerPoolBlockHistory:       config.Watcher.WorkerPoolBlockHistory,
		WorkerPoolTransactionHistory: config.Watcher.WorkerPoolTransactionHistory,
	}, repoEthereum, repoStorage)
	go func() {
		watcher.StartWatching()
	}()

	handler := handlerParser.New(watcher)
	routeAPI(handler)

	srv := http.Server{
		Addr: fmt.Sprintf(":%d", config.Http.Port),
	}
	go func() {
		log.Printf("http server started run at port %d\n", config.Http.Port)
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Println("err http ", err)
		}
	}()

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("received signal %s, graceful shutdown started", <-sig)

	{
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		done := make(chan bool, 1)
		go func() {
			srv.Shutdown(ctx)
			close(done)
		}()
		select {
		case <-ctx.Done():
		case <-done:
		}
	}

	watcher.StopWatching()

	log.Println("graceful shutdown success, good bye")
}

func routeAPI(handler *handlerParser.Handler) {
	http.HandleFunc("/blocks/latest", handler.GetCurrentBlock)
	http.HandleFunc("/subscribes", handler.Subscribe)
	http.HandleFunc("/transactions", handler.GetTransactions)
}
