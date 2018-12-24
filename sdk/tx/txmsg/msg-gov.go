package txmsg

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tendermint/go-amino"
	"time"
)

// name to idetify transaction types
const MsgRoute = "gov"

type VoteOption byte

//nolint
const (
	OptionEmpty      VoteOption = 0x00
	OptionYes        VoteOption = 0x01
	OptionAbstain    VoteOption = 0x02
	OptionNo         VoteOption = 0x03
	OptionNoWithVeto VoteOption = 0x04
)

// String to proposalType byte.  Returns ff if invalid.
func VoteOptionFromString(str string) (VoteOption, error) {
	switch str {
	case "Yes":
		return OptionYes, nil
	case "Abstain":
		return OptionAbstain, nil
	case "No":
		return OptionNo, nil
	case "NoWithVeto":
		return OptionNoWithVeto, nil
	default:
		return VoteOption(0xff), errors.Errorf("'%s' is not a valid vote option", str)
	}
}

// Marshal needed for protobuf compatibility
func (vo VoteOption) Marshal() ([]byte, error) {
	return []byte{byte(vo)}, nil
}

// Unmarshal needed for protobuf compatibility
func (vo *VoteOption) Unmarshal(data []byte) error {
	*vo = VoteOption(data[0])
	return nil
}

// Marshals to JSON using string
func (vo VoteOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(vo.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (vo *VoteOption) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := VoteOptionFromString(s)
	if err != nil {
		return err
	}
	*vo = bz2
	return nil
}

// Turns VoteOption byte to String
func (vo VoteOption) String() string {
	switch vo {
	case OptionYes:
		return "Yes"
	case OptionAbstain:
		return "Abstain"
	case OptionNo:
		return "No"
	case OptionNoWithVeto:
		return "NoWithVeto"
	default:
		return ""
	}
}

// For Printf / Sprintf, returns bech32 when using %s
// nolint: errcheck
func (vo VoteOption) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", vo.String())))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(vo))))
	}
}

//-----------------------------------------------------------
// ProposalKind

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
	ProposalTypeFeeChange ProposalKind = 0x05
)

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
	default:
		return ProposalKind(0xff), errors.Errorf("'%s' is not a valid proposal type", str)
	}
}

// is defined ProposalType?
func validProposalType(pt ProposalKind) bool {
	if pt == ProposalTypeText ||
		pt == ProposalTypeParameterChange ||
		pt == ProposalTypeSoftwareUpgrade ||
		pt == ProposalTypeListTradingPair ||
		pt == ProposalTypeFeeChange {
		return true
	}
	return false
}

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
	default:
		return ""
	}
}

// For Printf / Sprintf, returns bech32 when using %s
// nolint: errcheck
func (pt ProposalKind) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(pt.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(pt))))
	}
}

//-----------------------------------------------------------
type ListTradingPairParams struct {
	BaseAssetSymbol  string    `json:"base_asset_symbol"`  // base asset symbol
	QuoteAssetSymbol string    `json:"quote_asset_symbol"` // quote asset symbol
	InitPrice        int64     `json:"init_price"`         // init price
	Description      string    `json:"description"`        // description
	ExpireTime       time.Time `json:"expire_time"`        // expire time
}

