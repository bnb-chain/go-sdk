package msg

import (
	"fmt"

	"github.com/binance-chain/go-sdk/common/types"
)

const (
	TypeMsgSideChainUnjail      = "side_chain_unjail"

	SideChainSlashMsgRoute = "slashing"
)

type MsgSideChainUnjail struct {
	ValidatorAddr types.ValAddress `json:"address"`
	SideChainId   string           `json:"side_chain_id"`
}

func NewMsgSideChainUnjail(validatorAddr types.ValAddress, sideChainId string) MsgSideChainUnjail {
	return MsgSideChainUnjail{
		ValidatorAddr: validatorAddr,
		SideChainId:   sideChainId,
	}
}

func (msg MsgSideChainUnjail) Route() string { return SideChainSlashMsgRoute }
func (msg MsgSideChainUnjail) Type() string  { return TypeMsgSideChainUnjail }
func (msg MsgSideChainUnjail) GetSigners() []types.AccAddress {
	return []types.AccAddress{types.AccAddress(msg.ValidatorAddr)}
}

func (msg MsgSideChainUnjail) GetSignBytes() []byte {
	b := MsgCdc.MustMarshalJSON(msg)
	return MustSortJSON(b)
}

func (msg MsgSideChainUnjail) ValidateBasic() error {
	if msg.ValidatorAddr == nil {
		return fmt.Errorf("validator does not exist for that address")
	}
	if len(msg.SideChainId) == 0 || len(msg.SideChainId) > MaxSideChainIdLength {
		return fmt.Errorf(fmt.Sprintf("side chain id must be included and max length is %d bytes", MaxSideChainIdLength))
	}
	return nil
}

func (msg MsgSideChainUnjail) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}
