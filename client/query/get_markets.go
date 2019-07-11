package query

import (
	"encoding/json"

	"github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/common/types"
)

// GetMarkets returns list of trading pairs
func (c *client) GetMarkets(query *types.MarketsQuery) ([]types.TradingPair, error) {
	err := query.Check()
	if err != nil {
		return nil, err
	}
	qp, err := common.QueryParamToMap(*query)
	if err != nil {
		return nil, err
	}
	resp, _, err := c.baseClient.Get("/markets", qp)
	if err != nil {
		return nil, err
	}
	var listOfPairs []types.TradingPair
	if err := json.Unmarshal(resp, &listOfPairs); err != nil {
		return nil, err
	}

	return listOfPairs, nil
}
