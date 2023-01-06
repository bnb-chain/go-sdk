package types

import (
	"github.com/bnb-chain/node/plugins/tokens/swap"
)

const (
	NULL      = swap.NULL
	Open      = swap.Open
	Completed = swap.Completed
	Expired   = swap.Expired
)

var (
	NewSwapStatusFromString = swap.NewSwapStatusFromString
)

type (
	SwapStatus                 = swap.SwapStatus
	SwapBytes                  = swap.SwapBytes
	AtomicSwap                 = swap.AtomicSwap
	QuerySwapByID              swap.QuerySwapByID
	QuerySwapByCreatorParams   = swap.QuerySwapByCreatorParams
	QuerySwapByRecipientParams swap.QuerySwapByRecipientParams
)
