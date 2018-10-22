package sdk

import (
	"encoding/json"
	"fmt"
)

// MarketDepth to be broadcasted to the user
type MarketDepth struct {
	LastUpdateID int64      `json:"lastUpdateId"` // "lastUpdateId": 160,
	Symbol       string     `json:"symbol"`       // "symbol": "BNB_BTC"
	Bids         [][]string `json:"bids"`         // "bids": [ [ "0.0024", "10" ] ]
	Asks         [][]string `json:"asks"`         // "asks": [ [ "0.0024", "10" ] ]
}

// GetDepth returns market depth records
func (sdk *SDK) GetDepth(symbol string) (*MarketDepth, error) {
	if symbol == "" || len(symbol) < 7 {
		return nil, fmt.Errorf("Invalid symbol %s", symbol)
	}

	qp := map[string]string{}
	qp["symbol"] = symbol

	resp, err := sdk.dexAPI.Get("/depth", qp)
	if err != nil {
		return nil, err
	}

	var MarketDepth MarketDepth
	if err := json.Unmarshal(resp, &MarketDepth); err != nil {
		return nil, err
	}

	return &MarketDepth, nil
}
