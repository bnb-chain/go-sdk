package types

import (
	"errors"

	"github.com/binance-chain/go-sdk/common"
)

const (
	SideBuy  = "BUY"
	SideSell = "SELL"
)

var (
	// Param error
	AddressMissingError           = errors.New("Address is required ")
	SymbolMissingError            = errors.New("Symbol is required ")
	OffsetOutOfRangeError         = errors.New("offset out of range ")
	LimitOutOfRangeError          = errors.New("limit out of range ")
	TradeSideMisMatchError        = errors.New("Trade side is invalid ")
	StartTimeOutOfRangeError      = errors.New("start time out of range ")
	EndTimeOutOfRangeError        = errors.New("end time out of range ")
	IntervalMissingError          = errors.New("interval is required ")
	EndTimeLessThanStartTimeError = errors.New("end time should great than start time")
	OrderIdMissingError           = errors.New("order id is required ")
)

// ClosedOrdersQuery definition
type ClosedOrdersQuery struct {
	SenderAddress string  `json:"address"`                 // required
	Symbol        string  `json:"symbol,omitempty"`        //option
	Offset        *uint32 `json:"offset,omitempty,string"` //option
	Limit         *uint32 `json:"limit,omitempty,string"`  //option
	Start         *int64  `json:"start,omitempty,string"`  //option
	End           *int64  `json:"end,omitempty,string"`    //option
	Side          string  `json:"side,omitempty"`          //option
	Total         int     `json:"total,string"`            //0 for not required and 1 for required; default not required, return total=-1 in response
}

func NewClosedOrdersQuery(senderAddress string, withTotal bool) *ClosedOrdersQuery {
	totalQuery := 0
	if withTotal {
		totalQuery = 1
	}
	return &ClosedOrdersQuery{SenderAddress: senderAddress, Total: totalQuery}
}

func (param *ClosedOrdersQuery) Check() error {
	if param.SenderAddress == "" {
		return AddressMissingError
	}
	if param.Side != SideBuy && param.Side != SideSell && param.Side != "" {
		return TradeSideMisMatchError
	}
	if param.Limit != nil && *param.Limit <= 0 {
		return LimitOutOfRangeError
	}
	if param.Start != nil && *param.Start <= 0 {
		return StartTimeOutOfRangeError
	}
	if param.End != nil && *param.End <= 0 {
		return EndTimeOutOfRangeError
	}
	if param.Start != nil && param.End != nil && *param.Start > *param.End {
		return EndTimeLessThanStartTimeError
	}
	return nil
}

func (param *ClosedOrdersQuery) WithSymbol(baseAssetSymbol, quoteAssetSymbol string) *ClosedOrdersQuery {
	param.Symbol = common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol)
	return param
}

func (param *ClosedOrdersQuery) WithOffset(offset uint32) *ClosedOrdersQuery {
	param.Offset = &offset
	return param
}

func (param *ClosedOrdersQuery) WithLimit(limit uint32) *ClosedOrdersQuery {
	param.Limit = &limit
	return param
}

func (param *ClosedOrdersQuery) WithStart(start int64) *ClosedOrdersQuery {
	param.Start = &start
	return param
}

func (param *ClosedOrdersQuery) WithEnd(end int64) *ClosedOrdersQuery {
	param.End = &end
	return param
}

func (param *ClosedOrdersQuery) WithSide(side string) *ClosedOrdersQuery {
	param.Side = side
	return param
}

// TradesQuery definition
type TradesQuery struct {
	SenderAddress *string `json:"address,omitempty"`       // option
	Symbol        string  `json:"symbol,omitempty"`        //option
	Offset        *uint32 `json:"offset,omitempty,string"` //option
	Limit         *uint32 `json:"limit,omitempty,string"`  //option
	Start         *int64  `json:"start,omitempty,string"`  //option
	End           *int64  `json:"end,omitempty,string"`    //option
	Side          *string `json:"side,omitempty"`          //option
	Total         int     `json:"total,string"`            //0 for not required and 1 for required; default not required, return total=-1 in response
}

func NewTradesQuery(withTotal bool) *TradesQuery {
	totalQuery := 0
	if withTotal {
		totalQuery = 1
	}
	return &TradesQuery{Total: totalQuery}
}

func (param *TradesQuery) Check() error {
	if param.Side != nil && *param.Side != SideBuy && *param.Side != SideSell && *param.Side != "" {
		return TradeSideMisMatchError
	}
	if param.Limit != nil && *param.Limit <= 0 {
		return LimitOutOfRangeError
	}
	if param.Start != nil && *param.Start <= 0 {
		return StartTimeOutOfRangeError
	}
	if param.End != nil && *param.End <= 0 {
		return EndTimeOutOfRangeError
	}
	if param.Start != nil && param.End != nil && *param.Start > *param.End {
		return EndTimeLessThanStartTimeError
	}
	return nil
}

func (param *TradesQuery) WithSymbol(baseAssetSymbol, quoteAssetSymbol string) *TradesQuery {
	param.Symbol = common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol)
	return param
}

func (param *TradesQuery) WithOffset(offset uint32) *TradesQuery {
	param.Offset = &offset
	return param
}

