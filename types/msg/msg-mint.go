package msg

import (
	"encoding/json"
	"fmt"

	"github.com/binance-chain/go-sdk/common/types"
)

type MintMsg struct {
	From   types.AccAddress `json:"from"`
	Symbol string           `json:"symbol"`
	Amount int64            `json:"amount"`
}

func NewMintMsg(from types.AccAddress, symbol string, amount int64) MintMsg {
	return MintMsg{
		From:   from,
		Symbol: symbol,
		Amount: amount,
	}
}

func (msg MintMsg) ValidateBasic() error {
	if msg.From == nil {
		return fmt.Errorf("sender address cannot be empty")
	}
	if msg.Amount <= 0 {
		return fmt.Errorf("Amount cant be less than 0 ")
	}

	return nil
}

// Implements MintMsg.
func (msg MintMsg) Route() string                  { return "tokensIssue" }
func (msg MintMsg) Type() string                   { return "mintMsg" }
func (msg MintMsg) String() string                 { return fmt.Sprintf("MintMsg{%#v}", msg) }
func (msg MintMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }
func (msg MintMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}
func (msg MintMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}
