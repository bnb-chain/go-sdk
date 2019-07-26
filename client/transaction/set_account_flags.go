package transaction

import (
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type SetAccountFlagsResult struct {
	tx.TxCommitResult
}

func (c *client) SetAccountFlags(flags uint64, sync bool, options ...Option) (*SetAccountFlagsResult, error) {
	fromAddr := c.keyManager.GetAddr()

	setAccMsg := msg.NewSetAccountFlagsMsg(
		fromAddr,
		flags,
	)
	commit, err := c.broadcastMsg(setAccMsg, sync, options...)
	if err != nil {
		return nil, err
	}

	return &SetAccountFlagsResult{*commit}, nil
}
