package msg

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/bnb-chain/node/plugins/tokens/swap"
)

var (
	SortJSON            = cTypes.SortJSON
	MustSortJSON        = cTypes.MustSortJSON
	CalculateRandomHash = swap.CalculateRandomHash
	CalculateSwapID     = swap.CalculateSwapID
	HexAddress          = cTypes.HexAddress
	HexEncode           = cTypes.HexEncode
	HexDecode           = cTypes.HexDecode
	Has0xPrefix         = cTypes.Has0xPrefix
)
