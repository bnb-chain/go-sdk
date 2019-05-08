package types

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
