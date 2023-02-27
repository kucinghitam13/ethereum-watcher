package watcher

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	modelEth "github.com/kucinghitam/ethereum-watcher/model/ethereum"
)

func (this *watcher) StartWatching() {
	quitChan := this.watchTransaction(
		this.watchBlock(
			this.watchBlockchain(),
		),
	)

	// graceful stop
	<-quitChan
	this.quitDoneSignal <- true
}

func (this *watcher) StopWatching() {
	if this.quitSignal != nil {
		close(this.quitSignal)
		<-this.quitDoneSignal
	}
}

func (this *watcher) watchBlockchain() <-chan int {
	blockNumOut := make(chan int, this.config.BufferTransaction)
	go func() {
		defer close(blockNumOut)
		timerCh := time.NewTicker(this.config.WatchBlockchainInterval)
		fmt.Printf("watching blockchain at interval %s\n", this.config.WatchBlockchainInterval.String())
		var quit bool
		for !quit {
			select {
			case <-this.quitSignal:
				quit = true
				break
			case <-timerCh.C:
				ctx := context.Background()
				var wg sync.WaitGroup
				wg.Add(2)

				var latestBlock int
				var errLatestBlock error
				go func() {
					defer wg.Done()
					latestBlock, errLatestBlock = this.repoStorage.GetLatestBlockNumber(ctx)
					if errLatestBlock != nil {
						log.Print("err", errLatestBlock)
						return
					}
				}()

				var latestBlockEth int
				var errLatestBlockEth error
				go func() {
					defer wg.Done()
					latestBlockEth, errLatestBlockEth = this.repoEth.GetLatestBlockNumber(ctx)
					if errLatestBlockEth != nil {
						log.Print("err", errLatestBlockEth)
						return
					}
				}()

				wg.Wait()
				if errLatestBlock != nil || errLatestBlockEth != nil {
					continue
				}

				if latestBlock == 0 || latestBlock < latestBlockEth {
					// if latestBlock from storage is zero, we're only interested in adding newest block from eth
					// if not, we need to check all of the blocks between eth latest and storage latest too
					var blocks []int
					if latestBlock > 0 {
						for num := latestBlock + 1; num < latestBlockEth; num++ {
							blocks = append(blocks, num)
						}
					}
					blocks = append(blocks, latestBlockEth)

					for _, blockNum := range blocks {
						blockNumOut <- blockNum
					}

					err := this.repoStorage.SetLatestBlockNumber(ctx, latestBlockEth)
					if err != nil {
						log.Print("err", errLatestBlockEth)
					}
				}
			}
		}
		log.Println("done quit watching blockchain")
	}()
	return blockNumOut
}

func (this *watcher) watchBlock(blockIn <-chan int) <-chan modelEth.Transaction {
	txOut := make(chan modelEth.Transaction, this.config.BufferTransaction)
	go func() {
		defer close(txOut)

		buffer := make(chan struct{}, this.config.WorkerPoolBlock)
		for blockNumber := range blockIn {
			buffer <- struct{}{}
			go func(blockNumber int) {
				defer func() {
					<-buffer
				}()
				ctx := context.Background()

				isBlockLocked, blockLocker, err := this.repoStorage.LockBlockByNumber(ctx, blockNumber)
				if err != nil {
					log.Print("err", err)
					return
				}
				if isBlockLocked {
					block, err := this.repoEth.GetBlockByNumber(ctx, blockNumber)
					defer blockLocker.Unlock()
					if err != nil {
						log.Print("err", err)
						return
					}
					for _, tx := range block.Transactions {
						// send transaction to be checked
						txOut <- tx
					}
					log.Printf("block %d (%s) processed, found %d transactions", block.Number, block.Hash, len(block.Transactions))
				}

			}(blockNumber)
		}
		log.Println("done quit watching blocks")
	}()
	return txOut
}

func (this *watcher) watchTransaction(txIn <-chan modelEth.Transaction) <-chan bool {
	quitTxChan := make(chan bool, 1)
	go func() {
		defer close(quitTxChan)

		buffer := make(chan struct{}, this.config.WorkerPoolTransaction)
		for tx := range txIn {
			buffer <- struct{}{}
			go func(tx modelEth.Transaction) {
				defer func() {
					<-buffer
				}()
				ctx := context.Background()

				// check if transaction from or to address is registered in as our subscribers the call to storage
				// we can optimize this to batch the subscribing checking to reduce number of call
				var wg sync.WaitGroup
				wg.Add(2)
				for _, addr := range []string{tx.From, tx.To} {
					go func(addr string) {
						defer wg.Done()
						isSub, err := this.repoStorage.IsAddressSubscribed(ctx, addr)
						if err != nil {
							log.Print("err", err)
							return
						}
						if isSub {
							err = this.repoStorage.AddTransactionToAddress(ctx, addr, tx)
							if err != nil {
								log.Print("err", err)
								return
							}
							this.notifTransactionToAddr(addr, tx)
						}
					}(addr)
				}

				wg.Wait()

			}(tx)
		}
		log.Println("done quit watching transactions")
	}()
	return quitTxChan
}

func (this *watcher) notifTransactionToAddr(addr string, tx modelEth.Transaction) {
	log.Printf("new transaction for subscribed address of %s (from %s to %s value %d)\n", addr, tx.From, tx.To, tx.Value)
}
