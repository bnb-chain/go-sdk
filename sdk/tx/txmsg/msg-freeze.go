package txmsg

import (
	"fmt"
	// "github.com/BiJie/BinanceChain/plugins/tokens/base"
)

// var _ sdk.Msg = (*FreezeMsg)(nil)

type FreezeMsg struct {
	MsgBase
}

func NewFreezeMsg(from AccAddress, symbol string, amount int64) FreezeMsg {
	return FreezeMsg{MsgBase{From: from, Symbol: symbol, Amount: amount}}
}

func (msg FreezeMsg) Type() string { return "tokensFreeze" }

func (msg FreezeMsg) String() string {
	return fmt.Sprintf("Freeze{%v#%v}", msg.From, msg.Symbol)
}

// var _ sdk.Msg = (*UnfreezeMsg)(nil)

type UnfreezeMsg struct {
	MsgBase
}

func NewUnfreezeMsg(from AccAddress, symbol string, amount int64) UnfreezeMsg {
	return UnfreezeMsg{MsgBase{From: from, Symbol: symbol, Amount: amount}}
}

func (msg UnfreezeMsg) Type() string { return "tokensFreeze" }

func (msg UnfreezeMsg) String() string {
	return fmt.Sprintf("Unfreeze{%v#%v%v}", msg.From, msg.Amount, msg.Symbol)
}
