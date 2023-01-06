package query

import (
	"encoding/json"

	"github.com/bnb-chain/go-sdk/common"
	"github.com/bnb-chain/go-sdk/common/types"
)

// GetMiniTokens returns list of mini tokens
func (c *client) GetMiniTokens(query *types.TokensQuery) ([]types.MiniToken, error) {
	qp, err := common.QueryParamToMap(*query)
	if err != nil {
		return nil, err
	}
	resp, _, err := c.baseClient.Get("/mini/tokens", qp)
	if err != nil {
		return nil, err
	}

	var tokens []types.MiniToken
	if err := json.Unmarshal(resp, &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}
