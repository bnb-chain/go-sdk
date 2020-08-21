package msg

import (
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"github.com/tendermint/go-amino"

	"github.com/binance-chain/go-sdk/common/types"
)

const (
	MsgTypeSideSubmitProposal = "side_submit_proposal"
	MsgTypeSideDeposit        = "side_deposit"
	MsgTypeSideVote           = "side_vote"

	// side chain params change
	ProposalTypeSCParamsChange ProposalKind = 0x81
	// cross side chain param change
	ProposalTypeCSCParamsChange ProposalKind = 0x82

	MaxSideChainIdLength = 20
)

//-----------------------------------------------------------
// SideChainSubmitProposalMsg
type SideChainSubmitProposalMsg struct {
	Title          string           `json:"title"`           //  Title of the proposal
	Description    string           `json:"description"`     //  Description of the proposal
	ProposalType   ProposalKind     `json:"proposal_type"`   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       types.AccAddress `json:"proposer"`        //  Address of the proposer
	InitialDeposit types.Coins      `json:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive.
	VotingPeriod   time.Duration    `json:"voting_period"`   //  Length of the voting period (s)
	SideChainId    string           `json:"side_chain_id"`
}

func NewSideChainSubmitProposalMsg(title string, description string, proposalType ProposalKind, proposer types.AccAddress, initialDeposit types.Coins, votingPeriod time.Duration, sideChainId string) SideChainSubmitProposalMsg {
	return SideChainSubmitProposalMsg{
		Title:          title,
		Description:    description,
		ProposalType:   proposalType,
		Proposer:       proposer,
		InitialDeposit: initialDeposit,
		VotingPeriod:   votingPeriod,
		SideChainId:    sideChainId,
	}
}

//nolint
func (msg SideChainSubmitProposalMsg) Route() string { return MsgRoute }
func (msg SideChainSubmitProposalMsg) Type() string  { return MsgTypeSideSubmitProposal }

// Implements Msg.
func (msg SideChainSubmitProposalMsg) ValidateBasic() error {
	if len(msg.Title) == 0 {
		return fmt.Errorf("title can't be empty")
	}
	if len(msg.Title) > MaxTitleLength {
		return fmt.Errorf("Proposal title is longer than max length of %d", MaxTitleLength)
	}
	if len(msg.Description) == 0 {
		return fmt.Errorf("description can't be empty")
	}

	if len(msg.Description) > MaxDescriptionLength {
		return fmt.Errorf("Proposal description is longer than max length of %d", MaxDescriptionLength)
	}
	if !validSideProposalType(msg.ProposalType) {
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
	if msg.VotingPeriod <= 0 || msg.VotingPeriod > MaxVotingPeriod {
		return fmt.Errorf("voting period should between 0 and %d weeks", MaxVotingPeriod/(7*24*60*60*time.Second))
	}
	return nil
	return nil
}

func (msg SideChainSubmitProposalMsg) String() string {
	return fmt.Sprintf("SideChainSubmitProposalMsg{%s, %s, %s, %v, %s, %s}", msg.Title,
		msg.Description, msg.ProposalType, msg.InitialDeposit, msg.VotingPeriod, msg.SideChainId)
}

// Implements Msg.
func (msg SideChainSubmitProposalMsg) GetSignBytes() []byte {
	b, err := amino.NewCodec().MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return MustSortJSON(b)
}

// Implements Msg. Identical to MsgSubmitProposal, keep here for code readability.
func (msg SideChainSubmitProposalMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.Proposer}
}

// Implements Msg. Identical to MsgSubmitProposal, keep here for code readability.
func (msg SideChainSubmitProposalMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}

//-----------------------------------------------------------
// SideChainDepositMsg
type SideChainDepositMsg struct {
	ProposalID  int64            `json:"proposal_id"` // ID of the proposal
	Depositer   types.AccAddress `json:"depositer"`   // Address of the depositer
	Amount      types.Coins      `json:"amount"`      // Coins to add to the proposal's deposit
	SideChainId string           `json:"side_chain_id"`
}

func NewSideChainDepositMsg(depositer types.AccAddress, proposalID int64, amount types.Coins, sideChainId string) SideChainDepositMsg {
	return SideChainDepositMsg{
		ProposalID:  proposalID,
		Depositer:   depositer,
		Amount:      amount,
		SideChainId: sideChainId,
	}
}

// nolint
func (msg SideChainDepositMsg) Route() string { return MsgRoute }
func (msg SideChainDepositMsg) Type() string  { return MsgTypeSideDeposit }

// Implements Msg.
func (msg SideChainDepositMsg) ValidateBasic() error {
	if len(msg.SideChainId) == 0 || len(msg.SideChainId) > MaxSideChainIdLength {
		return fmt.Errorf("invalid side chain id %s", msg.SideChainId)
	}
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

func (msg SideChainDepositMsg) String() string {
	return fmt.Sprintf("SideChainDepositMsg{%s=>%v: %v, %s}", msg.Depositer, msg.ProposalID, msg.Amount, msg.SideChainId)
}

// Implements Msg.
func (msg SideChainDepositMsg) GetSignBytes() []byte {
	b, err := amino.NewCodec().MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return MustSortJSON(b)
}

// Implements Msg. Identical to MsgDeposit, keep here for code readability.
func (msg SideChainDepositMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.Depositer}
}

// Implements Msg. Identical to MsgDeposit, keep here for code readability.
func (msg SideChainDepositMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}

//-----------------------------------------------------------
// SideChainVoteMsg

type SideChainVoteMsg struct {
	ProposalID  int64            `json:"proposal_id"` // ID of the proposal
	Voter       types.AccAddress `json:"voter"`       //  address of the voter
	Option      VoteOption       `json:"option"`      //  option from OptionSet chosen by the voter
	SideChainId string           `json:"side_chain_id"`
}

func NewSideChainVoteMsg(voter types.AccAddress, proposalID int64, option VoteOption, sideChainId string) SideChainVoteMsg {
	return SideChainVoteMsg{
		ProposalID:  proposalID,
		Voter:       voter,
		Option:      option,
		SideChainId: sideChainId,
	}
}

func (msg SideChainVoteMsg) Route() string { return MsgRoute }
func (msg SideChainVoteMsg) Type() string  { return MsgTypeSideVote }

// Implements Msg.
func (msg SideChainVoteMsg) ValidateBasic() error {
	if len(msg.SideChainId) == 0 || len(msg.SideChainId) > MaxSideChainIdLength {
		return fmt.Errorf("invalid side chain id %s", msg.SideChainId)
	}
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

func (msg SideChainVoteMsg) String() string {
	return fmt.Sprintf("SideChainVoteMsg{%v - %s, %s}", msg.ProposalID, msg.Option, msg.SideChainId)
}

// Implements Msg.
func (msg SideChainVoteMsg) GetSignBytes() []byte {
	b, err := amino.NewCodec().MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return MustSortJSON(b)
}

// Implements Msg. Identical to MsgVote, keep here for code readability.
func (msg SideChainVoteMsg) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.Voter}
}

// Implements Msg. Identical to MsgVote, keep here for code readability.
func (msg SideChainVoteMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}

func validSideProposalType(pt ProposalKind) bool {
	if pt == ProposalTypeSCParamsChange ||
		pt == ProposalTypeCSCParamsChange {
		return true
	}
	return false
}

// ---------   Definition side chain prams change ------------------- //
type SCParam interface {
	UpdateCheck() error
	// native means weather the parameter stored in native store context or side chain store context
	//GetParamAttribute() (string, bool)
	GetParamAttribute() (string, bool)
}

type SCChangeParams struct {
	SCParams    []SCParam `json:"sc_params"`
	Description string    `json:"description"`
}

func (s *SCChangeParams) Check() error {
	// use literal string to avoid  import cycle
	supportParams := []string{"slash", "ibc", "oracle", "staking"}

	if len(s.SCParams) != len(supportParams) {
		return fmt.Errorf("the sc_params length mismatch, suppose %d", len(supportParams))
	}

	paramSet := make(map[string]bool)
	for _, s := range supportParams {
		paramSet[s] = true
	}

	for _, sc := range s.SCParams {
		if sc == nil {
			return fmt.Errorf("sc_params contains empty element")
		}
		err := sc.UpdateCheck()
		if err != nil {
			return err
		}
		paramType, _ := sc.GetParamAttribute()
		if exist := paramSet[paramType]; exist {
			delete(paramSet, paramType)
		} else {
			return fmt.Errorf("unsupported param type %s", paramType)
		}
	}
	return nil
}

type IbcParams struct {
	RelayerFee int64 `json:"relayer_fee"`
}

func (p *IbcParams) UpdateCheck() error {
	if p.RelayerFee <= 0 {
		return fmt.Errorf("the syn_package_fee should be greater than 0")
	}
	return nil
}

func (p *IbcParams) GetParamAttribute() (string, bool) {
	return "ibc", false
}

type OracleParams struct {
	ConsensusNeeded types.Dec `json:"ConsensusNeeded"` //  Minimum deposit for a proposal to enter voting period.
}

func (p *OracleParams) UpdateCheck() error {
	if p.ConsensusNeeded.IsNil() || p.ConsensusNeeded.GT(types.OneDec()) || p.ConsensusNeeded.LT(types.NewDecWithPrec(5, 1)) {
		return fmt.Errorf("the value should be in range 0.5 to 1")
	}
	return nil
}

func (p *OracleParams) GetParamAttribute() (string, bool) {
	return "oracle", true
}

type SlashParams struct {
	MaxEvidenceAge           time.Duration `json:"max_evidence_age"`
	SignedBlocksWindow       int64         `json:"signed_blocks_window"`
	MinSignedPerWindow       types.Dec     `json:"min_signed_per_window"`
	DoubleSignUnbondDuration time.Duration `json:"double_sign_unbond_duration"`
	DowntimeUnbondDuration   time.Duration `json:"downtime_unbond_duration"`
	TooLowDelUnbondDuration  time.Duration `json:"too_low_del_unbond_duration"`
	SlashFractionDoubleSign  types.Dec     `json:"slash_fraction_double_sign"`
	SlashFractionDowntime    types.Dec     `json:"slash_fraction_downtime"`
	DoubleSignSlashAmount    int64         `json:"double_sign_slash_amount"`
	DowntimeSlashAmount      int64         `json:"downtime_slash_amount"`
	SubmitterReward          int64         `json:"submitter_reward"`
	DowntimeSlashFee         int64         `json:"downtime_slash_fee"`
}

func (p *SlashParams) GetParamAttribute() (string, bool) {
	return "slash", false
}

func (p *SlashParams) UpdateCheck() error {
	// no check for SignedBlocksWindow, MinSignedPerWindow, SlashFractionDoubleSign, SlashFractionDowntime
	if p.MaxEvidenceAge < 1*time.Minute || p.MaxEvidenceAge > 100*24*time.Hour {
		return fmt.Errorf("the max_evidence_age should be in range 1 minutes to 100 day")
	}
	if p.DoubleSignUnbondDuration < 1*time.Hour {
		return fmt.Errorf("the double_sign_unbond_duration should be greate than 1 hour")
	}
	if p.DowntimeUnbondDuration < 60*time.Second || p.DowntimeUnbondDuration > 100*24*time.Hour {
		return fmt.Errorf("the downtime_unbond_duration should be in range 1 minutes to 100 day")
	}
	if p.TooLowDelUnbondDuration < 60*time.Second || p.TooLowDelUnbondDuration > 100*24*time.Hour {
		return fmt.Errorf("the too_low_del_unbond_duration should be in range 1 minutes to 100 day")
	}
	if p.DoubleSignSlashAmount < 1e8 {
		return fmt.Errorf("the double_sign_slash_amount should be larger than 1e8")
	}
	if p.DowntimeSlashAmount < 1e8 || p.DowntimeSlashAmount > 10000e8 {
		return fmt.Errorf("the downtime_slash_amount should be in range 1e8 to 10000e8")
	}
	if p.SubmitterReward < 1e7 || p.SubmitterReward > 1000e8 {
		return fmt.Errorf("the submitter_reward should be in range 1e7 to 1000e8")
	}
	if p.DowntimeSlashFee < 1e8 || p.DowntimeSlashFee > 1000e8 {
		return fmt.Errorf("the downtime_slash_fee should be in range 1e8 to 1000e8")
	}
	return nil
}

// Params defines the high level settings for staking
type StakeParams struct {
	UnbondingTime time.Duration `json:"unbonding_time"`

	MaxValidators       uint16 `json:"max_validators"`        // maximum number of validators
	BondDenom           string `json:"bond_denom"`            // bondable coin denomination
	MinSelfDelegation   int64  `json:"min_self_delegation"`   // the minimal self-delegation amount
	MinDelegationChange int64  `json:"min_delegation_change"` // the minimal delegation amount changed
}

func (p *StakeParams) GetParamAttribute() (string, bool) {
	return "staking", false
}

func (p *StakeParams) UpdateCheck() error {
	if p.BondDenom != NativeToken {
		return fmt.Errorf("only native token is availabe as bond_denom so far")
	}
	// the valid range is 1 minute to 100 day.
	if p.UnbondingTime > 100*24*time.Hour || p.UnbondingTime < time.Minute {
		return fmt.Errorf("the UnbondingTime should be in range 1 minute to 100 days")
	}
	if p.MaxValidators < 1 || p.MaxValidators > 500 {
		return fmt.Errorf("the max validator should be in range 1 to 500")
	}
	// BondDenom do not check here, it should be native token and do not support update so far.
	// Leave the check in node repo.

	if p.MinSelfDelegation > 10000000e8 || p.MinSelfDelegation < 1e8 {
		return fmt.Errorf("the min_self_delegation should be in range 1e8 to 10000000e8]")
	}
	if p.MinDelegationChange < 1e5 {
		return fmt.Errorf("the min_delegation_change should be no less than 1e5")
	}
	return nil
}

// ---------   Definition cross side chain prams change ------------------- //

type CSCParamChange struct {
	Key    string `json:"key"` // the name of the parameter
	Value  string `json:"value"`
	Target string `json:"target"`

	// Since byte slice is not friendly to show in proposal description, omit it.
	ValueBytes  []byte `json:"-"` // the value of the parameter
	TargetBytes []byte `json:"-"` // the address of the target contract
}

func (c *CSCParamChange) Check() error {
	targetBytes, err := hex.DecodeString(c.Target)
	if err != nil {
		return fmt.Errorf("target is not hex encoded, err %v", err)
	}
	c.TargetBytes = targetBytes

	valueBytes, err := hex.DecodeString(c.Value)
	if err != nil {
		return fmt.Errorf("value is not hex encoded, err %v", err)
	}
	c.ValueBytes = valueBytes
	keyBytes := []byte(c.Key)
	if len(keyBytes) <= 0 || len(keyBytes) > math.MaxUint8 {
		return fmt.Errorf("the length of key exceed the limitation")
	}
	if len(c.ValueBytes) <= 0 || len(c.ValueBytes) > math.MaxUint8 {
		return fmt.Errorf("the length of value exceed the limitation")
	}
	if len(c.TargetBytes) != types.AddrLen {
		return fmt.Errorf("the length of target address is not %d", types.AddrLen)
	}
	return nil
}
