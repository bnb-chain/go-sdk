package msg

import (
	"bytes"
	"fmt"

	"github.com/binance-chain/go-sdk/common/types"
	"github.com/tendermint/tendermint/crypto"
)

// Description - description fields for a validator
type Description struct {
	Moniker  string `json:"moniker"`  // name
	Identity string `json:"identity"` // optional identity signature (ex. UPort or Keybase)
	Website  string `json:"website"`  // optional website link
	Details  string `json:"details"`  // optional details
}

// EnsureLength ensures the length of a validator's description.
func (d Description) EnsureLength() (Description, error) {
	if len(d.Moniker) > 70 {
		return d, fmt.Errorf("len moniker %d can be more than %d ", len(d.Moniker), 70)
	}
	if len(d.Identity) > 3000 {
		return d, fmt.Errorf("len Identity %d can be more than %d ", len(d.Identity), 3000)
	}
	if len(d.Website) > 140 {
		return d, fmt.Errorf("len Website %d can be more than %d ", len(d.Website), 140)
	}
	if len(d.Details) > 280 {
		return d, fmt.Errorf("len Details %d can be more than %d ", len(d.Details), 280)
	}

	return d, nil
}

// MsgCreateValidator - struct for bonding transactions
type MsgCreateValidator struct {
	Description   Description
	Commission    types.CommissionMsg
	DelegatorAddr types.AccAddress `json:"delegator_address"`
	ValidatorAddr types.ValAddress `json:"validator_address"`
	PubKey        crypto.PubKey    `json:"pubkey"`
	Delegation    types.Coin       `json:"delegation"`
}

type MsgCreateValidatorProposal struct {
	MsgCreateValidator
	ProposalId int64 `json:"proposal_id"`
}

func (msg MsgCreateValidatorProposal) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return MustSortJSON(bz)
}

func (msg MsgCreateValidatorProposal) Type() string { return "create_validator" }

func (msg MsgCreateValidator) Route() string { return MsgRoute }
func (msg MsgCreateValidator) Type() string  { return "create_validator_open" }

// Return address(es) that must sign over msg.GetSignBytes()
func (msg MsgCreateValidator) GetSigners() []types.AccAddress {
	// delegator is first signer so delegator pays fees
	addrs := []types.AccAddress{msg.DelegatorAddr}

	if !bytes.Equal(msg.DelegatorAddr.Bytes(), msg.ValidatorAddr.Bytes()) {
		// if validator addr is not same as delegator addr, validator must sign
		// msg as well
		addrs = append(addrs, types.AccAddress(msg.ValidatorAddr))
	}
	return addrs
}

// get the bytes for the message signer to sign on
func (msg MsgCreateValidator) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(struct {
		Description
		Commission    types.CommissionMsg
		DelegatorAddr types.AccAddress `json:"delegator_address"`
		ValidatorAddr types.ValAddress `json:"validator_address"`
		PubKey        string           `json:"pubkey"`
		Delegation    types.Coin       `json:"delegation"`
	}{
		Description:   msg.Description,
		Commission:    msg.Commission,
		ValidatorAddr: msg.ValidatorAddr,
		PubKey:        types.MustBech32ifyConsPub(msg.PubKey),
		Delegation:    msg.Delegation,
	})
	if err != nil {
		panic(err)
	}
	return MustSortJSON(b)
}

func (msg MsgCreateValidator) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}

// quick validity check
func (msg MsgCreateValidator) ValidateBasic() error {
	if len(msg.DelegatorAddr) != types.AddrLen {
		return fmt.Errorf("Expected delegator address length is %d, actual length is %d", types.AddrLen, len(msg.DelegatorAddr))
	}
	if len(msg.ValidatorAddr) != types.AddrLen {
		return fmt.Errorf("Expected validator address length is %d, actual length is %d", types.AddrLen, len(msg.ValidatorAddr))
	}
	if msg.Delegation.Amount < 1e8 {
		return fmt.Errorf("self delegation must not be less than 1e8")
	}
	if msg.Description == (Description{}) {
		return fmt.Errorf("description must be included")
	}
	if _, err := msg.Description.EnsureLength(); err != nil {
		return err
	}
	commission := types.NewCommission(msg.Commission.Rate, msg.Commission.MaxRate, msg.Commission.MaxChangeRate)
	if err := commission.Validate(); err != nil {
		return err
	}

	return nil
}

type MsgRemoveValidator struct {
	LauncherAddr types.AccAddress  `json:"launcher_addr"`
	ValAddr      types.ValAddress  `json:"val_addr"`
	ValConsAddr  types.ConsAddress `json:"val_cons_addr"`
	ProposalId   int64             `json:"proposal_id"`
}

