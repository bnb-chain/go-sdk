package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/binance-chain/go-sdk/common/bech32"
	"github.com/tendermint/tendermint/crypto"
)

type ValAddress []byte

type BondStatus byte

// nolint
const (
	Unbonded  BondStatus = 0x00
	Unbonding BondStatus = 0x01
	Bonded    BondStatus = 0x02
)

type (
	// Commission defines a commission parameters for a given validator.
	Commission struct {
		Rate          Dec       `json:"rate"`            // the commission rate charged to delegators
		MaxRate       Dec       `json:"max_rate"`        // maximum commission rate which validator can ever charge
		MaxChangeRate Dec       `json:"max_change_rate"` // maximum daily increase of the validator commission
		UpdateTime    time.Time `json:"update_time"`     // the last time the commission rate was changed
	}

	// CommissionMsg defines a commission message to be used for creating a
	// validator.
	CommissionMsg struct {
		Rate          Dec `json:"rate"`            // the commission rate charged to delegators
		MaxRate       Dec `json:"max_rate"`        // maximum commission rate which validator can ever charge
		MaxChangeRate Dec `json:"max_change_rate"` // maximum daily increase of the validator commission
	}
)

func NewCommission(rate, maxRate, maxChangeRate Dec) Commission {
	return Commission{
		Rate:          rate,
		MaxRate:       maxRate,
		MaxChangeRate: maxChangeRate,
		UpdateTime:    time.Unix(0, 0).UTC(),
	}
}

// Validate performs basic sanity validation checks of initial commission
// parameters. If validation fails, an error is returned.
func (c Commission) Validate() error {
	switch {
	case c.MaxRate.LT(ZeroDec()):
		// max rate cannot be negative
		return fmt.Errorf("Commission maxrate %v is negative", c.MaxRate)

	case c.MaxRate.GT(OneDec()):
		// max rate cannot be greater than 100%
		return fmt.Errorf("Commission maxrate %v can't be greater than 100%", c.MaxRate)

	case c.Rate.LT(ZeroDec()):
		// rate cannot be negative
		return fmt.Errorf("Commission rate %v can't be negative ", c.Rate)

	case c.Rate.GT(c.MaxRate):
		// rate cannot be greater than the max rate
		return fmt.Errorf("Commission rate %v can't be greater than maxrate %v", c.Rate, c.MaxRate)

	case c.MaxChangeRate.LT(ZeroDec()):
		// change rate cannot be negative
		return fmt.Errorf("Commission change rate %v can't be negative", c.MaxChangeRate)

	case c.MaxChangeRate.GT(c.MaxRate):
		// change rate cannot be greater than the max rate
		return fmt.Errorf("Commission change rate %v can't be greater than MaxRat %v", c.MaxChangeRate, c.MaxRate)
	}

	return nil
}

// ValidateNewRate performs basic sanity validation checks of a new commission
// rate. If validation fails, an SDK error is returned.
func (c Commission) ValidateNewRate(newRate Dec, blockTime time.Time) error {
	switch {
	case blockTime.Sub(c.UpdateTime).Hours() < 24:
		// new rate cannot be changed more than once within 24 hours
		return fmt.Errorf("new rate %v cannot be changed more than once within 24 hours", blockTime.Sub(c.UpdateTime).Hours())

	case newRate.LT(ZeroDec()):
		// new rate cannot be negative
		return fmt.Errorf("new rate %v cannot be negative", newRate)

	case newRate.GT(c.MaxRate):
		// new rate cannot be greater than the max rate
		return fmt.Errorf("new rate %v cannot be greater than the max rate %v", newRate, c.MaxRate)

	case newRate.Sub(c.Rate).Abs().GT(c.MaxChangeRate):
		// new rate % points change cannot be greater than the max change rate
		return fmt.Errorf("new rate %v points change cannot be greater than the max change rate %v", newRate.Sub(c.Rate).Abs(), c.MaxChangeRate)
	}

	return nil
}

func (c Commission) String() string {
	return fmt.Sprintf("rate: %s, maxRate: %s, maxChangeRate: %s, updateTime: %s",
		c.Rate, c.MaxRate, c.MaxChangeRate, c.UpdateTime,
	)
}

// Description - description fields for a validator
type Description struct {
	Moniker  string `json:"moniker"`  // name
	Identity string `json:"identity"` // optional identity signature (ex. UPort or Keybase)
	Website  string `json:"website"`  // optional website link
	Details  string `json:"details"`  // optional details
}

