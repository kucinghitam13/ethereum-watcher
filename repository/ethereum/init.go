package ethereum

import (
	"net/http"

	"github.com/kucinghitam/ethereum-watcher/repository"
)

type (
	Config struct {
		Host string
	}

	eth struct {
		config *Config
		client *http.Client
	}
)

func New(config *Config) repository.Ethereum {
	return &eth{
		config: config,
		client: &http.Client{},
	}
}
