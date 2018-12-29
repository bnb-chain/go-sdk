package api

import (
	"fmt"

	"github.com/binance-chain/go-sdk/tx/txmsg"
)

type BurnTokenResult struct {
	TxCommitResult
}

func (dex *dexAPI) BurnToken(symbol string, amount int64, sync bool) (*BurnTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Burn token symbol can't be empty ")
	}
	fromAddr := dex.keyManager.GetAddr()

	burnMsg := txmsg.NewTokenBurnMsg(
		fromAddr,
		symbol,
		amount,
	)
	err := burnMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := dex.broadcastMsg(burnMsg, sync)
	if err != nil {
		return nil, err
	}

	return &BurnTokenResult{*commit}, nil

}
