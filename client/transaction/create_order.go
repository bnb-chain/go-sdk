package transaction

import (
	"fmt"

	"github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type CreateOrderResult struct {
	tx.TxCommitResult
	OrderId string
}

func (c *client) CreateOrder(baseAssetSymbol, quoteAssetSymbol string, op int8, price, quantity int64, sync bool, options ...Option) (*CreateOrderResult, error) {
	if baseAssetSymbol == "" || quoteAssetSymbol == "" {
		return nil, fmt.Errorf("BaseAssetSymbol or QuoteAssetSymbol is missing. ")
	}
	fromAddr := c.keyManager.GetAddr()
	acc, err := c.queryClient.GetAccount(fromAddr.String())
	if err != nil {
		return nil, err
	}
	sequence := acc.Sequence

	orderId := msg.GenerateOrderID(sequence+1, c.keyManager.GetAddr())
	newOrderMsg := msg.NewCreateOrderMsg(
		fromAddr,
		orderId,
		op,
		common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol),
		price,
		quantity,
	)
	err = newOrderMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(newOrderMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &CreateOrderResult{*commit, orderId}, nil
}
