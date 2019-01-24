package query

import (
	"encoding/json"

	"github.com/binance-chain/go-sdk/common"
)

// OpenOrdersQuery def
type OpenOrdersQuery struct {
	SenderAddress string  `json:"address"` // required
	Symbol        string  `json:"symbol,omitempty"`
	Offset        *uint32 `json:"offset,omitempty,string"`
	Limit         *uint32 `json:"limit,omitempty,string"`
}

func NewOpenOrdersQuery(senderAddress string) *OpenOrdersQuery {
	return &OpenOrdersQuery{SenderAddress: senderAddress}
}

func (param *OpenOrdersQuery) WithSymbol(symbol string) *OpenOrdersQuery {
	param.Symbol = symbol
	return param
}

func (param *OpenOrdersQuery) WithOffset(offset uint32) *OpenOrdersQuery {
	param.Offset = &offset
	return param
}

func (param *OpenOrdersQuery) WithLimit(limit uint32) *OpenOrdersQuery {
	param.Limit = &limit
	return param
}

func (param *OpenOrdersQuery) Check() error {
	if param.SenderAddress == "" {
		return AddressMissingError
	}
	if param.Limit != nil && *param.Limit <= 0 {
		return LimitOutOfRangeError
	}
	return nil
}

type OpenOrders struct {
	Order []Order `json:"order"`
	Total int     `json:"total"`
}

// GetOpenOrders returns array of open orders
func (c *client) GetOpenOrders(query *OpenOrdersQuery) (*OpenOrders, error) {
	err := query.Check()
	if err != nil {
		return nil, err
	}
	qp, err := common.QueryParamToMap(*query)
	if err != nil {
		return nil, err
	}

	resp, err := c.baseClient.Get("/orders/open", qp)
	if err != nil {
		return nil, err
	}

	var openOrders OpenOrders
	if err := json.Unmarshal(resp, &openOrders); err != nil {
		return nil, err
	}

	return &openOrders, nil
}
