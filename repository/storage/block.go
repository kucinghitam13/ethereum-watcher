package storage

import (
	"context"
	"sync/atomic"

	"github.com/kucinghitam/ethereum-watcher/repository"
)

type (
	blockLocker struct {
		unlock func()
	}
)

func (this *storage) GetLatestBlockNumber(ctx context.Context) (blockNumber int, err error) {
	blockNumber = int(atomic.LoadInt32(&this.currentBlock))
	return
}

func (this *storage) SetLatestBlockNumber(ctx context.Context, blockNumber int) (err error) {
	atomic.StoreInt32(&this.currentBlock, int32(blockNumber))
	return
}

func (this *storage) LockBlockByNumber(ctx context.Context, number int) (ok bool, locker repository.BlockLocker, err error) {
	lAny, isLoaded := this.blockLockers.LoadOrStore(number, &blockLocker{
		unlock: func() {
			this.blockLockers.Delete(number)
		},
	})
	ok = !isLoaded
	if !ok {
		return
	}
	locker = lAny.(*blockLocker)

	return
}

func (this *blockLocker) Unlock() {
	this.unlock()
}
