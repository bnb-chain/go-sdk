package sdk

import (
	"encoding/json"
	"fmt"
)

// OpenOrdersQuery def
type OpenOrdersQuery struct {
	SenderAddress string // required
	Symbol        string
}

// GetOpenOrders returns array of open orders
func (sdk *SDK) GetOpenOrders(query *OpenOrdersQuery) ([]*Order, error) {
	if query.SenderAddress == "" {
		return nil, fmt.Errorf("Query.SenderAddress is required")
	}

	dqj, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	var qp map[string]string
	json.Unmarshal(dqj, &qp)

	resp, err := sdk.dexAPI.Get("/orders/open", qp)
	if err != nil {
		return nil, err
	}

	var orders []*Order
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
