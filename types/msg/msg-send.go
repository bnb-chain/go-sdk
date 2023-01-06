package msg

import (
	"github.com/bnb-chain/go-sdk/common/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type (
	SendMsg = bank.MsgSend
	Input   = bank.Input
	Output  = bank.Output
)

var (
	NewMsgSend = bank.NewMsgSend
	NewInput   = bank.NewInput
	NewOutput  = bank.NewOutput
)

type Transfer struct {
	ToAddr types.AccAddress
	Coins  types.Coins
}

func CreateSendMsg(from types.AccAddress, fromCoins types.Coins, transfers []Transfer) SendMsg {
	input := NewInput(from, fromCoins)

	output := make([]Output, 0, len(transfers))
	for _, t := range transfers {
		t.Coins = t.Coins.Sort()
		output = append(output, NewOutput(t.ToAddr, t.Coins))
	}
	msg := NewMsgSend([]Input{input}, output)
	return msg
}
