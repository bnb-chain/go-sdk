package sdk

import (
	"encoding/json"
	"fmt"
)

// Ticker def
type Ticker struct {
	Symbol   string `json:"symbol"`
	AskPrice string `json:"askPrice"` // In decimal form, e.g. 1.00000000
	AskQty   string `json:"askQty"`   // In decimal form, e.g. 1.00000000
	BidPrice string `json:"bidPrice"` // In decimal form, e.g. 1.00000000
	BidQty   string `json:"bidQty"`   // In decimal form, e.g. 1.00000000
}

// GetTicker returns ticker
func (sdk *SDK) GetTicker(symbol string) (*Ticker, error) {

	if symbol == "" {
		return nil, fmt.Errorf("Symbol is required")
	}

	qp := map[string]string{}
	qp["symbol"] = symbol

	resp, err := sdk.dexAPI.Get("/ticker/ticker", qp)
	if err != nil {
		return nil, err
	}

	var ticker Ticker
	if err := json.Unmarshal(resp, &ticker); err != nil {
		return nil, err
	}

	return &ticker, nil
}
