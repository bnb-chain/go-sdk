package transaction

import (
	"fmt"

	"github.com/bnb-chain/go-sdk/types/msg"
	"github.com/bnb-chain/go-sdk/types/tx"
)

type BurnTokenResult struct {
	tx.TxCommitResult
}

func (c *client) BurnToken(symbol string, amount int64, sync bool, options ...Option) (*BurnTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Burn token symbol can't be empty ")
	}
	fromAddr := c.keyManager.GetAddr()

	burnMsg := msg.NewTokenBurnMsg(
		fromAddr,
		symbol,
		amount,
	)
	commit, err := c.broadcastMsg(burnMsg, sync, options...)
	if err != nil {
		return nil, err
	}

	return &BurnTokenResult{*commit}, nil

}
