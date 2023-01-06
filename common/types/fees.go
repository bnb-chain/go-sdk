package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	paramHubTypes "github.com/cosmos/cosmos-sdk/x/paramHub/types"
)

const (
	OperateFeeType  = paramHubTypes.OperateFeeType
	TransferFeeType = paramHubTypes.TransferFeeType
	DexFeeType      = paramHubTypes.DexFeeType

	FeeForProposer = types.FeeForProposer
	FeeForAll      = types.FeeForAll
	FeeFree        = types.FeeFree
)

type (
	FeeDistributeType = types.FeeDistributeType

	FeeParam         = paramHubTypes.FeeParam
	DexFeeParam      = paramHubTypes.DexFeeParam
	DexFeeField      = paramHubTypes.DexFeeField
	FixedFeeParams   = paramHubTypes.FixedFeeParams
	TransferFeeParam = paramHubTypes.TransferFeeParam
)
