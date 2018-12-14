package tx

import (
	"github.com/BiJie/bnc-go-sdk/sdk/tx/txmsg"
	"github.com/tendermint/go-amino"
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
