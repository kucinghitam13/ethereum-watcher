package repository

import (
	"context"

	modelEth "github.com/kucinghitam/ethereum-watcher/model/ethereum"
)

type (
	Storage interface {
		BlockStorage
		SubscriberStorage
	}
	BlockStorage interface {
		// get latest processed block
		GetLatestBlockNumber(ctx context.Context) (blockNumber int, err error)
		// store latest processed block
		SetLatestBlockNumber(ctx context.Context, blockNumber int) (err error)
		// lock block by number
		// if not locked ok will return true and the locker object will be non-nil
		// if ok is false locker object return nil
		// for distributed locking purpose
		LockBlockByNumber(ctx context.Context, number int) (ok bool, locker BlockLocker, err error)
	}
	BlockLocker interface {
		Unlock()
	}

	SubscriberStorage interface {
		// check whether given address currently subscribes
		IsAddressSubscribed(ctx context.Context, address string) (subscribed bool, err error)
		// subscribe given address
		SubscribeAddress(ctx context.Context, address string) (err error)

		AddTransactionToAddress(ctx context.Context, address string, transaction modelEth.Transaction) (err error)
	}
)
