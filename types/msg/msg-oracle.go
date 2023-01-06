package msg

import (
	"fmt"
	"math/big"

	"github.com/bnb-chain/go-sdk/common/rlp"
	bridgeTypes "github.com/bnb-chain/node/plugins/bridge/types"
	"github.com/cosmos/cosmos-sdk/types"
	oracleTypes "github.com/cosmos/cosmos-sdk/x/oracle/types"
	paramHubTypes "github.com/cosmos/cosmos-sdk/x/paramHub/types"
	sidechainTypes "github.com/cosmos/cosmos-sdk/x/sidechain/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing"
	crossStake "github.com/cosmos/cosmos-sdk/x/stake/cross_stake"
	stakeTypes "github.com/cosmos/cosmos-sdk/x/stake/types"

	sdk "github.com/bnb-chain/go-sdk/common/types"
)

const (
	RouteOracle     = oracleTypes.RouteOracle
	ClaimMsgType    = oracleTypes.ClaimMsgType
	OracleChannelId = oracleTypes.RelayPackagesChannelId
)

const (
	CrossChainFeeLength = sidechainTypes.CrossChainFeeLength
	PackageTypeLength   = sidechainTypes.PackageTypeLength
	PackageHeaderLength = sidechainTypes.PackageHeaderLength
)

var (
	GetClaimId  = oracleTypes.GetClaimId
	NewClaim    = oracleTypes.NewClaim
	NewClaimMsg = oracleTypes.NewClaimMsg
)

type (
	Claim    = oracleTypes.Claim
	ClaimMsg = oracleTypes.ClaimMsg
	Package  = oracleTypes.Package
	Packages = oracleTypes.Packages

	CrossChainPackageType = types.CrossChainPackageType
)

const (
	SynCrossChainPackageType     = types.SynCrossChainPackageType
	AckCrossChainPackageType     = types.AckCrossChainPackageType
	FailAckCrossChainPackageType = types.FailAckCrossChainPackageType
)

type (
	ApproveBindSynPackage       = bridgeTypes.ApproveBindSynPackage
	BindSynPackage              = bridgeTypes.BindSynPackage
	TransferOutRefundPackage    = bridgeTypes.TransferOutRefundPackage
	TransferOutSynPackage       = bridgeTypes.TransferOutSynPackage
	TransferInSynPackage        = bridgeTypes.TransferInSynPackage
	MirrorSynPackage            = bridgeTypes.MirrorSynPackage
	MirrorSyncSynPackage        = bridgeTypes.MirrorSyncSynPackage
	CommonAckPackage            = sidechainTypes.CommonAckPackage
	IbcValidatorSetPackage      = stakeTypes.IbcValidatorSetPackage
	IbcValidator                = stakeTypes.IbcValidator
	CrossParamChange            = paramHubTypes.CSCParamChange
	SideDowntimeSlashPackage    = slashingTypes.SideDowntimeSlashPackage
	CrossStakeSynPackageFromBSC = crossStake.CrossStakeSynPackageFromBSC
	CrossStakeRefundPackage     = stakeTypes.CrossStakeRefundPackage
)

type CrossChainPackage struct {
	PackageType CrossChainPackageType
	RelayFee    big.Int
	Content     interface{}
}

var (
	DecodePackageHeader = sidechainTypes.DecodePackageHeader
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
	sdk.IbcChannelID(4): {
		SynCrossChainPackageType: func() interface{} {
			return new(MirrorSynPackage)
		},
		AckCrossChainPackageType:     noneExistPackageProto,
		FailAckCrossChainPackageType: noneExistPackageProto,
	},
	sdk.IbcChannelID(5): {
		SynCrossChainPackageType: func() interface{} {
			return new(MirrorSyncSynPackage)
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
	sdk.IbcChannelID(16): {
		SynCrossChainPackageType: func() interface{} {
			return new(CrossStakeSynPackageFromBSC)
		},
		AckCrossChainPackageType: func() interface{} {
			return new(CrossStakeRefundPackage)
		},
		FailAckCrossChainPackageType: noneExistPackageProto,
	},
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
