package msg

import (
	"bytes"
	"fmt"

	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/common/types/bsc"
)

const (
	TypeSideChainSubmitEvidence = "submit_evidence"
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

type MsgBscSubmitEvidence struct {
	Submitter types.AccAddress `json:"submitter"`
	Headers   [2]*bsc.Header   `json:"headers"`
}

func NewMsgBscSubmitEvidence(submitter types.AccAddress, headers [2]*bsc.Header) MsgBscSubmitEvidence {

	return MsgBscSubmitEvidence{
		Submitter: submitter,
		Headers:   headers,
	}
}

func (MsgBscSubmitEvidence) Route() string {
	return SideChainSlashMsgRoute
}

func (MsgBscSubmitEvidence) Type() string {
	return TypeSideChainSubmitEvidence
}

func (msg MsgBscSubmitEvidence) ValidateBasic() error {
	if len(msg.Submitter) != types.AddrLen {
		return fmt.Errorf("Expected delegator address length is %d, actual length is %d", types.AddrLen, len(msg.Submitter))
	}

	if msg.Headers[0] == nil || msg.Headers[1] == nil {
		return fmt.Errorf("Both two block headers can not be nil")
	}
	if msg.Headers[0].Number != (msg.Headers[1].Number) {
		return fmt.Errorf("The numbers of two block headers are not the same")
	}
	if msg.Headers[0].ParentHash.Cmp(msg.Headers[1].ParentHash) != 0 {
		return fmt.Errorf("The parent hash of two block headers are not the same")
	}
	signature1, err := msg.Headers[0].GetSignature()
	if err != nil {
		return fmt.Errorf("Failed to get signature from block header, %s", err.Error())
	}
	signature2, err := msg.Headers[1].GetSignature()
	if err != nil {
		return fmt.Errorf("Failed to get signature from block header, %s", err.Error())
	}
	if bytes.Compare(signature1, signature2) == 0 {
		return fmt.Errorf("The two blocks are the same")
	}

	signer1, err := msg.Headers[0].ExtractSignerFromHeader()
	if err != nil {
		return fmt.Errorf("Failed to extract signer from block header, %s", err.Error())
	}
	signer2, err := msg.Headers[1].ExtractSignerFromHeader()
	if err != nil {
		return fmt.Errorf("Failed to extract signer from block header, %s", err.Error())
	}
	if bytes.Compare(signer1.Bytes(), signer2.Bytes()) != 0 {
		return fmt.Errorf("The signers of two block headers are not the same")
	}

	return nil
}

func (msg MsgBscSubmitEvidence) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return MustSortJSON(bz)
}

func (msg MsgBscSubmitEvidence) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.Submitter}
}

func (msg MsgBscSubmitEvidence) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}
