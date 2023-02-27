package storage

import (
	"bytes"
	"context"
	"encoding/gob"
	"sync"

	modelEth "github.com/kucinghitam/ethereum-watcher/model/ethereum"
)

type (
	subscriber struct {
		mutex           sync.RWMutex
		transactionsMap sync.Map
	}
)

func (this *storage) IsAddressSubscribed(ctx context.Context, address string) (subscribed bool, err error) {
	_, subscribed = this.subscribers.Load(address)
	return
}

func (this *storage) SubscribeAddress(ctx context.Context, address string) (err error) {
	this.subscribers.LoadOrStore(address, &subscriber{})
	return
}

func (this *storage) AddTransactionToAddress(ctx context.Context, address string, transaction modelEth.Transaction) (err error) {
	var d *[]byte
	{
		var buff bytes.Buffer
		err = gob.NewEncoder(&buff).Encode(transaction)
		if err != nil {
			return
		}
		d = new([]byte)
		*d = buff.Bytes()
	}
	dAny, _ := this.transactions.LoadOrStore(transaction.Hash, d)
	d = dAny.(*[]byte)

	subAny, ok := this.subscribers.Load(address)
	if ok {
		sub := subAny.(*subscriber)
		sub.mutex.Lock()
		defer sub.mutex.Unlock()

		sub.transactionsMap.Store(transaction.Hash, d)
	}
	return
}
