package msg

import (
	"encoding/json"
	"fmt"

	"github.com/binance-chain/go-sdk/common/types"
)

// TokenFreezeMsg def
type TokenFreezeMsg struct {
	From   types.AccAddress `json:"from"`
	Symbol string           `json:"symbol"`
	Amount int64            `json:"amount"`
}

// NewFreezeMsg for instance creation
func NewFreezeMsg(from types.AccAddress, symbol string, amount int64) TokenFreezeMsg {
	return TokenFreezeMsg{From: from, Symbol: symbol, Amount: amount}
}

// Route is part of Msg interface
func (msg TokenFreezeMsg) Route() string { return "tokensFreeze" }

// Type is part of Msg interface
func (msg TokenFreezeMsg) Type() string { return "tokensFreeze" }

// String is part of Msg interface
func (msg TokenFreezeMsg) String() string {
	return fmt.Sprintf("Freeze{%v#%v}", msg.From, msg.Symbol)
}

// GetInvolvedAddresses is part of Msg interface
func (msg TokenFreezeMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}

func (msg TokenFreezeMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.From}
}

func (msg TokenFreezeMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg TokenFreezeMsg) ValidateBasic() error {
	err := ValidateSymbol(msg.Symbol)
	if err != nil {
		return fmt.Errorf("ErrInvalidCoins %s", msg.Symbol)
	}

	if msg.Amount <= 0 {
		return fmt.Errorf("ErrInsufficientFunds, amount should be more than 0")
	}

	return nil
}

// TokenUnfreezeMsg def
type TokenUnfreezeMsg struct {
	From   types.AccAddress `json:"from"`
	Symbol string           `json:"symbol"`
	Amount int64            `json:"amount"`
}

// NewUnfreezeMsg for instance creation
func NewUnfreezeMsg(from types.AccAddress, symbol string, amount int64) TokenUnfreezeMsg {
	return TokenUnfreezeMsg{From: from, Symbol: symbol, Amount: amount}
}

// Route is part of Msg interface
func (msg TokenUnfreezeMsg) Route() string { return "tokensFreeze" }

// Type is part of Msg interface
func (msg TokenUnfreezeMsg) Type() string { return "tokensFreeze" }

// String is part of Msg interface
func (msg TokenUnfreezeMsg) String() string {
	return fmt.Sprintf("Unfreeze{%v#%v%v}", msg.From, msg.Amount, msg.Symbol)
}

// GetInvolvedAddresses is part of Msg interface
func (msg TokenUnfreezeMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}

func (msg TokenUnfreezeMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.From}
}

func (msg TokenUnfreezeMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg TokenUnfreezeMsg) ValidateBasic() error {
	err := ValidateSymbol(msg.Symbol)
	if err != nil {
		return fmt.Errorf("ErrInvalidCoins %s", msg.Symbol)
	}

	if msg.Amount <= 0 {
		return fmt.Errorf("ErrInsufficientFunds, amount should be more than 0")
	}

	return nil
}
