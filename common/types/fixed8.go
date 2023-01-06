package types

import (
	"github.com/bnb-chain/node/common/utils"
)

var (
	Fixed8Decimals = utils.Fixed8Decimals
	Fixed8One      = utils.Fixed8One
	Fixed8Zero     = utils.NewFixed8(0)
)

type Fixed8 = utils.Fixed8

var (
	NewFixed8          = utils.NewFixed8
	Fixed8DecodeString = utils.Fixed8DecodeString
)
