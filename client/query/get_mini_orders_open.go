package query

import (
	"encoding/json"

	"github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/common/types"
)

// GetMiniOpenOrders returns array of mini open orders
func (c *client) GetMiniOpenOrders(query *types.OpenOrdersQuery) (*types.OpenOrders, error) {
	err := query.Check()
	if err != nil {
		return nil, err
	}
	qp, err := common.QueryParamToMap(*query)
	if err != nil {
		return nil, err
	}

	resp, _, err := c.baseClient.Get("/mini/orders/open", qp)
	if err != nil {
		return nil, err
	}

	var openOrders types.OpenOrders
	if err := json.Unmarshal(resp, &openOrders); err != nil {
		return nil, err
	}

	return &openOrders, nil
}
