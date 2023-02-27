package config

import (
	"encoding/json"
	"os"
)

type (
	Config struct {
		Http struct {
			Port int `json:"port"`
		} `json:"http"`
		Ethereum struct {
			Host string `json:"host"`
		} `json:"ethereum"`

		Watcher struct {
			WatchBlockchainIntervalSeconds int `json:"watch_blockchain_interval_seconds"`

			BufferBlock       int `json:"buffer_block"`
			BufferTransaction int `json:"buffer_transaction"`

			WorkerPoolBlock       int `json:"worker_pool_block"`
			WorkerPoolTransaction int `json:"worker_pool_transaction"`
		} `json:"watcher"`
	}
)

var config *Config

func init() {
	f, err := os.Open("files/config.json")
	if err != nil {
		panic(err)
	}
	config = new(Config)
	err = json.NewDecoder(f).Decode(config)
	if err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	return config
}
