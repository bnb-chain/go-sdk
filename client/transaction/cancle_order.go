package transaction

import (
	"fmt"

	"github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type CancelOrderResult struct {
	tx.TxCommitResult
}

func (c *client) CancelOrder(baseAssetSymbol, quoteAssetSymbol, refId string, sync bool) (*CancelOrderResult, error) {
	if baseAssetSymbol == "" || quoteAssetSymbol == "" {
		return nil, fmt.Errorf("BaseAssetSymbol or QuoteAssetSymbol is missing. ")
	}
	if refId == "" {
		return nil, fmt.Errorf("OrderId or Order RefId is missing. ")
	}

	fromAddr := c.keyManager.GetAddr()

	cancelOrderMsg := msg.NewCancelOrderMsg(fromAddr, common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol), refId)
	err := cancelOrderMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(cancelOrderMsg, sync)
	if err != nil {
		return nil, err
	}

	return &CancelOrderResult{*commit}, nil

}
