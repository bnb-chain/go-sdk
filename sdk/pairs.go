package sdk

import (
	"encoding/json"
	"strconv"
)

// SymbolPair definition
type SymbolPair struct {
	TradeAsset string `json:"base_asset_symbol"`
	QuoteAsset string `json:"quote_asset_symbol"`
	Price      string `json:"price"`
	TickSize   string `json:"tick_size"`
	LotSize    string `json:"lot_size"`
}

// GetPairs returns list of trading pairs
func (sdk *SDK) GetPairs(limit int) ([]*SymbolPair, error) {
	qp := map[string]string{}

	if limit > 0 {
		qp["limit"] = strconv.Itoa(limit)
	}

	resp, err := sdk.dexAPI.Get("/pairs", qp)
	if err != nil {
		return nil, err
	}

	var listOfPairs []*SymbolPair
	if err := json.Unmarshal(resp, &listOfPairs); err != nil {
		return nil, err
	}

	return listOfPairs, nil
}
