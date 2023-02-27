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
		BlockHash        string
		BlockNumber      int
		TransactionIndex int
		Type             int
		From             string
		To               string
		Value            int
		Nonce            int
		Input            string
		Hash             string
	}
)
