package types

import "strings"

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
