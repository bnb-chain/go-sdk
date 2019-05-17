package transaction

import (
	"fmt"

	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type MintTokenResult struct {
	tx.TxCommitResult
}

func (c *client) MintToken(symbol string, amount int64, sync bool) (*MintTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Freeze token symbol can'c be empty ")
	}
	fromAddr := c.keyManager.GetAddr()

	mintMsg := msg.NewMintMsg(
		fromAddr,
		symbol,
		amount,
	)
	err := mintMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(mintMsg, "", sync)
	if err != nil {
		return nil, err
	}

	return &MintTokenResult{*commit}, nil

}
