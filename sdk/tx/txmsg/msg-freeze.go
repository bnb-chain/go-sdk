package txmsg

import (
	"fmt"
	// "github.com/BiJie/BinanceChain/plugins/tokens/base"
)

// var _ sdk.Msg = (*TokenFreezeMsg)(nil)

type TokenFreezeMsg struct {
	MsgBase
}

func NewFreezeMsg(from AccAddress, symbol string, amount int64) TokenFreezeMsg {
	return TokenFreezeMsg{MsgBase{From: from, Symbol: symbol, Amount: amount}}
}

func (msg TokenFreezeMsg) Type() string { return "tokensFreeze" }

func (msg TokenFreezeMsg) String() string {
	return fmt.Sprintf("Freeze{%v#%v}", msg.From, msg.Symbol)
}

// var _ sdk.Msg = (*TokenUnfreezeMsg)(nil)

type TokenUnfreezeMsg struct {
	MsgBase
}

func NewUnfreezeMsg(from AccAddress, symbol string, amount int64) TokenUnfreezeMsg {
	return TokenUnfreezeMsg{MsgBase{From: from, Symbol: symbol, Amount: amount}}
}

func (msg TokenUnfreezeMsg) Type() string { return "tokensFreeze" }

func (msg TokenUnfreezeMsg) String() string {
	return fmt.Sprintf("Unfreeze{%v#%v%v}", msg.From, msg.Amount, msg.Symbol)
}
