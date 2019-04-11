package msg

import (
	"encoding/json"
	"fmt"

	"github.com/binance-chain/go-sdk/common/types"
)

// TokenBurnMsg def
type TokenBurnMsg struct {
	From   types.AccAddress `json:"from"`
	Symbol string           `json:"symbol"`
	Amount int64            `json:"amount"`
}

// NewMsg for instance creation
func NewTokenBurnMsg(from types.AccAddress, symbol string, amount int64) TokenBurnMsg {
	return TokenBurnMsg{From: from, Symbol: symbol, Amount: amount}
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

func (msg TokenBurnMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners is part of Msg interface
func (msg TokenBurnMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.From}
}

// ValidateBasic does a simple validation check that
// doesn't require access to any other information.
func (msg TokenBurnMsg) ValidateBasic() error {
	err := ValidateSymbol(msg.Symbol)
	if err != nil {
		return fmt.Errorf("ErrInvalidCoins %s", msg.Symbol)
	}

	if msg.Amount <= 0 {
		return fmt.Errorf("ErrInsufficientFunds, amount should be more than 0")
	}

	return nil
}
