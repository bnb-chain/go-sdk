package msg

import (
	"encoding/json"
	"fmt"
	"github.com/binance-chain/go-sdk/types"
)

// MsgBase def
type MsgBase struct {
	From   types.AccAddress `json:"from"`
	Symbol string           `json:"symbol"`
	Amount int64            `json:"amount"`
}

// Type is part of Msg interface
func (msg MsgBase) Type() string {
	return ""
}

// String is part of Msg interface
func (msg MsgBase) String() string {
	return fmt.Sprintf("MsgBase{%v#%v%v}", msg.From, msg.Amount, msg.Symbol)
}

// Get is part of Msg interface
func (msg MsgBase) Get(key interface{}) (value interface{}) {
	return nil
}

// GetSignBytes is part of Msg interface
func (msg MsgBase) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners is part of Msg interface
func (msg MsgBase) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.From}
}

// ValidateBasic does a simple validation check that
// doesn't require access to any other information.
func (msg MsgBase) ValidateBasic() error {
	err := ValidateSymbol(msg.Symbol)
	if err != nil {
		return fmt.Errorf("ErrInvalidCoins %s", msg.Symbol)
	}

	if msg.Amount <= 0 {
		return fmt.Errorf("ErrInsufficientFunds, amount should be more than 0")
	}

	return nil
}
