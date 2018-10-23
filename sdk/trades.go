package sdk

import (
	"encoding/json"

	"github.com/fatih/structs"
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
	qp := structs.Map(query)
	resp, err := sdk.dexAPI.Get("/trades", ToMapStrStr(qp))
	if err != nil {
		return nil, err
	}

	var trades []*Trade
	if err := json.Unmarshal(resp, &trades); err != nil {
		return nil, err
	}

	return trades, nil
}
