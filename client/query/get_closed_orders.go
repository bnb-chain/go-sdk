package query

import (
	"encoding/json"

	"github.com/binance-chain/go-sdk/common"
)

const (
	SideBuy  = "BUY"
	SideSell = "SELL"
)

type ClosedOrdersQuery struct {
	SenderAddress string  `json:"address"`                 // required
	Symbol        string  `json:"symbol,omitempty"`        //option
	Offset        *uint32 `json:"offset,omitempty,string"` //option
	Limit         *uint32 `json:"limit,omitempty,string"`  //option
	Start         *int64  `json:"start,omitempty,string"`  //option
	End           *int64  `json:"end,omitempty,string"`    //option
	Side          string  `json:"side,omitempty"`          //option
	Total         int     `json:"total,string"`            //0 for not required and 1 for required; default not required, return total=-1 in response
}

func NewClosedOrdersQuery(senderAddress string, withTotal bool) *ClosedOrdersQuery {
	totalQuery := 0
	if withTotal {
		totalQuery = 1
	}
	return &ClosedOrdersQuery{SenderAddress: senderAddress, Total: totalQuery}
}

func (param *ClosedOrdersQuery) Check() error {
	if param.SenderAddress == "" {
		return AddressMissingError
	}
	if param.Side != SideBuy && param.Side != SideSell && param.Side != "" {
		return TradeSideMisMatchError
	}
	if param.Limit != nil && *param.Limit <= 0 {
		return LimitOutOfRangeError
	}
	if param.Start != nil && *param.Start <= 0 {
		return StartTimeOutOfRangeError
	}
	if param.End != nil && *param.End <= 0 {
		return EndTimeOutOfRangeError
	}
	if param.Start != nil && param.End != nil && *param.Start > *param.End {
		return EndTimeLessThanStartTimeError
	}
	return nil
}

func (param *ClosedOrdersQuery) WithSymbol(baseAssetSymbol, quoteAssetSymbol string) *ClosedOrdersQuery {
	param.Symbol = common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol)
	return param
}

func (param *ClosedOrdersQuery) WithOffset(offset uint32) *ClosedOrdersQuery {
	param.Offset = &offset
	return param
}

func (param *ClosedOrdersQuery) WithLimit(limit uint32) *ClosedOrdersQuery {
	param.Limit = &limit
	return param
}

func (param *ClosedOrdersQuery) WithStart(start int64) *ClosedOrdersQuery {
	param.Start = &start
	return param
}

func (param *ClosedOrdersQuery) WithEnd(end int64) *ClosedOrdersQuery {
	param.End = &end
	return param
}

func (param *ClosedOrdersQuery) WithSide(side string) *ClosedOrdersQuery {
	param.Side = side
	return param
}

type CloseOrders struct {
	Order []Order `json:"order"`
	Total int     `json:"total"`
}

// GetClosedOrders returns array of open orders
func (c *client) GetClosedOrders(query *ClosedOrdersQuery) (*CloseOrders, error) {
	err := query.Check()
	if err != nil {
		return nil, err
	}
	qp, err := common.QueryParamToMap(*query)
	if err != nil {
		return nil, err
	}
	resp, err := c.baseClient.Get("/orders/closed", qp)
	if err != nil {
		return nil, err
	}

	var orders CloseOrders
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, err
	}

	return &orders, nil
}
