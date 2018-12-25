package api

import (
	"fmt"
	"github.com/binance-chain/go-sdk/sdk/tx/txmsg"
)

type UnfreezeTokenResult struct {
	TxCommitResult
}

func (dex *dexAPI) UnfreezeToken(symbol string, amount int64, sync bool) (*UnfreezeTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Freeze token symbol can't be empty ")
	}
	fromAddr := dex.keyManager.GetAddr()

	unfreezeMsg := txmsg.NewUnfreezeMsg(
		fromAddr,
		symbol,
		amount,
	)
	err := unfreezeMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := dex.broadcastMsg(unfreezeMsg, sync)
	if err != nil {
		return nil, err
	}

	return &UnfreezeTokenResult{*commit}, nil

}
