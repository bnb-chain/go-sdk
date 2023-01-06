package types

import (
	"github.com/cosmos/cosmos-sdk/x/stake/querier"
	stakeTypes "github.com/cosmos/cosmos-sdk/x/stake/types"
)

type (
	Delegation         = stakeTypes.Delegation
	Redelegation       = stakeTypes.Redelegation
	DelegationResponse = stakeTypes.DelegationResponse

	QueryDelegatorParams    = querier.QueryDelegatorParams
	QueryRedelegationParams = querier.QueryRedelegationParams
)

var (
	UnmarshalRED = stakeTypes.UnmarshalRED
)
