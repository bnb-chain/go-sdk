package query

import (
	"encoding/json"
)

// Token definition
type Token struct {
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
	TotalSupply    string `json:"total_supply"`
	Owner          string `json:"owner"`
	OriginalSymbol string `json:"original_symbol"`
}

// GetTokens returns list of tokens
func (c *client) GetTokens() ([]Token, error) {
	qp := map[string]string{}
	resp, err := c.baseClient.Get("/tokens", qp)
	if err != nil {
		return nil, err
	}

	var tokens []Token
	if err := json.Unmarshal(resp, &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}
