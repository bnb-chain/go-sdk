package sdk

import (
	"encoding/json"
)

// TradesQuery def
type TradesQuery = ClosedOrdersQuery

// Trade def
type Trade struct {
	BuyerOrderID  string `json:"buyerOrderId"`
	BuyFee        string `json:"buyFee"`
	Price         string `json:"price"`
	Quantity      string `json:"quantity"`
	SellFee       string `json:"sellFee"`
	SellerOrderID string `json:"sellerOrderId"`
	Symbol        string `json:"symbol"`
	Time          int64  `json:"time"`
	TradeID       string `json:"tradeId"`
}

// GetTrades returns transaction details
func (sdk *SDK) GetTrades(query *TradesQuery) ([]*Trade, error) {
	dqj, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	var qp map[string]string
	json.Unmarshal(dqj, &qp)

	resp, err := sdk.dexAPI.Get("/trades", qp)
	if err != nil {
		return nil, err
	}

	var trades []*Trade
	if err := json.Unmarshal(resp, &trades); err != nil {
		return nil, err
	}

	return trades, nil
}
