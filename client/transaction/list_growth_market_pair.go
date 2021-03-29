package transaction

import (
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type ListGrowthMarketPairResult struct {
	tx.TxCommitResult
}

func (c *client) ListGrowthMarketPair(baseAssetSymbol string, quoteAssetSymbol string, initPrice int64, sync bool, options ...Option) (*ListGrowthMarketPairResult, error) {
	fromAddr := c.keyManager.GetAddr()

	listMsg := msg.NewListGrowthMarketMsg(fromAddr, baseAssetSymbol, quoteAssetSymbol, initPrice)
	commit, err := c.broadcastMsg(listMsg, sync, options...)
	if err != nil {
		return nil, err
	}

	return &ListGrowthMarketPairResult{*commit}, nil

}
