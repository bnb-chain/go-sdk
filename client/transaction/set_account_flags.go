package transaction

import (
	"fmt"

	"github.com/bnb-chain/go-sdk/common/types"
	"github.com/bnb-chain/go-sdk/types/msg"
	"github.com/bnb-chain/go-sdk/types/tx"
)

type SetAccountFlagsResult struct {
	tx.TxCommitResult
}

func (c *client) AddAccountFlags(flagOptions []types.FlagOption, sync bool, options ...Option) (*SetAccountFlagsResult, error) {
	fromAddr := c.keyManager.GetAddr()
	acc, err := c.queryClient.GetAccount(fromAddr.String())
	if err != nil {
		return nil, err
	}
	if len(flagOptions) == 0 {
		return nil, fmt.Errorf("missing flagOptions")
	}
	flags := acc.Flags
	for _, f := range flagOptions {
		flags = flags | uint64(f)
	}
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
