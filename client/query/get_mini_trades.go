package query

import (
	"encoding/json"

	"github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/common/types"
)

// GetMiniTrades returns trade details
func (c *client) GetMiniTrades(query *types.TradesQuery) (*types.Trades, error) {
	err := query.Check()
	if err != nil {
		return nil, err
	}
	qp, err := common.QueryParamToMap(*query)
	if err != nil {
		return nil, err
	}

	resp, _, err := c.baseClient.Get("/mini/trades", qp)
	if err != nil {
		return nil, err
	}

	var trades types.Trades
	if err := json.Unmarshal(resp, &trades); err != nil {
		return nil, err
	}

	return &trades, nil
}
