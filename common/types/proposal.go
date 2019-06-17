package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Type that represents Proposal Type as a byte
type ProposalKind byte

//nolint
const (
	ProposalTypeNil             ProposalKind = 0x00
	ProposalTypeText            ProposalKind = 0x01
	ProposalTypeParameterChange ProposalKind = 0x02
	ProposalTypeSoftwareUpgrade ProposalKind = 0x03
	ProposalTypeListTradingPair ProposalKind = 0x04
	// ProposalTypeFeeChange belongs to ProposalTypeParameterChange. We use this to make it easily to distinguish。
	ProposalTypeFeeChange       ProposalKind = 0x05
	ProposalTypeCreateValidator ProposalKind = 0x06
	ProposalTypeRemoveValidator ProposalKind = 0x07
)

// Marshal needed for protobuf compatibility
func (pt ProposalKind) Marshal() ([]byte, error) {
	return []byte{byte(pt)}, nil
}

// Unmarshal needed for protobuf compatibility
func (pt *ProposalKind) Unmarshal(data []byte) error {
	*pt = ProposalKind(data[0])
	return nil
}

// Marshals to JSON using string
func (pt ProposalKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(pt.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (pt *ProposalKind) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := ProposalTypeFromString(s)
	if err != nil {
		return err
	}
	*pt = bz2
	return nil
}

// String to proposalType byte.  Returns ff if invalid.
func ProposalTypeFromString(str string) (ProposalKind, error) {
	switch str {
	case "Text":
		return ProposalTypeText, nil
	case "ParameterChange":
		return ProposalTypeParameterChange, nil
	case "SoftwareUpgrade":
		return ProposalTypeSoftwareUpgrade, nil
	case "ListTradingPair":
		return ProposalTypeListTradingPair, nil
	case "FeeChange":
		return ProposalTypeFeeChange, nil
	case "CreateValidator":
		return ProposalTypeCreateValidator, nil
	case "RemoveValidator":
		return ProposalTypeRemoveValidator, nil
	default:
		return ProposalKind(0xff), errors.Errorf("'%s' is not a valid proposal type", str)
	}
}

// Turns VoteOption byte to String
func (pt ProposalKind) String() string {
	switch pt {
	case ProposalTypeText:
		return "Text"
	case ProposalTypeParameterChange:
		return "ParameterChange"
	case ProposalTypeSoftwareUpgrade:
		return "SoftwareUpgrade"
	case ProposalTypeListTradingPair:
		return "ListTradingPair"
	case ProposalTypeFeeChange:
		return "FeeChange"
	case ProposalTypeCreateValidator:
		return "CreateValidator"
	case ProposalTypeRemoveValidator:
		return "RemoveValidator"
	default:
		return ""
	}
}

type ProposalStatus byte

//nolint
const (
	StatusNil           ProposalStatus = 0x00
	StatusDepositPeriod ProposalStatus = 0x01
	StatusVotingPeriod  ProposalStatus = 0x02
	StatusPassed        ProposalStatus = 0x03
	StatusRejected      ProposalStatus = 0x04
)

// ProposalStatusToString turns a string into a ProposalStatus
func ProposalStatusFromString(str string) (ProposalStatus, error) {
	switch str {
	case "DepositPeriod":
		return StatusDepositPeriod, nil
	case "VotingPeriod":
		return StatusVotingPeriod, nil
	case "Passed":
		return StatusPassed, nil
	case "Rejected":
		return StatusRejected, nil
	case "":
		return StatusNil, nil
	default:
		return ProposalStatus(0xff), errors.Errorf("'%s' is not a valid proposal status", str)
	}
}

// is defined ProposalType?
func validProposalStatus(status ProposalStatus) bool {
	if status == StatusDepositPeriod ||
		status == StatusVotingPeriod ||
		status == StatusPassed ||
		status == StatusRejected {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (status ProposalStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *ProposalStatus) Unmarshal(data []byte) error {
	*status = ProposalStatus(data[0])
	return nil
}

// Marshals to JSON using string
func (status ProposalStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (status *ProposalStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := ProposalStatusFromString(s)
	if err != nil {
		return err
	}
	*status = bz2
	return nil
}

// Turns VoteStatus byte to String
func (status ProposalStatus) String() string {
	switch status {
	case StatusDepositPeriod:
		return "DepositPeriod"
	case StatusVotingPeriod:
		return "VotingPeriod"
	case StatusPassed:
		return "Passed"
	case StatusRejected:
		return "Rejected"
	default:
		return ""
	}
}

// For Printf / Sprintf, returns bech32 when using %s
// nolint: errcheck
func (status ProposalStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}

// Tally Results
type TallyResult struct {
	Yes        Dec `json:"yes"`
	Abstain    Dec `json:"abstain"`
	No         Dec `json:"no"`
	NoWithVeto Dec `json:"no_with_veto"`
	Total      Dec `json:"total"`
}

type Proposal interface {
	GetProposalID() int64
	SetProposalID(int64)

	GetTitle() string
	SetTitle(string)

	GetDescription() string
	SetDescription(string)

	GetProposalType() ProposalKind
	SetProposalType(ProposalKind)

	GetStatus() ProposalStatus
	SetStatus(ProposalStatus)

	GetTallyResult() TallyResult
	SetTallyResult(TallyResult)

	GetSubmitTime() time.Time
	SetSubmitTime(time.Time)

	GetTotalDeposit() Coins
	SetTotalDeposit(Coins)

	GetVotingStartTime() time.Time
	SetVotingStartTime(time.Time)

	GetVotingPeriod() time.Duration
	SetVotingPeriod(time.Duration)
}

// Text Proposals
type TextProposal struct {
	ProposalID   int64         `json:"proposal_id"`   //  ID of the proposal
	Title        string        `json:"title"`         //  Title of the proposal
	Description  string        `json:"description"`   //  Description of the proposal
	ProposalType ProposalKind  `json:"proposal_type"` //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	VotingPeriod time.Duration `json:"voting_period"` //  Length of the voting period

	Status      ProposalStatus `json:"proposal_status"` //  Status of the Proposal {Pending, Active, Passed, Rejected}
	TallyResult TallyResult    `json:"tally_result"`    //  Result of Tallys

	SubmitTime   time.Time `json:"submit_time"`   //  Height of the block where TxGovSubmitProposal was included
	TotalDeposit Coins     `json:"total_deposit"` //  Current deposit on this proposal. Initial value is set at InitialDeposit

	VotingStartTime time.Time `json:"voting_start_time"` //  Height of the block where MinDeposit was reached. -1 if MinDeposit is not reached
}

// Implements Proposal Interface
var _ Proposal = (*TextProposal)(nil)

// nolint
func (tp TextProposal) GetProposalID() int64                       { return tp.ProposalID }
func (tp *TextProposal) SetProposalID(proposalID int64)            { tp.ProposalID = proposalID }
func (tp TextProposal) GetTitle() string                           { return tp.Title }
func (tp *TextProposal) SetTitle(title string)                     { tp.Title = title }
func (tp TextProposal) GetDescription() string                     { return tp.Description }
func (tp *TextProposal) SetDescription(description string)         { tp.Description = description }
func (tp TextProposal) GetProposalType() ProposalKind              { return tp.ProposalType }
func (tp *TextProposal) SetProposalType(proposalType ProposalKind) { tp.ProposalType = proposalType }
func (tp TextProposal) GetStatus() ProposalStatus                  { return tp.Status }
func (tp *TextProposal) SetStatus(status ProposalStatus)           { tp.Status = status }
func (tp TextProposal) GetTallyResult() TallyResult                { return tp.TallyResult }
func (tp *TextProposal) SetTallyResult(tallyResult TallyResult)    { tp.TallyResult = tallyResult }
func (tp TextProposal) GetSubmitTime() time.Time                   { return tp.SubmitTime }
func (tp *TextProposal) SetSubmitTime(submitTime time.Time)        { tp.SubmitTime = submitTime }
func (tp TextProposal) GetTotalDeposit() Coins                     { return tp.TotalDeposit }
func (tp *TextProposal) SetTotalDeposit(totalDeposit Coins)        { tp.TotalDeposit = totalDeposit }
func (tp TextProposal) GetVotingStartTime() time.Time              { return tp.VotingStartTime }
func (tp *TextProposal) SetVotingStartTime(votingStartTime time.Time) {
	tp.VotingStartTime = votingStartTime
}
func (tp TextProposal) GetVotingPeriod() time.Duration { return tp.VotingPeriod }
func (tp *TextProposal) SetVotingPeriod(votingPeriod time.Duration) {
	tp.VotingPeriod = votingPeriod
}

type QueryProposalsParams struct {
	ProposalStatus     ProposalStatus
	NumLatestProposals int64
}

// Params for query 'custom/gov/proposal'
type QueryProposalParams struct {
	ProposalID int64
}
