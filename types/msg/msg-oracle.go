package msg

import (
	"encoding/json"
	"fmt"

	sdk "github.com/binance-chain/go-sdk/common/types"
)

const (
	RouteOracle = "oracle"

	ClaimMsgType = "oracleClaim"

	OracleChannelId sdk.IbcChannelID = 0x00
)

const (
	CrossChainFeeLength = 32
	PackageTypeLength   = 1
	PackageHeaderLength = CrossChainFeeLength + PackageTypeLength
)

func GetClaimId(chainId sdk.IbcChainID, channelId sdk.IbcChannelID, sequence int64) string {
	return fmt.Sprintf("%d:%d:%d", chainId, channelId, sequence)
}

// Claim contains an arbitrary claim with arbitrary content made by a given validator
type Claim struct {
	ID               string         `json:"id"`
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	Content          string         `json:"content"`
}

// NewClaim returns a new Claim
func NewClaim(id string, validatorAddress sdk.ValAddress, content string) Claim {
	return Claim{
		ID:               id,
		ValidatorAddress: validatorAddress,
		Content:          content,
	}
}

type ClaimMsg struct {
	ChainId          sdk.IbcChainID `json:"chain_id"`
	Sequence         uint64         `json:"sequence"`
	Payload          []byte         `json:"payload"`
	ValidatorAddress sdk.AccAddress `json:"validator_address"`
}

type Packages []Package

type Package struct {
	ChannelId sdk.IbcChannelID
	Sequence  uint64
	Payload   []byte
}

func NewClaimMsg(ChainId sdk.IbcChainID, sequence uint64, payload []byte, validatorAddr sdk.AccAddress) ClaimMsg {
	return ClaimMsg{
		ChainId:          ChainId,
		Sequence:         sequence,
		Payload:          payload,
		ValidatorAddress: validatorAddr,
	}
}

// nolint
func (msg ClaimMsg) Route() string { return RouteOracle }
func (msg ClaimMsg) Type() string  { return ClaimMsgType }
func (msg ClaimMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.ValidatorAddress}
}

func (msg ClaimMsg) String() string {
	return fmt.Sprintf("Claim{%v#%v#%v%v%x}",
		msg.ChainId, msg.Sequence, msg.ValidatorAddress.String(), msg.Payload)
}

// GetSignBytes - Get the bytes for the message signer to sign on
func (msg ClaimMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ClaimMsg) GetInvolvedAddresses() []sdk.AccAddress {
	return msg.GetSigners()
}

// ValidateBasic is used to quickly disqualify obviously invalid messages quickly
func (msg ClaimMsg) ValidateBasic() error {
	if len(msg.Payload) < PackageHeaderLength {
		return fmt.Errorf("length of payload is less than %d", PackageHeaderLength)
	}
	if len(msg.ValidatorAddress) != sdk.AddrLen {
		return fmt.Errorf("address length should be %d", sdk.AddrLen)
	}
	return nil
}
