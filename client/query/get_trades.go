package query

import (
	"encoding/json"

	"github.com/binance-go-sdk-candy/common"
	"github.com/binance-go-sdk-candy/common/types"
)

// GetTrades returns transaction details
func (c *client) GetTrades(query *types.TradesQuery) (*types.Trades, error) {
	err := query.Check()
	if err != nil {
		return nil, err
	}
	qp, err := common.QueryParamToMap(*query)
	if err != nil {
		return nil, err
	}

	resp, _, err := c.baseClient.Get("/trades", qp)
	if err != nil {
		return nil, err
	}

	var trades types.Trades
	if err := json.Unmarshal(resp, &trades); err != nil {
		return nil, err
	}

	return &trades, nil
}