func NewMsgRemoveValidator(launcherAddr types.AccAddress, valAddr types.ValAddress,
	valConsAddr types.ConsAddress, proposalId int64) MsgRemoveValidator {
	return MsgRemoveValidator{
		LauncherAddr: launcherAddr,
		ValAddr:      valAddr,
		ValConsAddr:  valConsAddr,
		ProposalId:   proposalId,
	}
}

//nolint
func (msg MsgRemoveValidator) Route() string { return MsgRoute }
func (msg MsgRemoveValidator) Type() string  { return "remove_validator" }
func (msg MsgRemoveValidator) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.LauncherAddr}
}

// get the bytes for the message signer to sign on
func (msg MsgRemoveValidator) GetSignBytes() []byte {
	b, err := MsgCdc.MarshalJSON(struct {
		LauncherAddr types.AccAddress  `json:"launcher_addr"`
		ValAddr      types.ValAddress  `json:"val_addr"`
		ValConsAddr  types.ConsAddress `json:"val_cons_addr"`
		ProposalId   int64             `json:"proposal_id"`
	}{
		LauncherAddr: msg.LauncherAddr,
		ValAddr:      msg.ValAddr,
		ValConsAddr:  msg.ValConsAddr,
		ProposalId:   msg.ProposalId,
	})
	if err != nil {
		panic(err)
	}
	return MustSortJSON(b)
}

// quick validity check
func (msg MsgRemoveValidator) ValidateBasic() error {
	if len(msg.LauncherAddr) != types.AddrLen {
		return fmt.Errorf("Expected launcher address length is %d, actual length is %d", types.AddrLen, len(msg.LauncherAddr))
	}
	if len(msg.ValAddr) != types.AddrLen {
		return fmt.Errorf("Expected validator address length is %d, actual length is %d", types.AddrLen, len(msg.ValAddr))
	}
	if len(msg.ValConsAddr) != types.AddrLen {
		return fmt.Errorf("Expected validator consensus address length is %d, actual length is %d", types.AddrLen, len(msg.ValConsAddr))
	}
	if msg.ProposalId <= 0 {
		return fmt.Errorf("Proposal id is expected to be positive, actual value is %d", msg.ProposalId)
	}
	return nil
}

func (msg MsgRemoveValidator) GetInvolvedAddresses() []types.AccAddress {
	return []types.AccAddress{msg.LauncherAddr}
}

// MsgEditValidator - struct for editing a validator
type MsgEditValidator struct {
	Description
	ValidatorAddr types.ValAddress `json:"address"`
	// We pass a reference to the new commission rate as it's not mandatory to
	// update. If not updated, the deserialized rate will be zero with no way to
	// distinguish if an update was intended.
	CommissionRate *types.Dec    `json:"commission_rate"`
	PubKey         crypto.PubKey `json:"pubkey"`
}

//nolint
func (msg MsgEditValidator) Route() string { return MsgRoute }
func (msg MsgEditValidator) Type() string  { return "edit_validator" }
func (msg MsgEditValidator) GetSigners() []types.AccAddress {
	return []types.AccAddress{types.AccAddress(msg.ValidatorAddr)}
}

// get the bytes for the message signer to sign on
func (msg MsgEditValidator) GetSignBytes() []byte {
	var pubkey string
	if msg.PubKey != nil {
		pubkey = types.MustBech32ifyConsPub(msg.PubKey)
	}
	bz := MsgCdc.MustMarshalJSON(struct {
		Description
		ValidatorAddr  types.ValAddress `json:"validator_address"`
		PubKey         string           `json:"pubkey,omitempty"`
		CommissionRate *types.Dec       `json:"commission_rate,omitempty"`
	}{
		Description:    msg.Description,
		ValidatorAddr:  msg.ValidatorAddr,
		PubKey:         pubkey,
		CommissionRate: msg.CommissionRate,
	})
	return MustSortJSON(bz)
}

// quick validity check
func (msg MsgEditValidator) ValidateBasic() error {
	if msg.ValidatorAddr == nil {
		return fmt.Errorf("nil validator address")
	}

	if msg.Description == (Description{}) {
		return fmt.Errorf("transaction must include some information to modify")
	}

	return nil
}

func (msg MsgEditValidator) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}

//______________________________________________________________________

