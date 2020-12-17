package msg

import (
	"github.com/tendermint/go-amino"
)

var MsgCdc = amino.NewCodec()

func RegisterCodec(cdc *amino.Codec) {

	cdc.RegisterInterface((*Msg)(nil), nil)

	cdc.RegisterConcrete(SideChainSubmitProposalMsg{}, "cosmos-sdk/MsgSideChainSubmitProposal", nil)
	cdc.RegisterConcrete(SideChainDepositMsg{}, "cosmos-sdk/MsgSideChainDeposit", nil)
	cdc.RegisterConcrete(SideChainVoteMsg{}, "cosmos-sdk/MsgSideChainVote", nil)

	cdc.RegisterInterface((*SCParam)(nil), nil)
	cdc.RegisterConcrete(&OracleParams{}, "params/OracleParamSet", nil)
	cdc.RegisterConcrete(&StakeParams{}, "params/StakeParamSet", nil)
	cdc.RegisterConcrete(&SlashParams{}, "params/SlashParamSet", nil)
	cdc.RegisterConcrete(&IbcParams{}, "params/IbcParamSet", nil)

	cdc.RegisterConcrete(CreateOrderMsg{}, "dex/NewOrder", nil)
	cdc.RegisterConcrete(CancelOrderMsg{}, "dex/CancelOrder", nil)
	cdc.RegisterConcrete(TokenIssueMsg{}, "tokens/IssueMsg", nil)
	cdc.RegisterConcrete(TokenBurnMsg{}, "tokens/BurnMsg", nil)
	cdc.RegisterConcrete(TransferOwnershipMsg{}, "tokens/TransferOwnershipMsg", nil)

	cdc.RegisterConcrete(TimeLockMsg{}, "tokens/TimeLockMsg", nil)
	cdc.RegisterConcrete(TokenFreezeMsg{}, "tokens/FreezeMsg", nil)
	cdc.RegisterConcrete(TokenUnfreezeMsg{}, "tokens/UnfreezeMsg", nil)

	cdc.RegisterConcrete(TimeUnlockMsg{}, "tokens/TimeUnlockMsg", nil)
	cdc.RegisterConcrete(TimeRelockMsg{}, "tokens/TimeRelockMsg", nil)

	cdc.RegisterConcrete(HTLTMsg{}, "tokens/HTLTMsg", nil)
	cdc.RegisterConcrete(DepositHTLTMsg{}, "tokens/DepositHTLTMsg", nil)
	cdc.RegisterConcrete(ClaimHTLTMsg{}, "tokens/ClaimHTLTMsg", nil)
	cdc.RegisterConcrete(RefundHTLTMsg{}, "tokens/RefundHTLTMsg", nil)

	cdc.RegisterConcrete(DexListMsg{}, "dex/ListMsg", nil)
	cdc.RegisterConcrete(MintMsg{}, "tokens/MintMsg", nil)
	//Must use cosmos-sdk.
	cdc.RegisterConcrete(SendMsg{}, "cosmos-sdk/Send", nil)

	cdc.RegisterConcrete(SubmitProposalMsg{}, "cosmos-sdk/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(DepositMsg{}, "cosmos-sdk/MsgDeposit", nil)
	cdc.RegisterConcrete(VoteMsg{}, "cosmos-sdk/MsgVote", nil)

	cdc.RegisterConcrete(SetAccountFlagsMsg{}, "scripts/SetAccountFlagsMsg", nil)

	cdc.RegisterConcrete(MsgCreateValidator{}, "cosmos-sdk/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgRemoveValidator{}, "cosmos-sdk/MsgRemoveValidator", nil)
	cdc.RegisterConcrete(MsgCreateValidatorProposal{}, "cosmos-sdk/MsgCreateValidatorProposal", nil)

	cdc.RegisterConcrete(CreateSideChainValidatorMsg{}, "cosmos-sdk/MsgCreateSideChainValidator", nil)
	cdc.RegisterConcrete(EditSideChainValidatorMsg{}, "cosmos-sdk/MsgEditSideChainValidator", nil)
	cdc.RegisterConcrete(SideChainDelegateMsg{}, "cosmos-sdk/MsgSideChainDelegate", nil)
	cdc.RegisterConcrete(SideChainRedelegateMsg{}, "cosmos-sdk/MsgSideChainRedelegate", nil)
	cdc.RegisterConcrete(SideChainUndelegateMsg{}, "cosmos-sdk/MsgSideChainUndelegate", nil)
	cdc.RegisterConcrete(MsgSideChainUnjail{}, "cosmos-sdk/MsgSideChainUnjail", nil)

	cdc.RegisterConcrete(BindMsg{}, "bridge/BindMsg", nil)
	cdc.RegisterConcrete(TransferOutMsg{}, "bridge/TransferOutMsg", nil)
	cdc.RegisterConcrete(Claim{}, "oracle/Claim", nil)
	cdc.RegisterConcrete(Prophecy{}, "oracle/Prophecy", nil)
	cdc.RegisterConcrete(Status{}, "oracle/Status", nil)
	cdc.RegisterConcrete(DBProphecy{}, "oracle/DBProphecy", nil)
	cdc.RegisterConcrete(ClaimMsg{}, "oracle/ClaimMsg", nil)

	cdc.RegisterConcrete(MiniTokenIssueMsg{}, "tokens/IssueMiniMsg", nil)
	cdc.RegisterConcrete(TinyTokenIssueMsg{}, "tokens/IssueTinyMsg", nil)
	cdc.RegisterConcrete(SetURIMsg{}, "tokens/SetURIMsg", nil)
	cdc.RegisterConcrete(ListMiniMsg{}, "dex/ListMiniMsg", nil)
}

func init() {
	RegisterCodec(MsgCdc)
}
