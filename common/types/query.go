package types

import (
	"errors"
)

const (
	SideBuy  = "BUY"
	SideSell = "SELL"
)

var (
	// Param error
	AddressMissingError   = errors.New("Address is required ")
	OffsetOutOfRangeError = errors.New("offset out of range ")
	LimitOutOfRangeError  = errors.New("limit out of range ")
)

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
