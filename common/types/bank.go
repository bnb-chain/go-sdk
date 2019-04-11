package types

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
)

// Token definition
type Token struct {
	Name        string     `json:"name"`
	Symbol      string     `json:"symbol"`
	OrigSymbol  string     `json:"original_symbol"`
	TotalSupply Fixed8     `json:"total_supply"`
	Owner       AccAddress `json:"owner"`
	Mintable    bool       `json:"mintable"`
}

// AppAccount definition
type AppAccount struct {
	BaseAccount `json:"base"`
	Name        string `json:"name"`
	FrozenCoins Coins  `json:"frozen"`
	LockedCoins Coins  `json:"locked"`
}

// Coin def
// Coin def
type Coin struct {
	Denom  string `json:"denom"`
	Amount int64  `json:"amount"`
}

func (coin Coin) IsZero() bool {
	return coin.Amount == 0
}

func (coin Coin) IsPositive() bool {
	return coin.Amount > 0
}

func (coin Coin) IsNotNegative() bool {
	return coin.Amount >= 0
}

func (coin Coin) SameDenomAs(other Coin) bool {
	return (coin.Denom == other.Denom)
}

func (coin Coin) Plus(coinB Coin) Coin {
	if !coin.SameDenomAs(coinB) {
		return coin
	}
	return Coin{coin.Denom, coin.Amount + coinB.Amount}
}

// Coins def
type Coins []Coin

func (coins Coins) IsValid() bool {
	switch len(coins) {
	case 0:
		return true
	case 1:
		return !coins[0].IsZero()
	default:
		lowDenom := coins[0].Denom
		for _, coin := range coins[1:] {
			if coin.Denom <= lowDenom {
				return false
			}
			if coin.IsZero() {
				return false
			}
			lowDenom = coin.Denom
		}
		return true
	}
}

func (coins Coins) IsPositive() bool {
	if len(coins) == 0 {
		return false
	}
	for _, coin := range coins {
		if !coin.IsPositive() {
			return false
		}
	}
	return true
}

func (coins Coins) Plus(coinsB Coins) Coins {
	sum := ([]Coin)(nil)
	indexA, indexB := 0, 0
	lenA, lenB := len(coins), len(coinsB)
	for {
		if indexA == lenA {
			if indexB == lenB {
				return sum
			}
			return append(sum, coinsB[indexB:]...)
		} else if indexB == lenB {
			return append(sum, coins[indexA:]...)
		}
		coinA, coinB := coins[indexA], coinsB[indexB]
		switch strings.Compare(coinA.Denom, coinB.Denom) {
		case -1:
			sum = append(sum, coinA)
			indexA++
		case 0:
			if coinA.Amount+coinB.Amount == 0 {
				// ignore 0 sum coin type
			} else {
				sum = append(sum, coinA.Plus(coinB))
			}
			indexA++
			indexB++
		case 1:
			sum = append(sum, coinB)
			indexB++
		}
	}
}

// IsEqual returns true if the two sets of Coins have the same value
func (coins Coins) IsEqual(coinsB Coins) bool {
	if len(coins) != len(coinsB) {
		return false
	}
	for i := 0; i < len(coins); i++ {
		if coins[i].Denom != coinsB[i].Denom || !(coins[i].Amount == coinsB[i].Amount) {
			return false
		}
	}
	return true
}

func (coins Coins) IsNotNegative() bool {
	if len(coins) == 0 {
		return true
	}
	for _, coin := range coins {
		if !coin.IsNotNegative() {
			return false
		}
	}
	return true
}

func (coins Coins) AmountOf(denom string) int64 {
	switch len(coins) {
	case 0:
		return 0
	case 1:
		coin := coins[0]
		if coin.Denom == denom {
			return coin.Amount
		}
		return 0
	default:
		midIdx := len(coins) / 2 // 2:1, 3:1, 4:2
		coin := coins[midIdx]
		if denom < coin.Denom {
			return coins[:midIdx].AmountOf(denom)
		} else if denom == coin.Denom {
			return coin.Amount
		} else {
			return coins[midIdx+1:].AmountOf(denom)
		}
	}
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

type NamedAcount interface {
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

func (acc *AppAccount) Clone() Account {
	baseAcc := acc.BaseAccount.Clone().(*BaseAccount)
	clonedAcc := &AppAccount{
		BaseAccount: *baseAcc,
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
func (acc *BaseAccount) Clone() Account {
	// given the fact PubKey and Address doesn't change,
	// it should be fine if not deep copy them. if both of
	// the two interfaces can provide a Clone() method would be terrific.
	clonedAcc := &BaseAccount{
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

type TokenBalance struct {
	Symbol string `json:"symbol"`
	Free   Fixed8 `json:"free"`
	Locked Fixed8 `json:"locked"`
	Frozen Fixed8 `json:"frozen"`
}