// Validator defines the total amount of bond shares and their exchange rate to
// coins. Accumulation of interest is modelled as an in increase in the
// exchange rate, and slashing as a decrease.  When coins are delegated to this
// validator, the validator is credited with a Delegation whose number of
// bond shares is based on the amount of coins delegated divided by the current
// exchange rate. Voting power can be calculated as total bonds multiplied by
// exchange rate.
type Validator struct {
	FeeAddr      AccAddress `json:"fee_addr"`         // address for fee collection
	OperatorAddr ValAddress `json:"operator_address"` // address of the validator's operator; bech encoded in JSON
	ConsPubKey   string     `json:"consensus_pubkey"` // the consensus public key of the validator; bech encoded in JSON
	Jailed       bool       `json:"jailed"`           // has the validator been jailed from bonded status?

	Status          BondStatus `json:"status"`           // validator status (bonded/unbonding/unbonded)
	Tokens          Dec        `json:"tokens"`           // delegated tokens (incl. self-delegation)
	DelegatorShares Dec        `json:"delegator_shares"` // total shares issued to a validator's delegators

	Description        Description `json:"description"`           // description terms for the validator
	BondHeight         int64       `json:"bond_height"`           // earliest height as a bonded validator
	BondIntraTxCounter int16       `json:"bond_intra_tx_counter"` // block-local tx index of validator change

	UnbondingHeight  int64     `json:"unbonding_height"` // if unbonding, height at which this validator has begun unbonding
	UnbondingMinTime time.Time `json:"unbonding_time"`   // if unbonding, min time for the validator to complete unbonding

	Commission Commission `json:"commission"` // commission parameters

	DistributionAddr AccAddress `json:"distribution_addr"` // the address receives rewards from the side address, and distribute rewards to delegators. It's auto generated
	SideChainId      string         `json:"side_chain_id"`     // side chain id to distinguish different side chains
	SideConsAddr     []byte         `json:"side_cons_addr"`    // consensus address of the side chain validator, this replaces the `ConsPubKey`
	SideFeeAddr      []byte         `json:"side_fee_addr"`     // fee address on the side chain
}

type UnbondingDelegation struct {
	DelegatorAddr  AccAddress `json:"delegator_addr"`  // delegator
	ValidatorAddr  ValAddress `json:"validator_addr"`  // validator unbonding from operator addr
	CreationHeight int64      `json:"creation_height"` // height which the unbonding took place
	MinTime        time.Time  `json:"min_time"`        // unix time for unbonding completion
	InitialBalance Coin       `json:"initial_balance"` // atoms initially scheduled to receive at completion
	Balance        Coin       `json:"balance"`         // atoms to receive at completion
}

func (va ValAddress) String() string {
	bech32PrefixValAddr := Network.Bech32ValidatorAddrPrefix()
	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixValAddr, va.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

func (va ValAddress) Bytes() []byte {
	return va
}

// MarshalJSON marshals to JSON using Bech32.
func (va ValAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(va.String())
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (va *ValAddress) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	va2, err := ValAddressFromBech32(s)
	if err != nil {
		return err
	}

	*va = va2
	return nil
}

func ValAddressFromBech32(address string) (addr ValAddress, err error) {
	bech32PrefixValAddr := Network.Bech32ValidatorAddrPrefix()
	bz, err := GetFromBech32(address, bech32PrefixValAddr)
	if err != nil {
		return nil, err
	}

	return ValAddress(bz), nil
}

// consensus node
// ----------------------------------------------------------------------------

// ConsAddress defines a wrapper around bytes meant to present a consensus node.
// When marshaled to a string or JSON, it uses Bech32.
type ConsAddress []byte

// ConsAddressFromHex creates a ConsAddress from a hex string.
func ConsAddressFromHex(address string) (addr ConsAddress, err error) {
	if len(address) == 0 {
		return addr, errors.New("decoding Bech32 address failed: must provide an address")
	}

	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}

	return ConsAddress(bz), nil
}

// ConsAddressFromBech32 creates a ConsAddress from a Bech32 string.
func ConsAddressFromBech32(address string) (addr ConsAddress, err error) {
	bz, err := GetFromBech32(address, bech32PrefixConsAddr)
	if err != nil {
		return nil, err
	}

	return ConsAddress(bz), nil
}

// get ConsAddress from pubkey
func GetConsAddress(pubkey crypto.PubKey) ConsAddress {
	return ConsAddress(pubkey.Address())
}

// Returns boolean for whether two ConsAddress are Equal
func (ca ConsAddress) Equals(ca2 ConsAddress) bool {
	if ca.Empty() && ca2.Empty() {
		return true
	}

	return bytes.Compare(ca.Bytes(), ca2.Bytes()) == 0
}

// Returns boolean for whether an ConsAddress is empty
func (ca ConsAddress) Empty() bool {
	if ca == nil {
		return true
	}

	ca2 := ConsAddress{}
	return bytes.Compare(ca.Bytes(), ca2.Bytes()) == 0
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (ca ConsAddress) Marshal() ([]byte, error) {
	return ca, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (ca *ConsAddress) Unmarshal(data []byte) error {
	*ca = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (ca ConsAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(ca.String())
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (ca *ConsAddress) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	ca2, err := ConsAddressFromBech32(s)
	if err != nil {
		return err
	}

	*ca = ca2
	return nil
}

// Bytes returns the raw address bytes.
func (ca ConsAddress) Bytes() []byte {
	return ca
}

// String implements the Stringer interface.
func (ca ConsAddress) String() string {
	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixConsAddr, ca.Bytes())
	if err != nil {
		panic(err)
	}

	return bech32Addr
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (ca ConsAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", ca.String())))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", ca)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(ca))))
	}
}


func NewBaseParams(sideChainId string) BaseParams {
	return BaseParams{SideChainId:sideChainId}
}

type QueryTopValidatorsParams struct {
	BaseParams
	Top int
}

type QueryBondsParams struct {
	BaseParams
	DelegatorAddr AccAddress
	ValidatorAddr ValAddress
}

type QueryValidatorParams struct {
	BaseParams
	ValidatorAddr ValAddress
}

type Pool struct {
	LooseTokens  Dec `json:"loose_tokens"`  // tokens which are not bonded in a validator
	BondedTokens Dec `json:"bonded_tokens"` // reserve of bonded tokens
}