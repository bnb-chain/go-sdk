package api

import (
	"fmt"
	"github.com/BiJie/bnc-go-sdk/sdk/tx/txmsg"
)

type SendTokenResult struct {
	TxCommitResult
}

func (dex *dexAPI) SendToken(dst txmsg.AccAddress, symbol string, quantity int64, sync bool) (*SendTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("symbol is missing. ")
	}
	fromAddr := dex.keyManager.GetAddr()
	coins := txmsg.Coins{txmsg.Coin{Denom:symbol, Amount:quantity}}
	sendMsg := txmsg.CreateSendMsg(fromAddr, dst, coins)
	err := sendMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := dex.broadcastMsg(sendMsg, sync)
	if err != nil {
		return nil, err
	}
	return &SendTokenResult{*commit}, err

}
