package ethereum

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	model "github.com/kucinghitam/ethereum-watcher/model/ethereum"
)

const (
	valJsonRPC = "2.0"
	valID      = 13
)

func (this *eth) GetLatestBlockNumber(ctx context.Context) (blockNumber int, err error) {
	reqPayload, _ := json.Marshal(reqGeneric{
		JsonRPC: valJsonRPC,
		Method:  "eth_blockNumber",
		Params:  []any{},
		ID:      valID,
	})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, this.config.Host, bytes.NewReader(reqPayload))
	resp, err := this.client.Do(req)
	if err != nil {
		log.Println("err", err)
		return
	} else if resp.StatusCode > 299 {
		log.Println("err status", resp.StatusCode)
		return
	}

	var respPayload respGenericString
	err = json.NewDecoder(resp.Body).Decode(&respPayload)
	if err != nil {
		log.Println("err", err)
		return
	}

	blockNumber = int(hexToNumber(respPayload.Result))

	return
}

func (this *eth) GetBlockByNumber(ctx context.Context, number int) (block *model.Block, err error) {
	reqPayload, _ := json.Marshal(reqGeneric{
		JsonRPC: valJsonRPC,
		Method:  "eth_getBlockByNumber",
		Params: []any{
			fmt.Sprintf("0x%s", strconv.FormatInt(int64(number), 16)),
			true,
		},
		ID: valID,
	})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, this.config.Host, bytes.NewReader(reqPayload))
	resp, err := this.client.Do(req)
	if err != nil {
		log.Println("err", err)
		return
	} else if resp.StatusCode > 299 {
		log.Println("err status", resp.StatusCode)
		return
	}

	var respPayload respGetBlockByNumber
	err = json.NewDecoder(resp.Body).Decode(&respPayload)
	if err != nil {
		log.Println("err", err)
		return
	}
	block = &model.Block{
		Number:          int(hexToNumber(respPayload.Result.Number)),
		ParentHash:      respPayload.Result.ParentHash,
		Hash:            respPayload.Result.Hash,
		TotalDifficulty: int(hexToNumber(respPayload.Result.TotalDifficulty)),
		Nonce:           int(hexToNumber(respPayload.Result.Nonce)),
	}
	for _, v := range respPayload.Result.Transactions {
		block.Transactions = append(block.Transactions, model.Transaction{
			BlockHash:        v.BlockHash,
			BlockNumber:      int(hexToNumber(v.BlockNumber)),
			TransactionIndex: int(hexToNumber(v.TransactionIndex)),
			Type:             int(hexToNumber(v.Type)),
			From:             v.From,
			To:               v.To,
			Value:            int(hexToNumber(v.Value)),
			Nonce:            int(hexToNumber(v.Nonce)),
			Input:            v.Input,
			Hash:             v.Hash,
		})
	}

	return
}
