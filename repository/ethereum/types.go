package ethereum

/*
	Note: ideally we should create our own custom encoder/decoder from and to hex
*/

type (
	reqGeneric struct {
		JsonRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  []any  `json:"params"`
		ID      int64  `json:"id"`
	}
	respGeneric struct {
		JsonRPC string     `json:"jsonrpc"`
		Error   *respError `json:"error,omitempty"`
		ID      int64      `json:"id"`
	}
	respError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	respGenericString struct {
		respGeneric
		Result string `json:"result"`
	}
	respGetBlockByNumber struct {
		respGeneric
		Result block `json:"result"`
	}
)

type (
	block struct {
		blockData
		Transactions []transaction `json:"transactions"`
	}
	blockData struct {
		BaseFeePerGas    string `json:"baseFeePerGas"`
		Difficulty       string `json:"difficulty"`
		ExtraData        string `json:"extraData"`
		GasLimit         string `json:"gasLimit"`
		GasUsed          string `json:"gasUsed"`
		Hash             string `json:"hash"`
		LogsBloom        string `json:"logsBloom"`
		Miner            string `json:"miner"`
		MixHash          string `json:"mixHash"`
		Nonce            string `json:"nonce"`
		Number           string `json:"number"`
		ParentHash       string `json:"parentHash"`
		ReceiptsRoot     string `json:"receiptsRoot"`
		Sha3Uncles       string `json:"sha3Uncles"`
		Size             string `json:"size"`
		StateRoot        string `json:"stateRoot"`
		Timestamp        string `json:"timestamp"`
		TotalDifficulty  string `json:"totalDifficulty"`
		TransactionsRoot string `json:"transactionsRoot"`
		Uncles           []any  `json:"uncles"`
	}
	transaction struct {
		BlockHash            string       `json:"blockHash"`
		BlockNumber          string       `json:"blockNumber"`
		From                 string       `json:"from"`
		Gas                  string       `json:"gas"`
		GasPrice             string       `json:"gasPrice"`
		MaxFeePerGas         string       `json:"maxFeePerGas"`
		MaxPriorityFeePerGas string       `json:"maxPriorityFeePerGas"`
		Hash                 string       `json:"hash"`
		Input                string       `json:"input"`
		Nonce                string       `json:"nonce"`
		To                   string       `json:"to"`
		TransactionIndex     string       `json:"transactionIndex"`
		Value                string       `json:"value"`
		Type                 string       `json:"type"`
		AccessList           []accessList `json:"accessList"`
		ChainID              string       `json:"chainId"`
		V                    string       `json:"v"`
		R                    string       `json:"r"`
		S                    string       `json:"s"`
	}
	accessList struct {
		Address     string   `json:"address"`
		StorageKeys []string `json:"storageKeys"`
	}
)
