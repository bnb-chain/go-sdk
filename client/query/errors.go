package query

import (
	"errors"
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
