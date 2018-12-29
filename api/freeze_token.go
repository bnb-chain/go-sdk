package api

import (
	"fmt"

	"github.com/binance-chain/go-sdk/tx/txmsg"
)

type FreezeTokenResult struct {
	TxCommitResult
}

func (dex *dexAPI) FreezeToken(symbol string, amount int64, sync bool) (*FreezeTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Freeze token symbol can't be empty ")
	}
	fromAddr := dex.keyManager.GetAddr()

	freezeMsg := txmsg.NewFreezeMsg(
		fromAddr,
		symbol,
		amount,
	)
	err := freezeMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := dex.broadcastMsg(freezeMsg, sync)
	if err != nil {
		return nil, err
	}

	return &FreezeTokenResult{*commit}, nil

}
