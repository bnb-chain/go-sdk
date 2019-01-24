package query

import (
	"encoding/json"
)

// OrderSide enum
var OrderSide = struct {
	BUY  string
	SELL string
}{
	"BUY",
	"SELL",
}

// TimeInForce enum
var TimeInForce = struct {
	GTC string
	IOC string
}{"GTC", "IOC"}

// OrderStatus enum
var OrderStatus = struct {
	ACK              string
	PARTIALLY_FILLED string
	IOC_NO_FILL      string
	FULLY_FILLED     string
	CANCELED         string
	EXPIRED          string
	UNKNOWN          string
}{
	"ACK",
	"PARTIALLY_FILLED",
	"IOC_NO_FILL",
	"FULLY_FILLED",
	"CANCELED",
	"EXPIRED",
	"UNKNOWN",
}

// OrderType enum
var OrderType = struct {
	LIMIT             string
	MARKET            string
	STOP_LOSS         string
	STOP_LOSS_LIMIT   string
	TAKE_PROFIT       string
	TAKE_PROFIT_LIMIT string
	LIMIT_MAKER       string
}{
	"LIMIT",
	"MARKET",
	"STOP_LOSS",
	"STOP_LOSS_LIMIT",
	"TAKE_PROFIT",
	"TAKE_PROFIT_LIMIT",
	"LIMIT_MAKER",
}

// Order def
type Order struct {
	ID                   string `json:"orderId"`
	Owner                string `json:"owner"`
	Symbol               string `json:"symbol"`
	Price                string `json:"price"`
	Quantity             string `json:"quantity"`
	CumulateQuantity     string `json:"cumulateQuantity"`
	Fee                  string `json:"fee"`
	Side                 int    `json:"side"` // BUY or SELL
	Status               string `json:"status"`
	TimeInForce          int    `json:"timeInForce"`
	Type                 int    `json:"type"`
	TradeId              string `json:"tradeId"`
	LastExecutedPrice    string `json:"last_executed_price"`
	LastExecutedQuantity string `json:"lastExecutedQuantity"`
	TransactionHash      string `json:"transactionHash"`
	OrderCreateTime      string `json:"orderCreateTime"`
	TransactionTime      string `json:"transactionTime"`
}

// GetOrder returns transaction details
func (c *client) GetOrder(orderID string) (*Order, error) {
	if orderID == "" {
		return nil, OrderIdMissingError
	}

	qp := map[string]string{}
	resp, err := c.baseClient.Get("/orders/"+orderID, qp)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(resp, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
