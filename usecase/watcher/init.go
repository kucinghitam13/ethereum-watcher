package watcher

import (
	"time"

	"github.com/kucinghitam/ethereum-watcher/repository"
	"github.com/kucinghitam/ethereum-watcher/usecase"
)

type (
	Config struct {
		WatchBlockchainInterval     time.Duration
		WatchNewSubscribersInterval time.Duration

		BufferBlock              int
		BufferBlockHistory       int
		BufferTransaction        int
		BufferTransactionHistory int

		WorkerPoolBlock              int
		WorkerPoolBlockHistory       int
		WorkerPoolTransaction        int
		WorkerPoolTransactionHistory int
	}
	watcher struct {
		config      *Config
		repoEth     repository.Ethereum
		repoStorage repository.Storage

		quitSignal          chan bool
		quitDoneSignal      chan bool
		newSubTriggerSignal chan bool
	}
)

func New(
	config *Config,
	repoEth repository.Ethereum,
	repoStorage repository.Storage,
) usecase.Watcher {
	return &watcher{
		config:              config,
		repoEth:             repoEth,
		repoStorage:         repoStorage,
		quitSignal:          make(chan bool, 1),
		quitDoneSignal:      make(chan bool, 1),
		newSubTriggerSignal: make(chan bool),
	}
}
