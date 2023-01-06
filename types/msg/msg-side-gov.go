package msg

import (
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	oracleTypes "github.com/cosmos/cosmos-sdk/x/oracle/types"
	paramHubTypes "github.com/cosmos/cosmos-sdk/x/paramHub/types"
	sidechainTypes "github.com/cosmos/cosmos-sdk/x/sidechain/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/stake"
)

const (
	MsgTypeSideSubmitProposal = gov.MsgTypeSideSubmitProposal
	MsgTypeSideDeposit        = gov.MsgTypeSideDeposit
	MsgTypeSideVote           = gov.MsgTypeSideVote

	ProposalTypeSCParamsChange  = gov.ProposalTypeSCParamsChange
	ProposalTypeCSCParamsChange = gov.ProposalTypeCSCParamsChange

	MaxSideChainIdLength = sidechainTypes.MaxSideChainIdLength
)

type (
	SideChainSubmitProposalMsg = gov.MsgSideChainSubmitProposal
	SideChainDepositMsg        = gov.MsgSideChainDeposit
	SideChainVoteMsg           = gov.MsgSideChainVote

	SCParam        = paramHubTypes.SCParam
	SCChangeParams = paramHubTypes.SCChangeParams
	IbcParams      = ibc.Params
	OracleParams   = oracleTypes.Params
	SlashParams    = slashing.Params
	StakeParams    = stake.Params
	CSCParamChange = paramHubTypes.CSCParamChange
)

var (
	NewSideChainSubmitProposalMsg = gov.NewMsgSideChainSubmitProposal
	NewSideChainDepositMsg        = gov.NewMsgSideChainDeposit
	NewSideChainVoteMsg           = gov.NewMsgSideChainVote
)
