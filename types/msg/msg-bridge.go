package msg

import (
	"github.com/bnb-chain/node/plugins/bridge/types"
	ctypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	RouteBridge = types.RouteBridge

	BindMsgType        = types.BindMsgType
	UnbindMsgType      = types.UnbindMsgType
	TransferOutMsgType = types.TransferOutMsgType

	MaxSymbolLength = types.MaxSymbolLength
)

type (
	SmartChainAddress = ctypes.SmartChainAddress
	BindStatus        = types.BindStatus
	BindMsg           = types.BindMsg
	TransferOutMsg    = types.TransferOutMsg
	UnbindMsg         = types.UnbindMsg
)

var (
	NewSmartChainAddress = ctypes.NewSmartChainAddress
	NewBindMsg           = types.NewBindMsg
	NewTransferOutMsg    = types.NewTransferOutMsg
	NewUnbindMsg         = types.NewUnbindMsg
)
