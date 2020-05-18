package msg

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	sdk "github.com/binance-chain/go-sdk/common/types"
)

const (
	RouteOracle = "oracle"

	ClaimMsgType = "oracleClaim"
)

// Type that represents Claim Type as a byte
type ClaimType byte

const (
	ClaimTypeSkipSequence      ClaimType = 0x1
	ClaimTypeUpdateBind        ClaimType = 0x2
	ClaimTypeTransferOutRefund ClaimType = 0x3
	ClaimTypeTransferIn        ClaimType = 0x4
	ClaimTypeDowntimeSlash     ClaimType = 0x5

	ClaimTypeSkipSequenceName      = "SkipSequence"
	ClaimTypeUpdateBindName        = "UpdateBind"
	ClaimTypeTransferOutRefundName = "TransferOutRefund"
	ClaimTypeTransferInName        = "TransferIn"
	ClaimTypeDowntimeSlashName     = "DowntimeSlash"
)

var claimTypeToName = map[ClaimType]string{
	ClaimTypeSkipSequence:      ClaimTypeSkipSequenceName,
	ClaimTypeUpdateBind:        ClaimTypeUpdateBindName,
	ClaimTypeTransferOutRefund: ClaimTypeTransferOutRefundName,
	ClaimTypeTransferIn:        ClaimTypeTransferInName,
	ClaimTypeDowntimeSlash:     ClaimTypeDowntimeSlashName,
}

var claimNameToType = map[string]ClaimType{
	ClaimTypeSkipSequenceName:      ClaimTypeSkipSequence,
	ClaimTypeUpdateBindName:        ClaimTypeUpdateBind,
	ClaimTypeTransferOutRefundName: ClaimTypeTransferOutRefund,
	ClaimTypeTransferInName:        ClaimTypeTransferIn,
	ClaimTypeDowntimeSlashName:     ClaimTypeDowntimeSlash,
}

// String to claim type byte.  Returns ff if invalid.
func ClaimTypeFromString(str string) (ClaimType, error) {
	claimType, ok := claimNameToType[str]
	if !ok {
		return ClaimType(0xff), errors.Errorf("'%s' is not a valid claim type", str)
	}
	return claimType, nil
}

func ClaimTypeToString(typ ClaimType) string {
	return claimTypeToName[typ]
}

func IsValidClaimType(ct ClaimType) bool {
	if _, ok := claimTypeToName[ct]; ok {
		return true
	}
	return false
}

func GetClaimId(claimType ClaimType, sequence int64) string {
	return fmt.Sprintf("%d:%d", claimType, sequence)
}

var claimTypeSequencePrefix = []byte("claimTypeSeq:")

func GetClaimTypeSequence(claimType ClaimType) []byte {
	return append(claimTypeSequencePrefix, byte(claimType))
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
	ClaimType        ClaimType      `json:"claim_type"`
	Sequence         int64          `json:"sequence"`
	Claim            string         `json:"claim"`
	ValidatorAddress sdk.AccAddress `json:"validator_address"`
}

func NewClaimMsg(claimType ClaimType, sequence int64, claim string, validatorAddr sdk.AccAddress) ClaimMsg {
	return ClaimMsg{
		ClaimType:        claimType,
		Sequence:         sequence,
		Claim:            claim,
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
	return fmt.Sprintf("Claim{%v#%v#%v#%v}",
		msg.ClaimType, msg.Sequence, msg.Claim, msg.ValidatorAddress.String())
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
	if !IsValidClaimType(msg.ClaimType) {
		return fmt.Errorf("claim type %v does not exist", msg.ClaimType)
	}

	if msg.Sequence < 0 {
		return fmt.Errorf("sequence should not be less than 0")
	}

	if len(msg.Claim) == 0 {
		return fmt.Errorf("claim should not be empty")
	}

	if len(msg.ValidatorAddress) != sdk.AddrLen {
		return fmt.Errorf("length of validator address should be %d", sdk.AddrLen)
	}
	return nil
}

type SideDowntimeSlashClaim struct {
	SideConsAddr  []byte `json:"side_cons_addr"`
	SideHeight    int64  `json:"side_height"`
	SideChainId   string `json:"side_chain_id"`
	SideTimestamp int64  `json:"side_timestamp"`
}
