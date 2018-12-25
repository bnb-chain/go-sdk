package api

import (
	"github.com/binance-chain/go-sdk/sdk/tx/txmsg"
)

type ListPairResult struct {
	TxCommitResult
}

func (dex *dexAPI) ListPair(proposalId int64, baseAssetSymbol string, quoteAssetSymbol string, initPrice int64, sync bool) (*ListPairResult, error) {
	fromAddr := dex.keyManager.GetAddr()

	burnMsg := txmsg.NewDexListMsg(fromAddr, proposalId, baseAssetSymbol, quoteAssetSymbol, initPrice)
	err := burnMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := dex.broadcastMsg(burnMsg, sync)
	if err != nil {
		return nil, err
	}

	return &ListPairResult{*commit}, nil

}
