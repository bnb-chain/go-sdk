package types

import (
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
)

// AppAccount definition
type AppAccount struct {
	BaseAccount `json:"base"`
	Name        string `json:"name"`
	FrozenCoins Coins  `json:"frozen"`
	LockedCoins Coins  `json:"locked"`
	Flags       uint64 `json:"flags"`
}

type Account interface {
	GetAddress() AccAddress
	SetAddress(address AccAddress) error // errors if already set.

	GetPubKey() crypto.PubKey // can return nil.
	SetPubKey(crypto.PubKey) error

	GetAccountNumber() int64
	SetAccountNumber(int64) error

	GetSequence() int64
	SetSequence(int64) error

	GetCoins() Coins
	SetCoins(Coins) error
	Clone() Account

	GetFlags() uint64
	SetFlags(flags uint64)
}

type NamedAccount interface {
	Account
	GetName() string
	SetName(string)

	GetFrozenCoins() Coins
	SetFrozenCoins(Coins)

	//TODO: this should merge into Coin
	GetLockedCoins() Coins
	SetLockedCoins(Coins)
}

func (acc AppAccount) GetName() string              { return acc.Name }
func (acc *AppAccount) SetName(name string)         { acc.Name = name }
func (acc AppAccount) GetFrozenCoins() Coins        { return acc.FrozenCoins }
func (acc *AppAccount) SetFrozenCoins(frozen Coins) { acc.FrozenCoins = frozen }
func (acc AppAccount) GetLockedCoins() Coins        { return acc.LockedCoins }
func (acc *AppAccount) SetLockedCoins(frozen Coins) { acc.LockedCoins = frozen }
func (acc *AppAccount) GetFlags() uint64            { return acc.Flags }
func (acc *AppAccount) SetFlags(flags uint64)       { acc.Flags = flags }

func (acc *AppAccount) Clone() Account {
	baseAcc := acc.BaseAccount.Clone()
	clonedAcc := &AppAccount{
		BaseAccount: baseAcc,
		Name:        acc.Name,
	}
	if acc.FrozenCoins == nil {
		clonedAcc.FrozenCoins = nil
	} else {
		coins := Coins{}
		for _, coin := range acc.FrozenCoins {
			coins = append(coins, Coin{Denom: coin.Denom, Amount: coin.Amount})
		}
		clonedAcc.FrozenCoins = coins
	}
	if acc.LockedCoins == nil {
		clonedAcc.LockedCoins = nil
	} else {
		coins := Coins{}
		for _, coin := range acc.LockedCoins {
			coins = append(coins, Coin{Denom: coin.Denom, Amount: coin.Amount})
		}
		clonedAcc.LockedCoins = coins
	}
	return clonedAcc
}

type BaseAccount struct {
	Address       AccAddress    `json:"address"`
	Coins         Coins         `json:"coins"`
	PubKey        crypto.PubKey `json:"public_key"`
	AccountNumber int64         `json:"account_number"`
	Sequence      int64         `json:"sequence"`
}

// Implements sdk.Account.
func (acc BaseAccount) GetAddress() AccAddress {
	return acc.Address
}

// Implements sdk.Account.
func (acc *BaseAccount) SetAddress(addr AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}
	acc.Address = addr
	return nil
}

// Implements sdk.Account.
func (acc BaseAccount) GetPubKey() crypto.PubKey {
	return acc.PubKey
}

// Implements sdk.Account.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	acc.PubKey = pubKey
	return nil
}

// Implements sdk.Account.
func (acc *BaseAccount) GetCoins() Coins {
	return acc.Coins
}

// Implements sdk.Account.
func (acc *BaseAccount) SetCoins(coins Coins) error {
	acc.Coins = coins
	return nil
}

// Implements Account
func (acc *BaseAccount) GetAccountNumber() int64 {
	return acc.AccountNumber
}

// Implements Account
func (acc *BaseAccount) SetAccountNumber(accNumber int64) error {
	acc.AccountNumber = accNumber
	return nil
}

// Implements sdk.Account.
func (acc *BaseAccount) GetSequence() int64 {
	return acc.Sequence
}

// Implements sdk.Account.
func (acc *BaseAccount) SetSequence(seq int64) error {
	acc.Sequence = seq
	return nil
}

// Implements sdk.Account.
func (acc *BaseAccount) Clone() BaseAccount {
	// given the fact PubKey and Address doesn't change,
	// it should be fine if not deep copy them. if both of
	// the two interfaces can provide a Clone() method would be terrific.
	clonedAcc := BaseAccount{
		PubKey:        acc.PubKey,
		Address:       acc.Address,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}

	if acc.Coins == nil {
		clonedAcc.Coins = nil
	} else {
		coins := make(Coins, 0, len(acc.Coins))
		for _, coin := range acc.Coins {
			coins = append(coins, Coin{Denom: coin.Denom, Amount: coin.Amount})
		}
		clonedAcc.Coins = coins
	}

	return clonedAcc
}

// Balance Account definition
type BalanceAccount struct {
	Number    int64          `json:"account_number"`
	Address   string         `json:"address"`
	Balances  []TokenBalance `json:"balances"`
	PublicKey []uint8        `json:"public_key"`
	Sequence  int64          `json:"sequence"`
	Flags     uint64         `json:"flags"`
}
