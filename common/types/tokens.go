package types

// Token definition
type Token struct {
	Name             string     `json:"name"`
	Symbol           string     `json:"symbol"`
	OrigSymbol       string     `json:"original_symbol"`
	TotalSupply      Fixed8     `json:"total_supply"`
	Owner            AccAddress `json:"owner"`
	Mintable         bool       `json:"mintable"`
	ContractAddress  string     `json:"contract_address,omitempty"`
	ContractDecimals int8       `json:"contract_decimals,omitempty"`
}

type TokenBalance struct {
	Symbol string `json:"symbol"`
	Free   Fixed8 `json:"free"`
	Locked Fixed8 `json:"locked"`
	Frozen Fixed8 `json:"frozen"`
}

// MiniToken definition
type MiniToken struct {
	Name             string     `json:"name"`
	Symbol           string     `json:"symbol"`
	OrigSymbol       string     `json:"original_symbol"`
	TotalSupply      Fixed8     `json:"total_supply"`
	Owner            AccAddress `json:"owner"`
	Mintable         bool       `json:"mintable"`
	TokenType        int8       `json:"token_type"`
	TokenURI         string     `json:"token_uri"`
	ContractAddress  string     `json:"contract_address,omitempty"`
	ContractDecimals int8       `json:"contract_decimals,omitempty"`
}
