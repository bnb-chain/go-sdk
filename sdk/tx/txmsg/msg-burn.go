package txmsg

import (
	"fmt"
)

// TokenBurnMsg def
type TokenBurnMsg struct {
	MsgBase
}

// NewMsg for instance creation
func NewMsg(from AccAddress, symbol string, amount int64) TokenBurnMsg {
	return TokenBurnMsg{MsgBase{From: from, Symbol: symbol, Amount: amount}}
}

// Type part of Msg interface
func (msg TokenBurnMsg) Type() string {
	return "tokensBurn"
}

// String part of Msg interface
func (msg TokenBurnMsg) String() string {
	return fmt.Sprintf("BurnMsg{%v#%v%v}", msg.From, msg.Amount, msg.Symbol)
}
