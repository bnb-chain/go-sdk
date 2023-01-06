package msg

import (
	"github.com/bnb-chain/node/plugins/tokens/freeze"
)

type (
	TokenFreezeMsg   = freeze.FreezeMsg
	TokenUnfreezeMsg = freeze.UnfreezeMsg
)

var (
	NewFreezeMsg   = freeze.NewFreezeMsg
	NewUnfreezeMsg = freeze.NewUnfreezeMsg
)
