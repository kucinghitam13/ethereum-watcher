package usecase

import (
	"context"

	modelEth "github.com/kucinghitam/ethereum-watcher/model/ethereum"
)

type Watcher interface {
	StartWatching()
	StopWatching()

	TriggerWatchNewSubscribers()

	Parser
}

type Parser interface {
	// last parsed block
	GetCurrentBlock(ctx context.Context) int

	// add address to observer
	Subscribe(ctx context.Context, address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(ctx context.Context, address string) []modelEth.Transaction
}
