package repository

import (
	"context"

	model "github.com/kucinghitam/ethereum-watcher/model/ethereum"
)

type Ethereum interface {
	GetLatestBlockNumber(ctx context.Context) (blockNumber int, err error)
	GetBlockByNumber(ctx context.Context, number int) (block *model.Block, err error)
}
