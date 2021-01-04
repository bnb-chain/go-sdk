package msg

import (
	"encoding/json"
	"fmt"

	"github.com/binance-chain/go-sdk/common/types"
)

// TransferOwnershipMsg def
type TransferOwnershipMsg struct {
	From     types.AccAddress `json:"from"`
	Symbol   string           `json:"symbol"`
	NewOwner types.AccAddress `json:"new_owner"`
}

func NewTransferOwnershipMsg(from types.AccAddress, symbol string, newOwner types.AccAddress) TransferOwnershipMsg {
	return TransferOwnershipMsg{
		From:     from,
		Symbol:   symbol,
		NewOwner: newOwner,
	}
}

func (msg TransferOwnershipMsg) Route() string  { return "tokensOwnershipTransfer" }
func (msg TransferOwnershipMsg) Type() string   { return "transferOwnership" }
func (msg TransferOwnershipMsg) String() string { return fmt.Sprintf("TransferOwnershipMsg{%#v}", msg) }

func (msg TransferOwnershipMsg) ValidateBasic() error {
	if len(msg.From) != types.AddrLen {
		return fmt.Errorf("Invalid from address, expected address length is %d, actual length is %d ", types.AddrLen, len(msg.From))
	}
	if len(msg.NewOwner) != types.AddrLen {
		return fmt.Errorf("Invalid newOwner, expected address length is %d, actual length is %d ", types.AddrLen, len(msg.NewOwner))
	}
	if !IsValidMiniTokenSymbol(msg.Symbol) {
		err := ValidateSymbol(msg.Symbol)
		if err != nil {
			return err
		}
	}
	return nil
}

func (msg TransferOwnershipMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg TransferOwnershipMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.From}
}

func (msg TransferOwnershipMsg) GetInvolvedAddresses() []types.AccAddress {
	return append(msg.GetSigners(), msg.NewOwner)
}
