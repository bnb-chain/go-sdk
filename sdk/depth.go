package sdk

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/structs"
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

	qp := structs.Map(query)
	resp, err := sdk.dexAPI.Get("/depth", ToMapStrStr(qp))
	if err != nil {
		return nil, err
	}

	var MarketDepth MarketDepth
	if err := json.Unmarshal(resp, &MarketDepth); err != nil {
		return nil, err
	}

	return &MarketDepth, nil
}
