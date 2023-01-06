package msg

import (
	stakeTypes "github.com/cosmos/cosmos-sdk/x/stake/types"
)

const (
	TypeCreateSideChainValidator = stakeTypes.MsgTypeCreateSideChainValidator
	TypeEditSideChainValidator   = stakeTypes.MsgTypeEditSideChainValidator
	TypeSideChainDelegate        = stakeTypes.MsgTypeSideChainDelegate
	TypeSideChainRedelegate      = stakeTypes.MsgTypeSideChainRedelegate
	TypeSideChainUndelegate      = stakeTypes.MsgTypeSideChainUndelegate

	SideChainStakeMsgRoute = stakeTypes.MsgRoute
	SideChainAddrLen       = 20

	MinDelegationAmount = 1e7
)

type (
	CreateSideChainValidatorMsg = stakeTypes.MsgCreateSideChainValidator
	EditSideChainValidatorMsg   = stakeTypes.MsgEditSideChainValidator
	SideChainDelegateMsg        = stakeTypes.MsgSideChainDelegate
	SideChainRedelegateMsg      = stakeTypes.MsgSideChainRedelegate
	SideChainUndelegateMsg      = stakeTypes.MsgSideChainUndelegate
)

var (
	NewCreateSideChainValidatorMsg           = stakeTypes.NewMsgCreateSideChainValidator
	NewMsgCreateSideChainValidatorOnBehalfOf = stakeTypes.NewMsgCreateSideChainValidatorOnBehalfOf
	NewEditSideChainValidatorMsg             = stakeTypes.NewMsgEditSideChainValidator
	NewSideChainDelegateMsg                  = stakeTypes.NewMsgSideChainDelegate
	NewSideChainRedelegateMsg                = stakeTypes.NewMsgSideChainRedelegate
	NewSideChainUndelegateMsg                = stakeTypes.NewMsgSideChainUndelegate
)
