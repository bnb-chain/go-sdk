package rpc

import (
	"crypto/sha256"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
)

const (
	maxABCIPathLength     = 1024
	maxABCIDataLength     = 1024 * 1024
	maxTxLength           = 1024 * 1024
	maxCommonStringLength = 1024
	maxUnConfirmedTxs     = 100
)

var (
	ExceedABCIPathLengthError       = fmt.Errorf("the abci path exceed max length %d ", maxABCIPathLength)
	ExceedABCIDataLengthError       = fmt.Errorf("the abci data exceed max length %d ", maxABCIDataLength)
	ExceedTxLengthError             = fmt.Errorf("the tx data exceed max length %d ", maxTxLength)
	LimitNegativeError              = fmt.Errorf("the limit can't be negative")
	ExceedMaxUnConfirmedTxsNumError = fmt.Errorf("the limit of unConfirmed tx exceed max limit %d ", maxUnConfirmedTxs)
	HeightNegativeError             = fmt.Errorf("the height can't be negative")
	MaxMinHeightConflictError       = fmt.Errorf("the min height can't be larger than max height")
	HashLengthError                 = fmt.Errorf("the length of hash is not 32")
	ExceedCommonStrLengthError      = fmt.Errorf("the query string exceed max length %d ", maxABCIPathLength)
)

func ValidateABCIPath(path string) error {
	if len(path) > maxABCIPathLength {
		return ExceedABCIPathLengthError
	}
	return nil
}

func ValidateABCIData(data common.HexBytes) error {
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
	if from < to {
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

func ValidateCommonStr(query string) error {
	if len(query) > maxCommonStringLength {
		return ExceedCommonStrLengthError
	}
	return nil
}
