package watcher

import (
	"context"
	"log"

	modelEth "github.com/kucinghitam/ethereum-watcher/model/ethereum"
)

func (this *watcher) GetCurrentBlock(ctx context.Context) (blockNumber int) {
	blockNumber, err := this.repoStorage.GetLatestBlockNumber(ctx)
	if err != nil {
		log.Println("err ", err)
	}
	return
}

func (this *watcher) Subscribe(ctx context.Context, address string) bool {
	isAlreadySubscribed, err := this.repoStorage.SubscribeAddress(ctx, address)
	if err != nil {
		log.Println("err ", err)
	}
	if !isAlreadySubscribed {
		// add to new subscribers to be included when fetching all blocks
		err = this.repoStorage.AddNewlySubscriber(ctx, address)
		if err != nil {
			log.Println("err ", err)
		}
	}

	return !isAlreadySubscribed
}

func (this *watcher) GetTransactions(ctx context.Context, address string) []modelEth.Transaction {
	transactions, err := this.repoStorage.GetTransactionsByAddress(ctx, address)
	if err != nil {
		log.Println("err ", err)
	}
	return transactions
}
