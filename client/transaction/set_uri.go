package transaction

import (
	"github.com/bnb-chain/go-sdk/types/msg"
	"github.com/bnb-chain/go-sdk/types/tx"
)

type SetUriResult struct {
	tx.TxCommitResult
}

func (c *client) SetURI(symbol, tokenURI string, sync bool, options ...Option) (*SetUriResult, error) {
	fromAddr := c.keyManager.GetAddr()

	setURIMsg := msg.NewSetUriMsg(fromAddr, symbol, tokenURI)
	commit, err := c.broadcastMsg(setURIMsg, sync, options...)
	if err != nil {
		return nil, err
	}

	return &SetUriResult{*commit}, nil

}
