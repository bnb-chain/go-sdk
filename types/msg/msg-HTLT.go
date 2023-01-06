package msg

import (
	"github.com/bnb-chain/node/plugins/tokens/swap"
)

const (
	AtomicSwapRoute = swap.AtomicSwapRoute
	HTLT            = swap.HTLT
	DepositHTLT     = swap.DepositHTLT
	ClaimHTLT       = swap.ClaimHTLT
	RefundHTLT      = swap.RefundHTLT

	Int64Size               = swap.Int64Size
	RandomNumberHashLength  = swap.RandomNumberHashLength
	RandomNumberLength      = swap.RandomNumberLength
	MaxOtherChainAddrLength = swap.MaxOtherChainAddrLength
	SwapIDLength            = swap.SwapIDLength
	MaxExpectedIncomeLength = swap.MaxExpectedIncomeLength
	MinimumHeightSpan       = swap.MinimumHeightSpan
	MaximumHeightSpan       = swap.MaximumHeightSpan
)

var (
	// bnb prefix address:  bnb1wxeplyw7x8aahy93w96yhwm7xcq3ke4f8ge93u
	// tbnb prefix address: tbnb1wxeplyw7x8aahy93w96yhwm7xcq3ke4ffasp3d
	AtomicSwapCoinsAccAddr = swap.AtomicSwapCoinsAccAddr
)

type (
	HTLTMsg        = swap.HTLTMsg
	DepositHTLTMsg = swap.DepositHTLTMsg
	ClaimHTLTMsg   = swap.ClaimHTLTMsg
	RefundHTLTMsg  = swap.RefundHTLTMsg
)

var (
	NewHTLTMsg        = swap.NewHTLTMsg
	NewDepositHTLTMsg = swap.NewDepositHTLTMsg
	NewClaimHTLTMsg   = swap.NewClaimHTLTMsg
	NewRefundHTLTMsg  = swap.NewRefundHTLTMsg
)
