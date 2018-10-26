package txmsg

import (
	"fmt"
)

// var _ sdk.TokenBurnMsg = (*TokenBurnMsg)(nil)

type TokenBurnMsg struct {
	MsgBase
}

func NewMsg(from AccAddress, symbol string, amount int64) TokenBurnMsg {
	return TokenBurnMsg{MsgBase{From: from, Symbol: symbol, Amount: amount}}
}

func (msg TokenBurnMsg) Type() string {
	return "tokensBurn"
}

func (msg TokenBurnMsg) String() string {
	return fmt.Sprintf("BurnMsg{%v#%v%v}", msg.From, msg.Amount, msg.Symbol)
}
