package usecase

import (
	"context"

	modelEth "github.com/kucinghitam/ethereum-watcher/model/ethereum"
)

type Watcher interface {
	// start watching blockchain for new block and its transactions
	// for subscribed addresses
	StartWatching()
	// graceful stop
	StopWatching()

	// manually trigger checking for new subscribers and fetch all of its transaction
	// from blockchain instead of waiting for next scheduled cycle
	// long running, will iterate from first known transaction block to latest
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
