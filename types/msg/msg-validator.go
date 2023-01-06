package msg

import (
	stakeTypes "github.com/cosmos/cosmos-sdk/x/stake/types"
)

type (
	Description            = stakeTypes.Description
	MsgCreateValidatorOpen = stakeTypes.MsgCreateValidatorOpen
	MsgRemoveValidator     = stakeTypes.MsgRemoveValidator
	MsgEditValidator       = stakeTypes.MsgEditValidator
	MsgDelegate            = stakeTypes.MsgDelegate
	MsgRedelegate          = stakeTypes.MsgRedelegate
	MsgUndelegate          = stakeTypes.MsgUndelegate
)

var NewMsgRemoveValidator = stakeTypes.NewMsgRemoveValidator
