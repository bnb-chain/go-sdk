package api

import (
	"encoding/json"
)

type Ticker24hQuery struct {
	Symbol string `json:"symbol,omitempty"`
}

func NewTicker24hQuery() *Ticker24hQuery {
	return &Ticker24hQuery{}
}

func (param *Ticker24hQuery) WithSymbol(baseAssetSymbol, quoteAssetSymbol string) *Ticker24hQuery {
	param.Symbol = CombineSymbol(baseAssetSymbol, quoteAssetSymbol)
	return param
}

// Ticker24h def
type Ticker24h struct {
	Symbol             string `json:"symbol"`
	AskPrice           string `json:"askPrice"`    // In decimal form, e.g. 1.00000000
	AskQuantity        string `json:"askQuantity"` // In decimal form, e.g. 1.00000000
	BidPrice           string `json:"bidPrice"`    // In decimal form, e.g. 1.00000000
	BidQuantity        string `json:"bidQuantity"` // In decimal form, e.g. 1.00000000
	CloseTime          int64  `json:"closeTime"`
	Count              int64  `json:"count"`
	FirstID            string `json:"firstId"`
	HighPrice          string `json:"highPrice"` // In decimal form, e.g. 1.00000000
	LastID             string `json:"lastId"`
	LastPrice          string `json:"lastPrice"`    // In decimal form, e.g. 1.00000000
	LastQuantity       string `json:"lastQuantity"` // In decimal form, e.g. 1.00000000
	LowPrice           string `json:"lowPrice"`     // In decimal form, e.g. 1.00000000
	OpenPrice          string `json:"openPrice"`    // In decimal form, e.g. 1.00000000
	OpenTime           int64  `json:"openTime"`
	PrevClosePrice     string `json:"prevClosePrice"` // In decimal form, e.g. 1.00000000
	PriceChange        string `json:"priceChange"`    // In decimal form, e.g. 1.00000000
	PriceChangePercent string `json:"priceChangePercent"`
	QuoteVolume        string `json:"quoteVolume"`      // In decimal form, e.g. 1.00000000
	Volume             string `json:"volume"`           // In decimal form, e.g. 1.00000000
	WeightedAvgPrice   string `json:"weightedAvgPrice"` // In decimal form, e.g. 1.00000000
}

// GetTicker24h returns ticker 24h
func (dex *dexAPI) GetTicker24h(query *Ticker24hQuery) ([]Ticker24h, error) {

	qp, err := QueryParamToMap(query)
	if err != nil {
		return nil, err
	}

	resp, err := dex.Get("/ticker/24hr", qp)
	if err != nil {
		return nil, err
	}

	tickers := []Ticker24h{}
	if err := json.Unmarshal(resp, &tickers); err != nil {
		return nil, err
	}

	return tickers, nil
}
