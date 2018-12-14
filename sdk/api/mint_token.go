package api

import (
	"fmt"
	"github.com/BiJie/bnc-go-sdk/sdk/tx/txmsg"
)

type MintTokenResult struct {
	TxCommitResult
}

func (dex *dexAPI) MintToken(symbol string, amount int64, sync bool) (*MintTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Freeze token symbol can't be empty ")
	}
	fromAddr := dex.keyManager.GetAddr()

	mintMsg := txmsg.NewMintMsg(
		fromAddr,
		symbol,
		amount,
	)
	err := mintMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := dex.broadcastMsg(mintMsg, sync)
	if err != nil {
		return nil, err
	}

	return &MintTokenResult{*commit}, nil

}
