package query

import (
	"encoding/json"

	"github.com/binance-chain/go-sdk/common"
)

// TradesQuery def
type TradesQuery = ClosedOrdersQuery

type Trades struct {
	Trade []Trade `json:"trade"`
	Total int     `json:"total"`
}

func NewTradesQuery(senderAddres string, withTotal bool) *TradesQuery {
	return NewClosedOrdersQuery(senderAddres, withTotal)
}

// Trade def
type Trade struct {
	BuyerOrderID  string `json:"buyerOrderId"`
	BuyFee        string `json:"buyFee"`
	BuyerId       string `json:"buyerId"`
	Price         string `json:"price"`
	Quantity      string `json:"quantity"`
	SellFee       string `json:"sellFee"`
	SellerId      string `json:"sellerId"`
	SellerOrderID string `json:"sellerOrderId"`
	Symbol        string `json:"symbol"`
	Time          int64  `json:"time"`
	TradeID       string `json:"tradeId"`
	BlockHeight   int64  `json:"blockHeight"`
	BaseAsset     string `json:"baseAsset"`
	QuoteAsset    string `json:"quoteAsset"`
}

// GetTrades returns transaction details
func (c *client) GetTrades(query *TradesQuery) (*Trades, error) {
	err := query.Check()
	if err != nil {
		return nil, err
	}
	qp, err := common.QueryParamToMap(*query)
	if err != nil {
		return nil, err
	}

	resp, err := c.baseClient.Get("/trades", qp)
	if err != nil {
		return nil, err
	}

	var trades Trades
	if err := json.Unmarshal(resp, &trades); err != nil {
		return nil, err
	}

	return &trades, nil
}
