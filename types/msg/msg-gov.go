package msg

import (
	"github.com/cosmos/cosmos-sdk/x/gov"
)

const (
	MsgRoute = gov.MsgRoute

	MaxTitleLength       = gov.MaxTitleLength
	MaxDescriptionLength = gov.MaxDescriptionLength
	MaxVotingPeriod      = gov.MaxVotingPeriod

	OptionEmpty      = gov.OptionEmpty
	OptionYes        = gov.OptionYes
	OptionAbstain    = gov.OptionAbstain
	OptionNo         = gov.OptionNo
	OptionNoWithVeto = gov.OptionNoWithVeto

	ProposalTypeNil             = gov.ProposalTypeNil
	ProposalTypeText            = gov.ProposalTypeText
	ProposalTypeParameterChange = gov.ProposalTypeParameterChange
	ProposalTypeSoftwareUpgrade = gov.ProposalTypeSoftwareUpgrade
	ProposalTypeListTradingPair = gov.ProposalTypeListTradingPair
	ProposalTypeFeeChange       = gov.ProposalTypeFeeChange
)

type (
	VoteOption   = gov.VoteOption
	ProposalKind = gov.ProposalKind

	ListTradingPairParams = gov.ListTradingPairParams

	SubmitProposalMsg = gov.MsgSubmitProposal
	DepositMsg        = gov.MsgDeposit
	VoteMsg           = gov.MsgVote
)

var (
	VoteOptionFromString   = gov.VoteOptionFromString
	ProposalTypeFromString = gov.ProposalTypeFromString

	NewDepositMsg        = gov.NewMsgDeposit
	NewMsgVote           = gov.NewMsgVote
	NewMsgSubmitProposal = gov.NewMsgSubmitProposal
)
