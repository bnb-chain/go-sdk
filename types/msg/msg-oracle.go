package msg

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/binance-chain/go-sdk/common/rlp"
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

type CrossChainPackageType uint8

const (
	SynCrossChainPackageType     CrossChainPackageType = 0x00
	AckCrossChainPackageType     CrossChainPackageType = 0x01
	FailAckCrossChainPackageType CrossChainPackageType = 0x02
)

func noneExistPackageProto() interface{} {
	panic("should not exist such package")
}

// package type
var protoMetrics = map[sdk.IbcChannelID]map[CrossChainPackageType]func() interface{}{
	sdk.IbcChannelID(1): {
		SynCrossChainPackageType: func() interface{} {
			return new(ApproveBindSynPackage)
		},
		AckCrossChainPackageType: noneExistPackageProto,
		FailAckCrossChainPackageType: func() interface{} {
			return new(BindSynPackage)
		},
	},
	sdk.IbcChannelID(2): {
		SynCrossChainPackageType: noneExistPackageProto,
		AckCrossChainPackageType: func() interface{} {
			return new(TransferOutRefundPackage)
		},
		FailAckCrossChainPackageType: func() interface{} {
			return new(TransferOutSynPackage)
		},
	},
	sdk.IbcChannelID(3): {
		SynCrossChainPackageType: func() interface{} {
			return new(TransferInSynPackage)
		},
		AckCrossChainPackageType:     noneExistPackageProto,
		FailAckCrossChainPackageType: noneExistPackageProto,
	},
	sdk.IbcChannelID(8): {
		SynCrossChainPackageType: noneExistPackageProto,
		AckCrossChainPackageType: func() interface{} {
			return new(CommonAckPackage)
		},
		FailAckCrossChainPackageType: func() interface{} {
			return new(IbcValidatorSetPackage)
		},
	},
	sdk.IbcChannelID(9): {
		SynCrossChainPackageType: noneExistPackageProto,
		AckCrossChainPackageType: func() interface{} {
			return new(CommonAckPackage)
		},
		FailAckCrossChainPackageType: func() interface{} {
			return new(CrossParamChange)
		},
	},
	sdk.IbcChannelID(11): {
		SynCrossChainPackageType: func() interface{} {
			return new(SideDowntimeSlashPackage)
		},
		AckCrossChainPackageType:     noneExistPackageProto,
		FailAckCrossChainPackageType: noneExistPackageProto,
	},
}

type ApproveBindSynPackage struct {
	Status          uint32
	Bep2TokenSymbol [32]byte
}

type BindSynPackage struct {
	PackageType     uint8
	Bep2TokenSymbol [32]byte
	ContractAddr    [20]byte
	TotalSupply     *big.Int
	PeggyAmount     *big.Int
	Decimals        uint8
	ExpireTime      uint64
}

type TransferOutRefundPackage struct {
	Bep2TokenSymbol [32]byte
	RefundAmount    *big.Int
	RefundAddr      []byte
	RefundReason    uint32
}

type TransferOutSynPackage struct {
	Bep2TokenSymbol [32]byte
	ContractAddress [20]byte
	Amount          *big.Int
	Recipient       [20]byte
	RefundAddress   []byte
	ExpireTime      uint64
}

type TransferInSynPackage struct {
	Bep2TokenSymbol   [32]byte
	ContractAddress   [20]byte
	Amounts           []*big.Int
	ReceiverAddresses [][]byte
	RefundAddresses   [][20]byte
	ExpireTime        uint64
}

type CommonAckPackage struct {
	Code uint32
}

type IbcValidatorSetPackage struct {
	Type         uint8
	ValidatorSet []IbcValidator
}

type IbcValidator struct {
	ConsAddr []byte
	FeeAddr  []byte
	DistAddr []byte
	Power    uint64
}

type CrossParamChange struct {
	Key    string
	Value  []byte
	Target []byte
}

type SideDowntimeSlashPackage struct {
	SideConsAddr  []byte `json:"side_cons_addr"`
	SideHeight    uint64 `json:"side_height"`
	SideChainId   uint16 `json:"side_chain_id"`
	SideTimestamp uint64 `json:"side_timestamp"`
}

type CrossChainPackage struct {
	PackageType CrossChainPackageType
	RelayFee    big.Int
	Content     interface{}
}

func ParseClaimPayload(payload []byte) ([]CrossChainPackage, error) {
	packages := Packages{}
	err := rlp.DecodeBytes(payload, &packages)
	if err != nil {
		return nil, err
	}
	decodedPackage := make([]CrossChainPackage, 0, len(packages))
	for _, pack := range packages {
		ptype, relayerFee, err := DecodePackageHeader(pack.Payload)
		if err != nil {
			return nil, err
		}
		if _, exist := protoMetrics[pack.ChannelId]; !exist {
			return nil, fmt.Errorf("channnel id do not exist")
		}
		proto, exist := protoMetrics[pack.ChannelId][ptype]
		if !exist || proto == nil {
			return nil, fmt.Errorf("package type do not exist")
		}
		content := proto()
		err = rlp.DecodeBytes(pack.Payload[PackageHeaderLength:], content)
		if err != nil {
			return nil, err
		}
		decodedPackage = append(decodedPackage, CrossChainPackage{
			PackageType: ptype,
			RelayFee:    relayerFee,
			Content:     content,
		})
	}
	return decodedPackage, nil
}

func DecodePackageHeader(packageHeader []byte) (packageType CrossChainPackageType, relayFee big.Int, err error) {
	if len(packageHeader) < PackageHeaderLength {
		err = fmt.Errorf("length of packageHeader is less than %d", PackageHeaderLength)
		return
	}
	packageType = CrossChainPackageType(packageHeader[0])
	relayFee.SetBytes(packageHeader[PackageTypeLength : CrossChainFeeLength+PackageTypeLength])
	return
}
