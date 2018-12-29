package api

import (
	"fmt"
	"github.com/binance-chain/go-sdk/sdk/tx/txmsg"
)

type CreateOrderResult struct {
	TxCommitResult
	OrderId string
}

func (dex *dexAPI) CreateOrder(baseAssetSymbol, quoteAssetSymbol string, op int8, price, quantity int64, sync bool) (*CreateOrderResult, error) {
	if baseAssetSymbol == "" || quoteAssetSymbol == "" {
		return nil, fmt.Errorf("BaseAssetSymbol or QuoteAssetSymbol is missing. ")
	}
	fromAddr := dex.keyManager.GetAddr()
	acc, err := dex.GetAccount(fromAddr.String())
	if err != nil {
		return nil, err
	}
	sequence := acc.Sequence

	orderId := txmsg.GenerateOrderID(sequence+1, dex.keyManager.GetAddr())
	newOrderMsg := txmsg.NewCreateOrderMsg(
		fromAddr,
		orderId,
		op,
		CombineSymbol(baseAssetSymbol, quoteAssetSymbol),
		price,
		quantity,
	)
	err = newOrderMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := dex.broadcastMsg(newOrderMsg, sync)
	if err != nil {
		return nil, err
	}
	return &CreateOrderResult{*commit, orderId}, nil
}
