package transaction

import (
	"fmt"

	"github.com/binance-chain/go-sdk/types"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type SendTokenResult struct {
	tx.TxCommitResult
}

func (c *client) SendToken(dst types.AccAddress, symbol string, quantity int64, sync bool) (*SendTokenResult, error) {
	if symbol == "" {
		return nil, fmt.Errorf("symbol is missing. ")
	}
	fromAddr := c.keyManager.GetAddr()
	coins := types.Coins{types.Coin{Denom: symbol, Amount: quantity}}
	sendMsg := msg.CreateSendMsg(fromAddr, dst, coins)
	err := sendMsg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	commit, err := c.broadcastMsg(sendMsg, sync)
	if err != nil {
		return nil, err
	}
	return &SendTokenResult{*commit}, err

}
