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

func (msg MsgCreateValidator) Route() string { return MsgRoute }
func (msg MsgCreateValidator) Type() string  { return "create_validator" }

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
		DelegatorAddr types.AccAddress `json:"delegator_address"`
		ValidatorAddr types.ValAddress `json:"validator_address"`
		PubKey        string           `json:"pubkey"`
		Delegation    types.Coin       `json:"delegation"`
	}{
		Description:   msg.Description,
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
	if !(msg.Delegation.Amount > 0) {
		return fmt.Errorf("DelegationAmount %d is invalid", msg.Delegation.Amount)
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
