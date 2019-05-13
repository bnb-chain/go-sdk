package types

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
