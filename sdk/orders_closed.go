package sdk

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/structs"
)

// ClosedOrdersQuery def
type ClosedOrdersQuery struct {
	SenderAddress string // required
	Symbol        string
	Offset        int32
	Limit         int32
	Start         int64
	End           int64
	Side          string
}

// GetClosedOrders returns array of open orders
func (sdk *SDK) GetClosedOrders(query *ClosedOrdersQuery) ([]*Order, error) {
	if query.SenderAddress == "" {
		return nil, fmt.Errorf("Query.SenderAddress is required")
	}

	if query.Side != "" && query.Side != OrderSide.SELL && query.Side != OrderSide.BUY {
		return nil, fmt.Errorf("Invalid `Query.Side` param")
	}

	qp := structs.Map(query)
	resp, err := sdk.dexAPI.Get("/orders/closed", ToMapStrStr(qp))
	if err != nil {
		return nil, err
	}

	var orders []*Order
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
