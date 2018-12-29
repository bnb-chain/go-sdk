package tx

import (
	"github.com/tendermint/go-amino"

	"github.com/binance-chain/go-sdk/tx/txmsg"
)

// Cdc global variable
var Cdc = amino.NewCodec()


func RegisterCodec(cdc *amino.Codec) {
	cdc.RegisterInterface((*Tx)(nil), nil)
	cdc.RegisterConcrete(StdTx{}, "auth/StdTx", nil)

	txmsg.RegisterCodec(cdc)
}

func init() {
	RegisterCodec(Cdc)
}
