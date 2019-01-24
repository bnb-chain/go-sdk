package transaction

import (
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type ListPairResult struct {
	tx.TxCommitResult
}

func (c *client) ListPair(proposalId int64, baseAssetSymbol string, quoteAssetSymbol string, initPrice int64, sync bool) (*ListPairResult, error) {
	fromAddr := c.keyManager.GetAddr()

	burnMsg := msg.NewDexListMsg(fromAddr, proposalId, baseAssetSymbol, quoteAssetSymbol, initPrice)
	err := burnMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(burnMsg, sync)
	if err != nil {
		return nil, err
	}

	return &ListPairResult{*commit}, nil

}
