package types

import (
	"github.com/cosmos/cosmos-sdk/types"
)

type AccAddress = types.AccAddress

type ChainNetwork uint8

const (
	TestNetwork ChainNetwork = iota
	ProdNetwork
	TmpTestNetwork
	GangesNetwork
)

const (
	AddrLen = types.AddrLen
)

var Network = ProdNetwork

func (this ChainNetwork) Bech32Prefixes() string {
	switch this {
	case TestNetwork:
		return "tbnb"
	case TmpTestNetwork:
		return "tbnb"
	case GangesNetwork:
		return "tbnb"
	case ProdNetwork:
		return "bnb"
	default:
		panic("Unknown network type")
	}
}

func (this ChainNetwork) Bech32ValidatorAddrPrefix() string {
	return "bva"
}

var (
	AccAddressFromHex    = types.AccAddressFromHex
	AccAddressFromBech32 = types.AccAddressFromBech32
	GetFromBech32        = types.GetFromBech32
	MustBech32ifyConsPub = types.MustBech32ifyConsPub
	Bech32ifyConsPub     = types.Bech32ifyConsPub
	GetConsPubKeyBech32  = types.GetConsPubKeyBech32
)
