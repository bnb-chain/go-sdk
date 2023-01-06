package msg

import (
	nTypes "github.com/bnb-chain/node/common/types"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	oracleTypes "github.com/cosmos/cosmos-sdk/x/oracle/types"
	stakeTypes "github.com/cosmos/cosmos-sdk/x/stake/types"
)

// constants
const (
	DotBSuffix              = nTypes.TokenSymbolDotBSuffix
	NativeToken             = nTypes.NativeTokenSymbol
	NativeTokenDotBSuffixed = nTypes.NativeTokenSymbolDotBSuffixed
	Decimals                = nTypes.TokenDecimals
	MaxTotalSupply          = nTypes.TokenMaxTotalSupply

	TokenSymbolMaxLen          = nTypes.TokenSymbolMaxLen
	TokenSymbolMinLen          = nTypes.TokenSymbolMinLen
	TokenSymbolTxHashSuffixLen = nTypes.TokenSymbolTxHashSuffixLen

	MiniTokenSymbolMaxLen          = nTypes.MiniTokenSymbolMaxLen
	MiniTokenSymbolMinLen          = nTypes.MiniTokenSymbolMinLen
	MiniTokenSymbolSuffixLen       = nTypes.MiniTokenSymbolSuffixLen
	MiniTokenSymbolMSuffix         = nTypes.MiniTokenSymbolMSuffix
	MiniTokenSymbolTxHashSuffixLen = nTypes.MiniTokenSymbolTxHashSuffixLen
	MaxMiniTokenNameLength         = 32
	MaxTokenURILength              = nTypes.MaxTokenURILength
)

type (
	Msg        = cTypes.Msg
	StatusText = oracleTypes.StatusText
)

const (
	PendingStatusText = oracleTypes.PendingStatusText
	SuccessStatusText = oracleTypes.SuccessStatusText
	FailedStatusText  = oracleTypes.FailedStatusText
)

var (
	StatusTextToString = oracleTypes.StatusTextToString
	StringToStatusText = oracleTypes.StringToStatusText
)

type (
	Status        = oracleTypes.Status
	Prophecy      = oracleTypes.Prophecy
	DBProphecy    = oracleTypes.DBProphecy
	OracleRelayer = stakeTypes.OracleRelayer
)

var (
	IsValidMiniTokenSymbol  = nTypes.IsValidMiniTokenSymbol
	ValidateMiniTokenSymbol = nTypes.ValidateMiniTokenSymbol
	ValidateSymbol          = nTypes.ValidateTokenSymbol
)
