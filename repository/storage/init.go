package storage

import (
	"sync"

	"github.com/kucinghitam/ethereum-watcher/repository"
)

type (
	Config  struct{}
	storage struct {
		config *Config

		currentBlock int32
		blockLockers sync.Map

		subscribers  sync.Map
		transactions sync.Map
	}
)

func New(
	config *Config,
) repository.Storage {
	return &storage{
		config: config,
	}
}
