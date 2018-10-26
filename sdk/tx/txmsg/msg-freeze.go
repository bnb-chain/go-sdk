package txmsg

import (
	"fmt"
)

// TokenFreezeMsg def
type TokenFreezeMsg struct {
	MsgBase
}

// NewFreezeMsg for instance creation
func NewFreezeMsg(from AccAddress, symbol string, amount int64) TokenFreezeMsg {
	return TokenFreezeMsg{MsgBase{From: from, Symbol: symbol, Amount: amount}}
}

// Type is part of Msg interface
func (msg TokenFreezeMsg) Type() string { return "tokensFreeze" }

// String is part of Msg interface
func (msg TokenFreezeMsg) String() string {
	return fmt.Sprintf("Freeze{%v#%v}", msg.From, msg.Symbol)
}

// TokenUnfreezeMsg def
type TokenUnfreezeMsg struct {
	MsgBase
}

// NewUnfreezeMsg for instance creation
func NewUnfreezeMsg(from AccAddress, symbol string, amount int64) TokenUnfreezeMsg {
	return TokenUnfreezeMsg{MsgBase{From: from, Symbol: symbol, Amount: amount}}
}

// Type is part of Msg interface
func (msg TokenUnfreezeMsg) Type() string { return "tokensFreeze" }

// String is part of Msg interface
func (msg TokenUnfreezeMsg) String() string {
	return fmt.Sprintf("Unfreeze{%v#%v%v}", msg.From, msg.Amount, msg.Symbol)
}
