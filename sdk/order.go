package sdk

import (
	"encoding/json"
	"fmt"
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
}{
	"GTC",
	"IOC",
}

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
	ID               string `json:"orderId"`
	Owner            string `json:"owner"`
	Symbol           string `json:"symbol"`
	Price            string `json:"price"`
	Quantity         string `json:"quantity"`
	ExecutedQuantity string `json:"executedQuantity"`
	Side             string `json:"side"` // BUY or SELL
	Status           string `json:"status"`
	TimeInForce      string `json:"timeInForce"`
	Type             string `json:"type"`
}

// GetOrder returns transaction details
func (sdk *SDK) GetOrder(orderId string) (*Order, error) {
	if orderId == "" {
		return nil, fmt.Errorf("Invalid order ID %s", orderId)
	}

	qp := map[string]string{}
	resp, err := sdk.dexAPI.Get("/orders/"+orderId, qp)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(resp, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
