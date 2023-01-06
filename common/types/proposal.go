package types

import (
	"github.com/cosmos/cosmos-sdk/x/gov"
)

type ProposalKind = gov.ProposalKind

// nolint
const (
	ProposalTypeNil                  = gov.ProposalTypeNil
	ProposalTypeText                 = gov.ProposalTypeText
	ProposalTypeParameterChange      = gov.ProposalTypeParameterChange
	ProposalTypeSoftwareUpgrade      = gov.ProposalTypeSoftwareUpgrade
	ProposalTypeListTradingPair      = gov.ProposalTypeListTradingPair
	ProposalTypeFeeChange            = gov.ProposalTypeFeeChange
	ProposalTypeCreateValidator      = gov.ProposalTypeCreateValidator
	ProposalTypeRemoveValidator      = gov.ProposalTypeRemoveValidator
	ProposalTypeDelistTradingPair    = gov.ProposalTypeDelistTradingPair
	ProposalTypeManageChanPermission = gov.ProposalTypeManageChanPermission

	ProposalTypeSCParamsChange  = gov.ProposalTypeSCParamsChange
	ProposalTypeCSCParamsChange = gov.ProposalTypeCSCParamsChange
)

var (
	ProposalTypeFromString = gov.ProposalTypeFromString
)

type ProposalStatus = gov.ProposalStatus

// nolint
const (
	StatusNil           = gov.StatusNil
	StatusDepositPeriod = gov.StatusDepositPeriod
	StatusVotingPeriod  = gov.StatusVotingPeriod
	StatusPassed        = gov.StatusPassed
	StatusRejected      = gov.StatusRejected
	StatusExecuted      = gov.StatusExecuted
)

var (
	ProposalStatusFromString = gov.ProposalStatusFromString
)

type (
	TallyResult          = gov.TallyResult
	Proposal             = gov.Proposal
	TextProposal         = gov.TextProposal
	BaseParams           = gov.BaseParams
	QueryProposalsParams = gov.QueryProposalsParams
	QueryProposalParams  = gov.QueryProposalParams
)
