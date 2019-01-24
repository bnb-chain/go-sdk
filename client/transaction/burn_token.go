package transaction

import (
	"fmt"

	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type BurnTokenResult struct {
	tx.TxCommitResult
}

func (c *client) BurnToken(symbol string, amount int64, sync bool) (*BurnTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Burn token symbol can'c be empty ")
	}
	fromAddr := c.keyManager.GetAddr()

	burnMsg := msg.NewTokenBurnMsg(
		fromAddr,
		symbol,
		amount,
	)
	err := burnMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(burnMsg, sync)
	if err != nil {
		return nil, err
	}

	return &BurnTokenResult{*commit}, nil

}
