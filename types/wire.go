package types

import (
	ntypes "github.com/cbarraford/go-sdk/common/types"
	"github.com/cbarraford/go-sdk/types/tx"
	amino "github.com/tendermint/go-amino"
	types "github.com/tendermint/tendermint/rpc/core/types"
)

func NewCodec() *amino.Codec {
	cdc := amino.NewCodec()
	types.RegisterAmino(cdc)
	ntypes.RegisterWire(cdc)
	tx.RegisterCodec(cdc)
	return cdc
}
