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
			WatchBlockchainIntervalSeconds   int `json:"watch_blockchain_interval_seconds"`
			WatchNewSubscribersIntervalHours int `json:"watch_new_subscribers_interval_hours"`

			BufferBlock              int `json:"buffer_block"`
			BufferBlockHistory       int `json:"buffer_block_history"`
			BufferTransaction        int `json:"buffer_transaction"`
			BufferTransactionHistory int `json:"buffer_transaction_history"`

			WorkerPoolBlock              int `json:"worker_pool_block"`
			WorkerPoolBlockHistory       int `json:"worker_pool_block_history"`
			WorkerPoolTransaction        int `json:"worker_pool_transaction"`
			WorkerPoolTransactionHistory int `json:"worker_pool_transaction_history"`
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
