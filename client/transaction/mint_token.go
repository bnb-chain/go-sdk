package transaction

import (
	"fmt"

	"github.com/bnb-chain/go-sdk/types/msg"
	"github.com/bnb-chain/go-sdk/types/tx"
)

type MintTokenResult struct {
	tx.TxCommitResult
}

func (c *client) MintToken(symbol string, amount int64, sync bool, options ...Option) (*MintTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Mint token symbol can't be empty ")
	}
	fromAddr := c.keyManager.GetAddr()

	mintMsg := msg.NewMintMsg(
		fromAddr,
		symbol,
		amount,
	)
	commit, err := c.broadcastMsg(mintMsg, sync, options...)
	if err != nil {
		return nil, err
	}

	return &MintTokenResult{*commit}, nil

}
