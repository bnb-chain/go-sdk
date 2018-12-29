package api

import (
	"encoding/json"
	"fmt"
)

// KlineQuery def
type KlineQuery struct {
	Symbol    string  `json:"symbol"`   // required
	Interval  string  `json:"interval"` // required, interval (5m, 1h, 1d, 1w, etc.)
	Limit     *uint32 `json:"limit,omitempty,string"`
	StartTime *int64  `json:"start_time,omitempty,string"`
	EndTime   *int64  `json:"end_time,omitempty,string"`
}

func NewKlineQuery(baseAssetSymbol, quoteAssetSymbol, interval string) *KlineQuery {
	return &KlineQuery{Symbol: CombineSymbol(baseAssetSymbol, quoteAssetSymbol), Interval: interval}
}

func (param *KlineQuery) WithStartTime(start int64) *KlineQuery {
	param.StartTime = &start
	return param
}

func (param *KlineQuery) WithEndTime(end int64) *KlineQuery {
	param.EndTime = &end
	return param
}

func (param *KlineQuery) WithLimit(limit uint32) *KlineQuery {
	param.Limit = &limit
	return param
}

func (param *KlineQuery) Check() error {
	if param.Symbol == "" {
		return SymbolMissingError
	}
	if param.Interval == "" {
		return IntervalMissingError
	}
	if param.Limit != nil && *param.Limit <= 0 {
		return LimitOutOfRangeError
	}
	if param.StartTime != nil && *param.StartTime <= 0 {
		return StartTimeOutOfRangeError
	}
	if param.EndTime != nil && *param.EndTime <= 0 {
		return EndTimeOutOfRangeError
	}
	if param.StartTime != nil && param.EndTime != nil && *param.StartTime > *param.EndTime {
		return EndTimeLessThanStartTimeError
	}
	return nil
}

// Kline def
type Kline struct {
	Close            float64 `json:"close,string"`
	CloseTime        int64   `json:"closeTime"`
	High             float64 `json:"high,string"`
	Low              float64 `json:"low,string"`
	NumberOfTrades   int32   `json:"numberOfTrades"`
	Open             float64 `json:"open,string"`
	OpenTime         int64   `json:"openTime"`
	QuoteAssetVolume float64 `json:"quoteAssetVolume,string"`
	Volume           float64 `json:"volume,string"`
}

// GetKlines returns transaction details
func (dex *dexAPI) GetKlines(query *KlineQuery) ([]Kline, error) {
	err := query.Check()
	if err != nil {
		return nil, err
	}
	qp, err := QueryParamToMap(*query)
	if err != nil {
		return nil, err
	}

	resp, err := dex.Get("/klines", qp)
	if err != nil {
		return nil, err
	}

	iklines := [][]interface{}{}
	if err := json.Unmarshal(resp, &iklines); err != nil {
		return nil, err
	}
	klines := make([]Kline, len(iklines))
	// Todo
	for index, ikline := range iklines {
		kl := Kline{}
		imap := make(map[string]interface{}, 9)
		if len(ikline) >= 9 {
			imap["openTime"] = ikline[0]
			imap["open"] = ikline[1]
			imap["high"] = ikline[2]
			imap["low"] = ikline[3]
			imap["close"] = ikline[4]
			imap["volume"] = ikline[5]
			imap["closeTime"] = ikline[6]
			imap["quoteAssetVolume"] = ikline[7]
			imap["NumberOfTrades"] = ikline[8]
		} else {
			return nil, fmt.Errorf("Receive kline scheme is unexpected ")
		}
		bz, err := json.Marshal(imap)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bz, &kl)
		if err != nil {
			return nil, err
		}
		klines[index] = kl
	}
	return klines, nil
}
