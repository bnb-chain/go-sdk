package query

import (
	"encoding/json"

	"github.com/bnb-chain/go-sdk/common"
	"github.com/bnb-chain/go-sdk/common/types"
)

// GetTokens returns list of tokens
func (c *client) GetTokens(query *types.TokensQuery) ([]types.Token, error) {
	qp, err := common.QueryParamToMap(*query)
	if err != nil {
		return nil, err
	}
	resp, _, err := c.baseClient.Get("/tokens", qp)
	if err != nil {
		return nil, err
	}

	var tokens []types.Token
	if err := json.Unmarshal(resp, &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}
