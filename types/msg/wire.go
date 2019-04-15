package msg

import (
	"github.com/tendermint/go-amino"
)

var MsgCdc = amino.NewCodec()

func RegisterCodec(cdc *amino.Codec) {

	cdc.RegisterInterface((*Msg)(nil), nil)

	cdc.RegisterConcrete(CreateOrderMsg{}, "dex/NewOrder", nil)
	cdc.RegisterConcrete(CancelOrderMsg{}, "dex/CancelOrder", nil)
	cdc.RegisterConcrete(TokenIssueMsg{}, "tokens/IssueMsg", nil)
	cdc.RegisterConcrete(TokenBurnMsg{}, "tokens/BurnMsg", nil)
	cdc.RegisterConcrete(TokenFreezeMsg{}, "tokens/FreezeMsg", nil)
	cdc.RegisterConcrete(TokenUnfreezeMsg{}, "tokens/UnfreezeMsg", nil)
	cdc.RegisterConcrete(DexListMsg{}, "dex/ListMsg", nil)
	cdc.RegisterConcrete(MintMsg{}, "tokens/MintMsg", nil)
	//Must use cosmos-sdk.
	cdc.RegisterConcrete(SendMsg{}, "cosmos-sdk/Send", nil)

	cdc.RegisterConcrete(SubmitProposalMsg{}, "cosmos-sdk/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(DepositMsg{}, "cosmos-sdk/MsgDeposit", nil)
	cdc.RegisterConcrete(VoteMsg{}, "cosmos-sdk/MsgVote", nil)
}

func init() {
	RegisterCodec(MsgCdc)
}
