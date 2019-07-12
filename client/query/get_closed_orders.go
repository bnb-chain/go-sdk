package query

import (
	"encoding/json"

	"github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/common/types"
)

// GetClosedOrders returns array of open orders
func (c *client) GetClosedOrders(query *types.ClosedOrdersQuery) (*types.CloseOrders, error) {
	err := query.Check()
	if err != nil {
		return nil, err
	}
	qp, err := common.QueryParamToMap(*query)
	if err != nil {
		return nil, err
	}
	resp, _, err := c.baseClient.Get("/orders/closed", qp)
	if err != nil {
		return nil, err
	}

	var orders types.CloseOrders
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, err
	}

	return &orders, nil
}
