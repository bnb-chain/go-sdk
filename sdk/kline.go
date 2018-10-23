package sdk

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/structs"
)

// KlineQuery def
type KlineQuery struct {
	Symbol    string // required
	Interval  string // required, interval (5m, 1h, 1d, 1w, etc.)
	Limit     int32
	StartTime int64
	EndTime   int64
}

// Kline def
type Kline struct {
	Close            int64 `json:"close"`
	CloseTime        int64 `json:"closeTime"`
	High             int64 `json:"high"`
	Low              int64 `json:"low"`
	NumberOfTrades   int32 `json:"numberOfTrades"`
	Open             int64 `json:"open"`
	OpenTime         int64 `json:"openTime"`
	QuoteAssetVolume int64 `json:"quoteAssetVolume"`
	Volume           int64 `json:"volume"`
}

// GetKlines returns transaction details
func (sdk *SDK) GetKlines(query *KlineQuery) ([]*Kline, error) {

	if query.Symbol == "" {
		return nil, fmt.Errorf("Query.Symbol is required")
	}

	if query.Interval == "" {
		return nil, fmt.Errorf("Query.Interval is required")
	}

	qp := structs.Map(query)
	resp, err := sdk.dexAPI.Get("/klines", ToMapStrStr(qp))
	if err != nil {
		return nil, err
	}

	var klines []*Kline
	if err := json.Unmarshal(resp, &klines); err != nil {
		return nil, err
	}

	return klines, nil
}
