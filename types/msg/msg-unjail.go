package msg

import (
	"fmt"
	sdk "github.com/binance-chain/go-sdk/common/types"
)

// name to identify transaction types
const (
	TypeMsgUnjail = "unjail"
	SlashMsgRoute = "slashing"
)

// MsgUnjail - struct for unjailing jailed validator
type MsgUnjail struct {
	ValidatorAddr sdk.ValAddress `json:"address"` // address of the validator operator
}

func NewMsgUnjail(validatorAddr sdk.ValAddress) MsgUnjail {
	return MsgUnjail{
		ValidatorAddr: validatorAddr,
	}
}

//nolint
func (msg MsgUnjail) Route() string { return SlashMsgRoute }
func (msg MsgUnjail) Type() string  { return TypeMsgUnjail }
func (msg MsgUnjail) GetSigners() []sdk.AccAddress {
	res := []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr)}
	fmt.Printf("MsgUnjail.GetSigners: %+v\n", res)
	return res
}

// get the bytes for the message signer to sign on
func (msg MsgUnjail) GetSignBytes() []byte {
	b := MsgCdc.MustMarshalJSON(msg)
	res := MustSortJSON(b)
	fmt.Printf("MsgUnjail.GetSignBytes: %X\n", res)
	return res
}

// quick validity check
func (msg MsgUnjail) ValidateBasic() error {
	if len(msg.ValidatorAddr) != sdk.AddrLen {
		return fmt.Errorf("validator does not exist for that address")
	}
	return nil
}

func (msg MsgUnjail) GetInvolvedAddresses() []sdk.AccAddress {
	return msg.GetSigners()
}