// MsgDelegate - struct for bonding transactions
type MsgDelegate struct {
	DelegatorAddr types.AccAddress `json:"delegator_addr"`
	ValidatorAddr types.ValAddress `json:"validator_addr"`
	Delegation    types.Coin       `json:"delegation"`
}

//nolint
func (msg MsgDelegate) Route() string { return MsgRoute }
func (msg MsgDelegate) Type() string  { return "delegate" }
func (msg MsgDelegate) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.DelegatorAddr}
}

// get the bytes for the message signer to sign on
func (msg MsgDelegate) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return MustSortJSON(bz)
}

// quick validity check
func (msg MsgDelegate) ValidateBasic() error {
	if msg.DelegatorAddr == nil {
		return fmt.Errorf("delegator address is nil")
	}
	if msg.ValidatorAddr == nil {
		return fmt.Errorf("validator address is nil")
	}
	if msg.Delegation.Amount < 1e8 {
		return fmt.Errorf("delegation amount should not be less than 1e8")
	}
	return nil
}

func (msg MsgDelegate) GetInvolvedAddresses() []types.AccAddress {
	return []types.AccAddress{msg.DelegatorAddr, types.AccAddress(msg.ValidatorAddr)}
}

// MsgDelegate - struct for bonding transactions
type MsgRedelegate struct {
	DelegatorAddr    types.AccAddress `json:"delegator_addr"`
	ValidatorSrcAddr types.ValAddress `json:"validator_src_addr"`
	ValidatorDstAddr types.ValAddress `json:"validator_dst_addr"`
	Amount           types.Coin       `json:"amount"`
}

//nolint
func (msg MsgRedelegate) Route() string { return MsgRoute }
func (msg MsgRedelegate) Type() string  { return "redelegate" }
func (msg MsgRedelegate) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.DelegatorAddr}
}

// get the bytes for the message signer to sign on
func (msg MsgRedelegate) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return MustSortJSON(bz)
}

func (msg MsgRedelegate) GetInvolvedAddresses() []types.AccAddress {
	return []types.AccAddress{msg.DelegatorAddr, types.AccAddress(msg.ValidatorSrcAddr), types.AccAddress(msg.DelegatorAddr)}
}

// ValidateBasic is used to quickly disqualify obviously invalid messages quickly
func (msg MsgRedelegate) ValidateBasic() error {
	if len(msg.DelegatorAddr) != types.AddrLen {
		return fmt.Errorf("Expected delegator address length is %d, actual length is %d", types.AddrLen, len(msg.DelegatorAddr))
	}
	if len(msg.ValidatorSrcAddr) != types.AddrLen {
		return fmt.Errorf("Expected validator source address length is %d, actual length is %d", types.AddrLen, len(msg.ValidatorSrcAddr))
	}
	if len(msg.ValidatorDstAddr) != types.AddrLen {
		return fmt.Errorf("Expected validator destination address length is %d, actual length is %d", types.AddrLen, len(msg.ValidatorDstAddr))
	}
	if msg.Amount.Amount <= 0 {
		return fmt.Errorf("Expected positive amount, actual amount is %v", msg.Amount.Amount)
	}
	return nil
}

type MsgUndelegate struct {
	DelegatorAddr types.AccAddress `json:"delegator_addr"`
	ValidatorAddr types.ValAddress `json:"validator_addr"`
	Amount        types.Coin       `json:"amount"`
}

//nolint
func (msg MsgUndelegate) Route() string { return MsgRoute }
func (msg MsgUndelegate) Type() string  { return "undelegate" }
func (msg MsgUndelegate) GetSigners() []types.AccAddress {
	return []types.AccAddress{msg.DelegatorAddr}
}

// get the bytes for the message signer to sign on
func (msg MsgUndelegate) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return MustSortJSON(bz)
}

// quick validity check
func (msg MsgUndelegate) ValidateBasic() error {
	if len(msg.DelegatorAddr) != types.AddrLen {
		return fmt.Errorf("Expected delegator address length is %d, actual length is %d", types.AddrLen, len(msg.DelegatorAddr))
	}
	if len(msg.ValidatorAddr) != types.AddrLen {
		return fmt.Errorf("Expected validator address length is %d, actual length is %d", types.AddrLen, len(msg.ValidatorAddr))
	}
	if msg.Amount.Amount <= 0 {
		return fmt.Errorf("undelegation amount must be positive: %d", msg.Amount.Amount)
	}
	return nil
}

func (msg MsgUndelegate) GetInvolvedAddresses() []types.AccAddress {
	return []types.AccAddress{msg.DelegatorAddr, types.AccAddress(msg.ValidatorAddr)}
}
