package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/binance-chain/go-sdk/common/bech32"
)

type ValAddress []byte

type BondStatus byte

// nolint
const (
	Unbonded  BondStatus = 0x00
	Unbonding BondStatus = 0x01
	Bonded    BondStatus = 0x02
)

// Description - description fields for a validator
type Description struct {
	Moniker  string `json:"moniker"`  // name
	Identity string `json:"identity"` // optional identity signature (ex. UPort or Keybase)
	Website  string `json:"website"`  // optional website link
	Details  string `json:"details"`  // optional details
}

type Commission struct {
	Rate          Dec       `json:"rate"`            // the commission rate charged to delegators
	MaxRate       Dec       `json:"max_rate"`        // maximum commission rate which validator can ever charge
	MaxChangeRate Dec       `json:"max_change_rate"` // maximum daily increase of the validator commission
	UpdateTime    time.Time `json:"update_time"`     // the last time the commission rate was changed
}

func (c Commission) String() string {
	return fmt.Sprintf("rate: %s, maxRate: %s, maxChangeRate: %s, updateTime: %s",
		c.Rate, c.MaxRate, c.MaxChangeRate, c.UpdateTime,
	)
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
