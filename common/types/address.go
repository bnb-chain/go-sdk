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

func SetNetwork(network ChainNetwork) {
	Network = network
	if network != ProdNetwork {
		sdkConfig := types.GetConfig()
		sdkConfig.SetBech32PrefixForAccount("tbnb", "bnbp")
	}
}

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

func init() {
	sdkConfig := types.GetConfig()
	sdkConfig.SetBech32PrefixForAccount("bnb", "bnbp")
	sdkConfig.SetBech32PrefixForValidator("bva", "bvap")
	sdkConfig.SetBech32PrefixForConsensusNode("bca", "bcap")
}

var (
	AccAddressFromHex    = types.AccAddressFromHex
	AccAddressFromBech32 = types.AccAddressFromBech32
	GetFromBech32        = types.GetFromBech32
	MustBech32ifyConsPub = types.MustBech32ifyConsPub
	Bech32ifyConsPub     = types.Bech32ifyConsPub
	GetConsPubKeyBech32  = types.GetConsPubKeyBech32
)
