package storage

import (
	"bytes"
	"context"
	"encoding/gob"
	"sort"
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

func (this *storage) SubscribeAddress(ctx context.Context, address string) (isAlreadySubscribed bool, err error) {
	_, isAlreadySubscribed = this.subscribers.LoadOrStore(address, &subscriber{})

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

func (this *storage) GetTransactionsByAddress(ctx context.Context, address string) (transactions []modelEth.Transaction, err error) {
	if subAny, ok := this.subscribers.Load(address); ok {
		sub := subAny.(*subscriber)
		sub.transactionsMap.Range(func(k, v any) bool {
			d := v.(*[]byte)
			var transaction modelEth.Transaction
			err = gob.NewDecoder(bytes.NewReader(*d)).Decode(&transaction)
			if err == nil {
				transactions = append(transactions, transaction)
			}
			return true
		})
	}
	// we sort the transactions on demand based on nonce
	// TODO: improves by creating an index by nonce and update it on AddTransactionToAddress
	// also can cause memory and cpu spike if transactions are huge
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Nonce < transactions[j].Nonce
	})
	return
}
