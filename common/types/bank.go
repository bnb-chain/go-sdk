package types

import "github.com/binance-chain/go-sdk/types"

// Token definition
type Token struct {
	Name        string           `json:"name"`
	Symbol      string           `json:"symbol"`
	OrigSymbol  string           `json:"original_symbol"`
	TotalSupply Fixed8           `json:"total_supply"`
	Owner       types.AccAddress `json:"owner"`
	Mintable    bool             `json:"mintable"`
}
