package watcher

import (
	"context"
	"log"
	"sync"
	"time"

	modelEth "github.com/kucinghitam/ethereum-watcher/model/ethereum"
)

func (this *watcher) StartWatching() {
	quitChan := this.watchTransaction(
		this.watchBlock(
			this.watchBlockchain(),
			this.config.BufferTransaction,
			this.config.WorkerPoolBlock,
		),
		this.config.WorkerPoolTransaction,
		this.repoStorage.IsAddressSubscribed,
	)

	var quitHistoryChan <-chan bool
	{
		blockHistoryChan, isSubFunc := this.watchNewSubscribers()
		quitHistoryChan = this.watchTransaction(
			this.watchBlock(
				blockHistoryChan,
				this.config.BufferTransactionHistory,
				this.config.WorkerPoolBlockHistory,
			),
			this.config.WorkerPoolTransactionHistory,
			isSubFunc,
		)
	}

	// graceful stop
	<-quitChan
	<-quitHistoryChan
	this.quitDoneSignal <- true
}

func (this *watcher) StopWatching() {
	if this.quitSignal != nil {
		close(this.quitSignal)
		<-this.quitDoneSignal
	}
}

func (this *watcher) TriggerWatchNewSubscribers() {
	if len(this.newSubTriggerSignal) == 0 {
		this.newSubTriggerSignal <- true
	}
}

func (this *watcher) watchBlockchain() <-chan int {
	blockNumOut := make(chan int, this.config.BufferBlock)
	go func() {
		defer close(blockNumOut)
		timerCh := time.NewTicker(this.config.WatchBlockchainInterval)
		log.Printf("watching blockchain at interval %s\n", this.config.WatchBlockchainInterval.String())
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

func (this *watcher) watchBlock(blockIn <-chan int, buffer, workerPool int) <-chan modelEth.Transaction {
	txOut := make(chan modelEth.Transaction, buffer)
	go func() {
		defer close(txOut)

		buffer := make(chan struct{}, workerPool)
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
					} else if block == nil {
						log.Print("block is empty")
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

func (this *watcher) watchTransaction(
	txIn <-chan modelEth.Transaction,
	workerPool int,
	isSubFunc func(context.Context, string) (bool, error),
) <-chan bool {
	quitTxChan := make(chan bool, 1)
	go func() {
		defer close(quitTxChan)

		buffer := make(chan struct{}, workerPool)
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
						isSub, err := isSubFunc(ctx, addr)
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

const blockNumberEarliestTransaction = 46147

func (this *watcher) watchNewSubscribers() (
	<-chan int,
	func(context.Context, string) (bool, error),
) {
	blockNumOut := make(chan int, this.config.BufferBlockHistory)
	newSubs := make(map[string]struct{})
	go func() {
		defer close(blockNumOut)
		log.Printf("watching new subscribers at interval %s\n", this.config.WatchNewSubscribersInterval.String())

		timerCh := time.NewTicker(this.config.WatchNewSubscribersInterval)
		var addresses []string
		processFunc := func() {
			ctx := context.Background()
			latestBlock, err := this.repoStorage.GetLatestBlockNumber(ctx)
			if err != nil {
				log.Print("err", err)
				return
			}
			if latestBlock == 0 {
				latestBlock, err = this.repoEth.GetLatestBlockNumber(ctx)
				if err != nil {
					log.Print("err", err)
					return
				}
			}

			//delete previous new subs if exists
			if len(addresses) > 0 {
				err = this.repoStorage.DeleteNewSubscribers(ctx, addresses)
				if err != nil {
					log.Print("err", err)
					return
				}
			}

			addresses, err = this.repoStorage.FetchhAllNewSubscribers(ctx)
			if err != nil {
				log.Print("err", err)
				return
			}

			newSubs = make(map[string]struct{})
			for _, addr := range addresses {
				newSubs[addr] = struct{}{}
			}

			log.Printf("process new subscribers history parsing started, blocks: %d, new subscribers: %d\n", latestBlock-1, len(newSubs))

			for blockNum := blockNumberEarliestTransaction; blockNum < latestBlock; blockNum++ {
				blockNumOut <- blockNum
			}
		}
		var quit bool
		for !quit {
			select {
			case <-this.quitSignal:
				quit = true
				break

			case <-this.newSubTriggerSignal:
				processFunc()
			case <-timerCh.C:
				processFunc()
			}
		}
		log.Println("done quit watching new subscribers")
	}()
	return blockNumOut, func(ctx context.Context, address string) (bool, error) {
		_, ok := newSubs[address]
		return ok, nil
	}
}