//-----------------------------------------------------------
type SubmitProposalMsg struct {
	Title          string       `json:"title"`           //  Title of the proposal
	Description    string       `json:"description"`     //  Description of the proposal
	ProposalType   ProposalKind `json:"proposal_type"`   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       AccAddress   `json:"proposer"`        //  Address of the proposer
	InitialDeposit Coins        `json:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive.
}

func NewMsgSubmitProposal(title string, description string, proposalType ProposalKind, proposer AccAddress, initialDeposit Coins) SubmitProposalMsg {
	return SubmitProposalMsg{
		Title:          title,
		Description:    description,
		ProposalType:   proposalType,
		Proposer:       proposer,
		InitialDeposit: initialDeposit,
	}
}

//nolint
func (msg SubmitProposalMsg) Route() string { return MsgRoute }
func (msg SubmitProposalMsg) Type() string  { return "submit_proposal" }

// Implements Msg.
func (msg SubmitProposalMsg) ValidateBasic() error {
	if len(msg.Title) == 0 {
		return fmt.Errorf("title can't be empty")
	}
	if len(msg.Description) == 0 {
		return fmt.Errorf("description can't be empty")
	}
	if !validProposalType(msg.ProposalType) {
		return fmt.Errorf("invalid proposal type %v ", msg.ProposalType)
	}
	if len(msg.Proposer) == 0 {
		return fmt.Errorf("proposer can't be empty")
	}
	if !msg.InitialDeposit.IsValid() {
		return fmt.Errorf("initial deposit %v is invalid. ", msg.InitialDeposit)
	}
	if !msg.InitialDeposit.IsNotNegative() {
		return fmt.Errorf("initial deposit %v is negative. ", msg.InitialDeposit)
	}
	return nil
}

func (msg SubmitProposalMsg) String() string {
	return fmt.Sprintf("SubmitProposalMsg{%s, %s, %s, %v}", msg.Title, msg.Description, msg.ProposalType, msg.InitialDeposit)
}

// Implements Msg.
func (msg SubmitProposalMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg SubmitProposalMsg) GetSignBytes() []byte {
	b, err := amino.NewCodec().MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return MustSortJSON(b)
}

// Implements Msg.
func (msg SubmitProposalMsg) GetSigners() []AccAddress {
	return []AccAddress{msg.Proposer}
}

func (msg SubmitProposalMsg) GetInvolvedAddresses() []AccAddress {
	return msg.GetSigners()
}

//-----------------------------------------------------------
// DepositMsg
type DepositMsg struct {
	ProposalID int64      `json:"proposal_id"` // ID of the proposal
	Depositer  AccAddress `json:"depositer"`   // Address of the depositer
	Amount     Coins      `json:"amount"`      // Coins to add to the proposal's deposit
}

func NewDepositMsg(depositer AccAddress, proposalID int64, amount Coins) DepositMsg {
	return DepositMsg{
		ProposalID: proposalID,
		Depositer:  depositer,
		Amount:     amount,
	}
}

// Implements Msg.
// nolint
func (msg DepositMsg) Route() string { return MsgRoute }
func (msg DepositMsg) Type() string  { return "deposit" }

// Implements Msg.
func (msg DepositMsg) ValidateBasic() error {
	if len(msg.Depositer) == 0 {
		return fmt.Errorf("depositer can't be empty ")
	}
	if !msg.Amount.IsValid() {
		return fmt.Errorf("amount is invalid ")
	}
	if !msg.Amount.IsNotNegative() {
		return fmt.Errorf("amount can't be negative ")
	}
	if msg.ProposalID < 0 {
		return fmt.Errorf("proposalId can't be negative ")
	}
	return nil
}

func (msg DepositMsg) String() string {
	return fmt.Sprintf("DepositMsg{%s=>%v: %v}", msg.Depositer, msg.ProposalID, msg.Amount)
}

// Implements Msg.
func (msg DepositMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg DepositMsg) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return MustSortJSON(b)
}

// Implements Msg.
func (msg DepositMsg) GetSigners() []AccAddress {
	return []AccAddress{msg.Depositer}
}

func (msg DepositMsg) GetInvolvedAddresses() []AccAddress {
	return msg.GetSigners()
}

//-----------------------------------------------------------
// VoteMsg
type VoteMsg struct {
	ProposalID int64      `json:"proposal_id"` // ID of the proposal
	Voter      AccAddress `json:"voter"`       //  address of the voter
	Option     VoteOption `json:"option"`      //  option from OptionSet chosen by the voter
}

func NewMsgVote(voter AccAddress, proposalID int64, option VoteOption) VoteMsg {
	return VoteMsg{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     option,
	}
}

// Implements Msg.
// nolint
func (msg VoteMsg) Route() string { return MsgRoute }
func (msg VoteMsg) Type() string  { return "vote" }

// Implements Msg.
func (msg VoteMsg) ValidateBasic() error {
	if len(msg.Voter.Bytes()) == 0 {
		return fmt.Errorf("vaoter can't be empty ")
	}
	if msg.ProposalID < 0 {
		return fmt.Errorf("proposalId can't be less than 0")
	}
	if !validVoteOption(msg.Option) {
		return fmt.Errorf("invalid msg option %v", msg.Option)
	}
	return nil
}

func (msg VoteMsg) String() string {
	return fmt.Sprintf("VoteMsg{%v - %s}", msg.ProposalID, msg.Option)
}

// Implements Msg.
func (msg VoteMsg) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg VoteMsg) GetSignBytes() []byte {
	b, err := amino.NewCodec().MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return MustSortJSON(b)
}

// Implements Msg.
func (msg VoteMsg) GetSigners() []AccAddress {
	return []AccAddress{msg.Voter}
}

func (msg VoteMsg) GetInvolvedAddresses() []AccAddress {
	return msg.GetSigners()
}

func validVoteOption(option VoteOption) bool {
	if option == OptionYes ||
		option == OptionAbstain ||
		option == OptionNo ||
		option == OptionNoWithVeto {
		return true
	}
	return false
}
