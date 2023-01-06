package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake"
	stakeTypes "github.com/cosmos/cosmos-sdk/x/stake/types"
)

// nolint
const (
	Unbonded  = types.Unbonded
	Unbonding = types.Unbonding
	Bonded    = types.Bonded
)

type (
	Commission          = stakeTypes.Commission
	CommissionMsg       = stakeTypes.CommissionMsg
	Description         = stakeTypes.Description
	Validator           = stakeTypes.Validator
	UnbondingDelegation = stakeTypes.UnbondingDelegation
	Pool                = stakeTypes.Pool

	ValAddress  = types.ValAddress
	BondStatus  = types.BondStatus
	ConsAddress = types.ConsAddress

	QueryTopValidatorsParams = stake.QueryTopValidatorsParams
	QueryBondsParams         = stake.QueryBondsParams
	QueryValidatorParams     = stake.QueryValidatorParams
)

var (
	NewCommission = stakeTypes.NewCommission

	ValAddressFromBech32  = types.ValAddressFromBech32
	ConsAddressFromHex    = types.ConsAddressFromHex
	ConsAddressFromBech32 = types.ConsAddressFromBech32
	GetConsAddress        = types.GetConsAddress

	NewBaseParams = stake.NewBaseParams
)
