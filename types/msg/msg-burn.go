package msg

import (
	"fmt"
	"github.com/binance-chain/go-sdk/types"
)

// TokenBurnMsg def
type TokenBurnMsg struct {
	MsgBase
}

// NewMsg for instance creation
func NewTokenBurnMsg(from types.AccAddress, symbol string, amount int64) TokenBurnMsg {
	return TokenBurnMsg{MsgBase{From: from, Symbol: symbol, Amount: amount}}
}

// Route is part of Msg interface
func (msg TokenBurnMsg) Route() string {
	return "tokensBurn"
}

// Type is part of Msg interface
func (msg TokenBurnMsg) Type() string {
	return "tokensBurn"
}

// String is part of Msg interface
func (msg TokenBurnMsg) String() string {
	return fmt.Sprintf("BurnMsg{%v#%v%v}", msg.From, msg.Amount, msg.Symbol)
}

// GetInvolvedAddresses is part of Msg interface
func (msg TokenBurnMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}
