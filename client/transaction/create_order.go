package transaction

import (
	"encoding/json"
	"fmt"

	"github.com/bnb-chain/go-sdk/common"
	"github.com/bnb-chain/go-sdk/types/msg"
	"github.com/bnb-chain/go-sdk/types/tx"
)

type CreateOrderResult struct {
	tx.TxCommitResult
	OrderId string
}

// CreateOrder deprecated
func (c *client) CreateOrder(baseAssetSymbol, quoteAssetSymbol string, op int8, price, quantity int64, sync bool, options ...Option) (*CreateOrderResult, error) {
	if baseAssetSymbol == "" || quoteAssetSymbol == "" {
		return nil, fmt.Errorf("BaseAssetSymbol or QuoteAssetSymbol is missing. ")
	}
	fromAddr := c.keyManager.GetAddr()
	newOrderMsg := msg.NewCreateOrderMsg(
		fromAddr,
		"",
		op,
		common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol),
		price,
		quantity,
	)
	commit, err := c.broadcastMsg(newOrderMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	type commitData struct {
		OrderId string `json:"order_id"`
	}
	var cdata commitData
	if sync {
		err = json.Unmarshal([]byte(commit.Data), &cdata)
		if err != nil {
			return nil, err
		}
	}

	return &CreateOrderResult{*commit, cdata.OrderId}, nil
}
