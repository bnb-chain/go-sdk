package transaction

import (
	"fmt"

	"github.com/bnb-chain/go-sdk/types/msg"
	"github.com/bnb-chain/go-sdk/types/tx"
)

type UnfreezeTokenResult struct {
	tx.TxCommitResult
}

func (c *client) UnfreezeToken(symbol string, amount int64, sync bool, options ...Option) (*UnfreezeTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Unfreeze token symbol can't be empty ")
	}
	fromAddr := c.keyManager.GetAddr()

	unfreezeMsg := msg.NewUnfreezeMsg(
		fromAddr,
		symbol,
		amount,
	)
	commit, err := c.broadcastMsg(unfreezeMsg, sync, options...)
	if err != nil {
		return nil, err
	}

	return &UnfreezeTokenResult{*commit}, nil

}
