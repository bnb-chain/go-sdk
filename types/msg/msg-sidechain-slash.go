package msg

import (
	"github.com/cosmos/cosmos-sdk/x/slashing"
)

const (
	TypeMsgSideChainUnjail = slashing.TypeMsgSideChainUnjail

	SideChainSlashMsgRoute = slashing.MsgRoute
)

type MsgSideChainUnjail = slashing.MsgSideChainUnjail

var NewMsgSideChainUnjail = slashing.NewMsgSideChainUnjail
