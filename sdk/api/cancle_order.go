package api

import (
	"fmt"
	"github.com/BiJie/bnc-go-sdk/sdk/tx/txmsg"
)

type CancelOrderResult struct {
	TxCommitResult
}

func (dex *dexAPI) CancelOrder(baseAssetSymbol, quoteAssetSymbol, id, refId string, sync bool) (*CancelOrderResult, error) {
	if baseAssetSymbol == "" || quoteAssetSymbol == "" {
		return nil, fmt.Errorf("BaseAssetSymbol or QuoteAssetSymbol is missing. ")
	}
	if id == "" || refId == "" {
		return nil, fmt.Errorf("OrderId or Order RefId is missing. ")
	}

	fromAddr := dex.keyManager.GetAddr()

	cancelOrderMsg := txmsg.NewCancelOrderMsg(fromAddr, CombineSymbol(baseAssetSymbol, quoteAssetSymbol), id, refId)
	err := cancelOrderMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := dex.broadcastMsg(cancelOrderMsg, sync)
	if err != nil {
		return nil, err
	}

	return &CancelOrderResult{*commit}, nil

}
