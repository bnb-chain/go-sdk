package msg

import (
	"math/big"

	sdk "github.com/bnb-chain/go-sdk/common/types"
	bridgeTypes "github.com/bnb-chain/node/plugins/bridge/types"
	"github.com/bnb-chain/node/plugins/dex/order"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	oracleTypes "github.com/cosmos/cosmos-sdk/x/oracle/types"
	paramHubTypes "github.com/cosmos/cosmos-sdk/x/paramHub/types"
	sidechainTypes "github.com/cosmos/cosmos-sdk/x/sidechain/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
	crossStake "github.com/cosmos/cosmos-sdk/x/stake/cross_stake"
	stakeTypes "github.com/cosmos/cosmos-sdk/x/stake/types"
)

type (
	Msg = cTypes.Msg
)

// ===================  gov module ====================
const (
	OptionEmpty      = gov.OptionEmpty
	OptionYes        = gov.OptionYes
	OptionAbstain    = gov.OptionAbstain
	OptionNo         = gov.OptionNo
	OptionNoWithVeto = gov.OptionNoWithVeto

	ProposalTypeNil             = gov.ProposalTypeNil
	ProposalTypeText            = gov.ProposalTypeText
	ProposalTypeParameterChange = gov.ProposalTypeParameterChange
	ProposalTypeSoftwareUpgrade = gov.ProposalTypeSoftwareUpgrade
	ProposalTypeListTradingPair = gov.ProposalTypeListTradingPair
	ProposalTypeFeeChange       = gov.ProposalTypeFeeChange

	ProposalTypeSCParamsChange  = gov.ProposalTypeSCParamsChange
	ProposalTypeCSCParamsChange = gov.ProposalTypeCSCParamsChange
)

type (
	VoteOption   = gov.VoteOption
	ProposalKind = gov.ProposalKind

	ListTradingPairParams = gov.ListTradingPairParams
)

type (
	SCParam        = paramHubTypes.SCParam
	SCChangeParams = paramHubTypes.SCChangeParams
	CSCParamChange = paramHubTypes.CSCParamChange
	IbcParams      = ibc.Params
	OracleParams   = oracleTypes.Params
	SlashParams    = slashingTypes.Params
	StakeParams    = stake.Params
)

// ===================  trade module ====================
var (
	OrderSide       = order.Side
	GenerateOrderID = order.GenerateOrderID
)

// ===================  oracle module ====================
const (
	OracleChannelId     = oracleTypes.RelayPackagesChannelId
	PackageHeaderLength = sidechainTypes.PackageHeaderLength

	SynCrossChainPackageType     = cTypes.SynCrossChainPackageType
	AckCrossChainPackageType     = cTypes.AckCrossChainPackageType
	FailAckCrossChainPackageType = cTypes.FailAckCrossChainPackageType
)

var (
	GetClaimId          = oracleTypes.GetClaimId
	DecodePackageHeader = sidechainTypes.DecodePackageHeader
)

type (
	Package  = oracleTypes.Package
	Packages = oracleTypes.Packages

	CrossChainPackageType = cTypes.CrossChainPackageType
)

type (
	Status        = oracleTypes.Status
	Prophecy      = oracleTypes.Prophecy
	DBProphecy    = oracleTypes.DBProphecy
	OracleRelayer = stakeTypes.OracleRelayer
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
	SideDowntimeSlashPackage    = slashingTypes.SideSlashPackage
	CrossStakeSynPackageFromBSC = crossStake.CrossStakeSynPackageFromBSC
	CrossStakeRefundPackage     = stakeTypes.CrossStakeRefundPackage
)

type CrossChainPackage struct {
	PackageType CrossChainPackageType
	RelayFee    big.Int
	Content     interface{}
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

// ===================  bank module ====================
type (
	Input  = bank.Input
	Output = bank.Output
)

var (
	NewInput  = bank.NewInput
	NewOutput = bank.NewOutput
)

type Transfer struct {
	ToAddr cTypes.AccAddress
	Coins  cTypes.Coins
}

// ===================  staking module ====================
type (
	Description = stakeTypes.Description
)
