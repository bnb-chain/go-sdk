package msg

import (
	"github.com/cosmos/cosmos-sdk/x/slashing"
)

const (
	TypeMsgUnjail = slashing.TypeMsgUnjail
	SlashMsgRoute = slashing.MsgRoute
)

type MsgUnjail = slashing.MsgUnjail

var NewMsgUnjail = slashing.NewMsgUnjail
