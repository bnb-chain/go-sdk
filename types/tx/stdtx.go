package tx

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	context "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
)

const Source int64 = 0

type (
	Tx           = types.Tx
	StdTx        = auth.StdTx
	StdSignDoc   = auth.StdSignDoc
	StdSignature = auth.StdSignature
	StdSignMsg   = context.StdSignMsg
)

var (
	StdSignBytes = auth.StdSignBytes
	NewStdTx     = auth.NewStdTx
)
