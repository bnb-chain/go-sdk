package api

import (
	"encoding/json"
)

// TradesQuery def
type TradesQuery = ClosedOrdersQuery

type Trades struct {
	Trade []Trade `json:"trade"`
	Total string  `json:"total"`
}

func NewTradesQuery(senderAddres string) *TradesQuery {
	return NewClosedOrdersQuery(senderAddres)
}

// Trade def
type Trade struct {
	BuyerOrderID  string  `json:"buyerOrderId"`
	BuyFee        string  `json:"buyFee"`
	BuyerId       string  `json:"buyerId"`
	Price         float64 `json:"price"`
	Quantity      float64 `json:"quantity"`
	SellFee       string  `json:"sellFee"`
	SellerId      string  `json:"sellerId"`
	SellerOrderID string  `json:"sellerOrderId"`
	Symbol        string  `json:"symbol"`
	Time          int64   `json:"time"`
	TradeID       string  `json:"tradeId"`
	BlockHeight   int64   `json:"blockHeight"`
	BaseAsset     string  `json:"baseAsset"`
	QuoteAsset    string  `json:"quoteAsset"`
}

// GetTrades returns transaction details
func (dex *dexAPI) GetTrades(query *TradesQuery) (*Trades, error) {
	err := query.Check()
	if err != nil {
		return nil, err
	}
	qp, err := QueryParamToMap(*query)
	if err != nil {
		return nil, err
	}

	resp, err := dex.Get("/trades", qp)
	if err != nil {
		return nil, err
	}

	var trades Trades
	if err := json.Unmarshal(resp, &trades); err != nil {
		return nil, err
	}

	return &trades, nil
}
