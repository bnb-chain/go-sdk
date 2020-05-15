package types

import (
	"fmt"
	"time"

	"github.com/tendermint/go-amino"
)

// Delegation represents the bond with tokens held by an account.  It is
// owned by one delegator, and is associated with the voting power of one
// pubKey.
type Delegation struct {
	DelegatorAddr AccAddress `json:"delegator_addr"`
	ValidatorAddr ValAddress `json:"validator_addr"`
	Shares        Dec        `json:"shares"`
}

type DelegationValue struct {
	Shares Dec
	Height int64
}

type Redelegation struct {
	DelegatorAddr    AccAddress `json:"delegator_addr"`     // delegator
	ValidatorSrcAddr ValAddress `json:"validator_src_addr"` // validator redelegation source operator addr
	ValidatorDstAddr ValAddress `json:"validator_dst_addr"` // validator redelegation destination operator addr
	CreationHeight   int64      `json:"creation_height"`    // height which the redelegation took place
	MinTime          time.Time  `json:"min_time"`           // unix time for redelegation completion
	InitialBalance   Coin       `json:"initial_balance"`    // initial balance when redelegation started
	Balance          Coin       `json:"balance"`            // current balance
	SharesSrc        Dec        `json:"shares_src"`         // amount of source shares redelegating
	SharesDst        Dec        `json:"shares_dst"`         // amount of destination shares redelegating
}

type redValue struct {
	CreationHeight int64
	MinTime        time.Time
	InitialBalance Coin
	Balance        Coin
	SharesSrc      Dec
	SharesDst      Dec
}

// unmarshal a redelegation from a store key and value
func UnmarshalRED(cdc *amino.Codec, key, value []byte) (red Redelegation, err error) {
	var storeValue redValue
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &storeValue)
	if err != nil {
		return
	}

	addrs := key[1:] // remove prefix bytes
	if len(addrs) != 3*AddrLen {
		err = fmt.Errorf("unexpected address length for this (address, srcValidator, dstValidator) tuple")
		return
	}
	delAddr := AccAddress(addrs[:AddrLen])
	valSrcAddr := ValAddress(addrs[AddrLen : 2*AddrLen])
	valDstAddr := ValAddress(addrs[2*AddrLen:])

	return Redelegation{
		DelegatorAddr:    delAddr,
		ValidatorSrcAddr: valSrcAddr,
		ValidatorDstAddr: valDstAddr,
		CreationHeight:   storeValue.CreationHeight,
		MinTime:          storeValue.MinTime,
		InitialBalance:   storeValue.InitialBalance,
		Balance:          storeValue.Balance,
		SharesSrc:        storeValue.SharesSrc,
		SharesDst:        storeValue.SharesDst,
	}, nil
}

type DelegationResponse struct {
	Delegation
	Balance Coin `json:"balance"`
}

type QueryDelegatorParams struct {
	BaseParams
	DelegatorAddr AccAddress
}