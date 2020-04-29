package rpc

import (
	"crypto/sha256"
	"fmt"
	"strings"

	libbytes "github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/types"
)

const (
	maxABCIPathLength     = 1024
	maxABCIDataLength     = 1024 * 1024
	maxTxLength           = 1024 * 1024
	maxABCIQueryStrLength = 1024
	maxTxSearchStrLength  = 1024
	maxUnConfirmedTxs     = 100
	maxDepthLevel         = 1000

	tokenSymbolMaxLen = 14
	tokenSymbolMinLen = 3
)

var (
	ExceedABCIPathLengthError         = fmt.Errorf("the abci path exceed max length %d ", maxABCIPathLength)
	ExceedABCIDataLengthError         = fmt.Errorf("the abci data exceed max length %d ", maxABCIDataLength)
	ExceedTxLengthError               = fmt.Errorf("the tx data exceed max length %d ", maxTxLength)
	LimitNegativeError                = fmt.Errorf("the limit can't be negative")
	ExceedMaxUnConfirmedTxsNumError   = fmt.Errorf("the limit of unConfirmed tx exceed max limit %d ", maxUnConfirmedTxs)
	HeightNegativeError               = fmt.Errorf("the height can't be negative")
	MaxMinHeightConflictError         = fmt.Errorf("the min height can't be larger than max height")
	HashLengthError                   = fmt.Errorf("the length of hash is not 32")
	ExceedABCIQueryStrLengthError     = fmt.Errorf("the query string exceed max length %d ", maxABCIPathLength)
	ExceedTxSearchQueryStrLengthError = fmt.Errorf("the query string exceed max length %d ", maxTxSearchStrLength)
	OffsetNegativeError               = fmt.Errorf("offset can't be less than 0")
	SymbolLengthExceedRangeError      = fmt.Errorf("length of symbol should be in range [%d,%d]", tokenSymbolMinLen, tokenSymbolMaxLen)
	PairFormatError                   = fmt.Errorf("the pair should in format 'symbol1_symbol2'")
	DepthLevelExceedRangeError        = fmt.Errorf("the level is out of range [%d, %d]", 0, maxDepthLevel)
	KeyMissingError                   = fmt.Errorf("BaseAssetSymbol or QuoteAssetSymbol is missing. ")
)

func ValidateABCIPath(path string) error {
	if len(path) > maxABCIPathLength {
		return ExceedABCIPathLengthError
	}
	return nil
}

func ValidateABCIData(data libbytes.HexBytes) error {
	if len(data) > maxABCIDataLength {
		return ExceedABCIPathLengthError
	}
	return nil
}

func ValidateTx(tx types.Tx) error {
	if len(tx) > maxTxLength {
		return ExceedTxLengthError
	}
	return nil
}

func ValidateUnConfirmedTxsLimit(limit int) error {
	if limit < 0 {
		return LimitNegativeError
	} else if limit > maxUnConfirmedTxs {
		return ExceedMaxUnConfirmedTxsNumError
	}
	return nil
}

func ValidateHeightRange(from int64, to int64) error {
	if from < 0 || to < 0 {
		return HeightNegativeError
	}
	if from > to {
		return MaxMinHeightConflictError
	}
	return nil
}

func ValidateHeight(height *int64) error {
	if height != nil && *height < 0 {
		return HeightNegativeError
	}
	return nil
}

func ValidateHash(hash []byte) error {
	if len(hash) != sha256.Size {
		return HashLengthError
	}
	return nil
}

func ValidateABCIQueryStr(query string) error {
	if len(query) > maxABCIQueryStrLength {
		return ExceedABCIQueryStrLengthError
	}
	return nil
}

func ValidateTxSearchQueryStr(query string) error {
	if len(query) > maxTxSearchStrLength {
		return ExceedTxSearchQueryStrLengthError
	}
	return nil
}

func ValidateOffset(offset int) error {
	if offset < 0 {
		return OffsetNegativeError
	}
	return nil
}

func ValidateLimit(limit int) error {
	if limit < 0 {
		return LimitNegativeError
	}
	return nil
}

func ValidateSymbol(symbol string) error {
	if len(symbol) > tokenSymbolMaxLen || len(symbol) < tokenSymbolMinLen {
		return SymbolLengthExceedRangeError
	}
	return nil
}

func ValidatePair(pair string) error {
	symbols := strings.Split(pair, "_")
	if len(symbols) != 2 {
		return PairFormatError
	}
	if err := ValidateSymbol(symbols[0]); err != nil {
		return err
	}
	if err := ValidateSymbol(symbols[1]); err != nil {
		return err
	}
	return nil
}

func ValidateDepthLevel(level int) error {
	if level < 0 || level > maxDepthLevel {
		return DepthLevelExceedRangeError
	}
	return nil
}
