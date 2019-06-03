package transaction

import (
	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type SendTokenResult struct {
	tx.TxCommitResult
}

func (c *client) SendToken(transfers []msg.Transfer, sync bool, options ...Option) (*SendTokenResult, error) {
	fromAddr := c.keyManager.GetAddr()
	fromCoins := types.Coins{}
	for _, t := range transfers {
		fromCoins = fromCoins.Plus(t.Coins)
	}
	sendMsg := msg.CreateSendMsg(fromAddr, fromCoins, transfers)
	err := sendMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(sendMsg, sync, options...)
	if err != nil {
		return nil, err
	}
	return &SendTokenResult{*commit}, err

}
