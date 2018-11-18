package sdk

import (
	"encoding/json"
)

// Token definition
type Token struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	TotalSupply string `json:"total_supply"`
	Owner       string `json:"owner"`
}

// GetTokens returns list of tokens
func (sdk *SDK) GetTokens() ([]*Token, error) {
	qp := map[string]string{}
	resp, err := sdk.dexAPI.Get("/tokens", qp)
	if err != nil {
		return nil, err
	}

	var tokens []*Token
	if err := json.Unmarshal(resp, &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}
