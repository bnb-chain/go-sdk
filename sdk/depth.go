package sdk

import (
	"encoding/json"
	"fmt"
)

// DepthQuery def
type DepthQuery struct {
	Symbol string
	Limit  int32
}

// MarketDepth to be broadcasted to the user
type MarketDepth struct {
	LastUpdateID int64      `json:"lastUpdateId"` // "lastUpdateId": 160,
	Symbol       string     `json:"symbol"`       // "symbol": "BNB_BTC"
	Bids         [][]string `json:"bids"`         // "bids": [ [ "0.0024", "10" ] ]
	Asks         [][]string `json:"asks"`         // "asks": [ [ "0.0024", "10" ] ]
}

// GetDepth returns market depth records
func (sdk *SDK) GetDepth(query *DepthQuery) (*MarketDepth, error) {
	if query.Symbol == "" {
		return nil, fmt.Errorf("Query.Symbol is required")
	}

	dqj, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	var qp map[string]string
	json.Unmarshal(dqj, &qp)

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