func (param *TradesQuery) WithLimit(limit uint32) *TradesQuery {
	param.Limit = &limit
	return param
}

func (param *TradesQuery) WithStart(start int64) *TradesQuery {
	param.Start = &start
	return param
}

func (param *TradesQuery) WithEnd(end int64) *TradesQuery {
	param.End = &end
	return param
}

func (param *TradesQuery) WithSide(side string) *TradesQuery {
	param.Side = &side
	return param
}

func (param *TradesQuery) WithAddress(addr string) *TradesQuery {
	param.SenderAddress = &addr
	return param
}

// Ticker24hQuery definition
type Ticker24hQuery struct {
	Symbol string `json:"symbol,omitempty"`
}

func NewTicker24hQuery() *Ticker24hQuery {
	return &Ticker24hQuery{}
}

func (param *Ticker24hQuery) WithSymbol(baseAssetSymbol, quoteAssetSymbol string) *Ticker24hQuery {
	param.Symbol = common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol)
	return param
}

// OpenOrdersQuery definition
type OpenOrdersQuery struct {
	SenderAddress string  `json:"address"` // required
	Symbol        string  `json:"symbol,omitempty"`
	Offset        *uint32 `json:"offset,omitempty,string"`
	Limit         *uint32 `json:"limit,omitempty,string"`
	Total         int     `json:"total,string"` //0 for not required and 1 for required; default not required, return total=-1 in response
}

func NewOpenOrdersQuery(senderAddress string, withTotal bool) *OpenOrdersQuery {
	totalQuery := 0
	if withTotal {
		totalQuery = 1
	}
	return &OpenOrdersQuery{SenderAddress: senderAddress, Total: totalQuery}
}

func (param *OpenOrdersQuery) WithSymbol(symbol string) *OpenOrdersQuery {
	param.Symbol = symbol
	return param
}

func (param *OpenOrdersQuery) WithOffset(offset uint32) *OpenOrdersQuery {
	param.Offset = &offset
	return param
}

func (param *OpenOrdersQuery) WithLimit(limit uint32) *OpenOrdersQuery {
	param.Limit = &limit
	return param
}

func (param *OpenOrdersQuery) Check() error {
	if param.SenderAddress == "" {
		return AddressMissingError
	}
	if param.Limit != nil && *param.Limit <= 0 {
		return LimitOutOfRangeError
	}
	return nil
}

// DepthQuery
type DepthQuery struct {
	Symbol string  `json:"symbol"`
	Limit  *uint32 `json:"limit,omitempty,string"`
}

func NewDepthQuery(baseAssetSymbol, quoteAssetSymbol string) *DepthQuery {
	return &DepthQuery{Symbol: common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol)}
}

func (param *DepthQuery) WithLimit(limit uint32) *DepthQuery {
	param.Limit = &limit
	return param
}

func (param *DepthQuery) Check() error {
	if param.Symbol == "" {
		return SymbolMissingError
	}
	if param.Limit != nil && *param.Limit <= 0 {
		return LimitOutOfRangeError
	}
	return nil
}

// KlineQuery definition
type KlineQuery struct {
	Symbol    string  `json:"symbol"`   // required
	Interval  string  `json:"interval"` // required, interval (5m, 1h, 1d, 1w, etc.)
	Limit     *uint32 `json:"limit,omitempty,string"`
	StartTime *int64  `json:"start_time,omitempty,string"`
	EndTime   *int64  `json:"end_time,omitempty,string"`
}

func NewKlineQuery(baseAssetSymbol, quoteAssetSymbol, interval string) *KlineQuery {
	return &KlineQuery{Symbol: common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol), Interval: interval}
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

// MarketsQuery definition
type MarketsQuery struct {
	Offset *uint32 `json:"offset,omitempty,string"` //Option
	Limit  *uint32 `json:"limit,omitempty,string"`  //Option
}

func NewMarketsQuery() *MarketsQuery {
	return &MarketsQuery{}
}

func (param *MarketsQuery) WithOffset(offset uint32) *MarketsQuery {
	param.Offset = &offset
	return param
}

func (param *MarketsQuery) WithLimit(limit uint32) *MarketsQuery {
	param.Limit = &limit
	return param
}

func (param *MarketsQuery) Check() error {
	if param.Limit != nil && *param.Limit <= 0 {
		return LimitOutOfRangeError
	}
	return nil
}

// TokensQuery definition
type TokensQuery struct {
	Offset *uint32 `json:"offset,omitempty,string"` //Option
	Limit  *uint32 `json:"limit,omitempty,string"`  //Option
}

func NewTokensQuery() *TokensQuery {
	return &TokensQuery{}
}

func (param *TokensQuery) WithOffset(offset uint32) *TokensQuery {
	param.Offset = &offset
	return param
}

func (param *TokensQuery) WithLimit(limit uint32) *TokensQuery {
	param.Limit = &limit
	return param
}

func (param *TokensQuery) Check() error {
	if param.Limit != nil && *param.Limit <= 0 {
		return LimitOutOfRangeError
	}
	if param.Offset != nil && *param.Offset < 0 {
		return OffsetOutOfRangeError
	}
	return nil
}
