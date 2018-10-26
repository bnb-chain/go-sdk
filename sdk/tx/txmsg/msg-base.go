package txmsg

import (
	"encoding/json"
	"fmt"
)

type MsgBase struct {
	From   AccAddress `json:"from"`
	Symbol string     `json:"symbol"`
	Amount int64      `json:"amount"`
}

func (msg MsgBase) Type() string {
	return ""
}

// ValidateBasic does a simple validation check that
// doesn't require access to any other information.
func (msg MsgBase) ValidateBasic() error {
	// err := types.ValidateSymbol(msg.Symbol)
	// if err != nil {
	// 	return sdk.ErrInvalidCoins(err.Error())
	// }

	// if msg.Amount <= 0 {
	// 	// TODO: maybe we need to define our own errors
	// 	return sdk.ErrInsufficientFunds("amount should be more than 0")
	// }

	return nil
}

func (msg MsgBase) String() string {
	return fmt.Sprintf("MsgBase{%v#%v%v}", msg.From, msg.Amount, msg.Symbol)
}

func (msg MsgBase) Get(key interface{}) (value interface{}) {
	return nil
}

func (msg MsgBase) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg MsgBase) GetSigners() []AccAddress {
	return []AccAddress{msg.From}
}
