package transaction

import (
	"fmt"

	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type UnfreezeTokenResult struct {
	tx.TxCommitResult
}

func (c *client) UnfreezeToken(symbol string, amount int64, sync bool) (*UnfreezeTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Freeze token symbol can'c be empty ")
	}
	fromAddr := c.keyManager.GetAddr()

	unfreezeMsg := msg.NewUnfreezeMsg(
		fromAddr,
		symbol,
		amount,
	)
	err := unfreezeMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(unfreezeMsg, sync)
	if err != nil {
		return nil, err
	}

	return &UnfreezeTokenResult{*commit}, nil

}
