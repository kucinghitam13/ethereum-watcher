package ethereum

type (
	Block struct {
		Number          int
		ParentHash      string
		Hash            string
		TotalDifficulty int
		Nonce           int
		Transactions    []Transaction
	}
	Transaction struct {
		BlockHash        string `json:"block_hash"`
		BlockNumber      int    `json:"block_number"`
		TransactionIndex int    `json:"transaction_index"`
		Type             int    `json:"type"`
		From             string `json:"from"`
		To               string `json:"to"`
		Value            int    `json:"value"`
		Nonce            int    `json:"nonce"`
		Input            string `json:"input"`
		Hash             string `json:"hash"`
	}
)
