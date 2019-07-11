package query

import (
	"encoding/json"
	"github.com/binance-chain/go-sdk/common/types"
)

// GetOrder returns transaction details
func (c *client) GetOrder(orderID string) (*types.Order, error) {
	if orderID == "" {
		return nil, types.OrderIdMissingError
	}

	qp := map[string]string{}
	resp, _, err := c.baseClient.Get("/orders/"+orderID, qp)
	if err != nil {
		return nil, err
	}

	var order types.Order
	if err := json.Unmarshal(resp, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
