package types

import (
	nodeTypes "github.com/bnb-chain/node/common/types"
	"github.com/bnb-chain/node/plugins/tokens/client/rest"
)

type (
	Token        = nodeTypes.Token
	MiniToken    = nodeTypes.MiniToken
	TokenBalance = rest.TokenBalance
)
