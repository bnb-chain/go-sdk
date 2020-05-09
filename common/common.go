package common

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/binance-chain/go-sdk/common/types/bsc"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	isAlphaNumFunc = regexp.MustCompile(`^[[:alnum:]]+$`).MatchString
)

func QueryParamToMap(qp interface{}) (map[string]string, error) {
	queryMap := make(map[string]string, 0)
	bz, err := json.Marshal(qp)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bz, &queryMap)
	if err != nil {
		return nil, err
	}
	return queryMap, nil
}

func CombineSymbol(baseAssetSymbol, quoteAssetSymbol string) string {
	return fmt.Sprintf("%s_%s", baseAssetSymbol, quoteAssetSymbol)
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func IsAlphaNum(s string) bool {
	return isAlphaNumFunc(s)
}

func EthHeaderToBscHeader(ethHeader ethtypes.Header) bsc.Header {
	var addr bsc.Address
	copy(addr[:], ethHeader.Coinbase[:])

	var bloom bsc.Bloom
	copy(bloom[:], ethHeader.Bloom[:])

	var blockNonce bsc.BlockNonce
	copy(blockNonce[:], ethHeader.Nonce[:])

	return bsc.Header{
		ParentHash:  bsc.BytesToHash(ethHeader.ParentHash[:]),
		UncleHash:   bsc.BytesToHash(ethHeader.UncleHash[:]),
		Coinbase:    addr,
		Root:        bsc.BytesToHash(ethHeader.Root[:]),
		TxHash:      bsc.BytesToHash(ethHeader.TxHash[:]),
		ReceiptHash: bsc.BytesToHash(ethHeader.ReceiptHash[:]),
		Bloom:       bloom,
		Difficulty:  ethHeader.Difficulty.Int64(),
		Number:      ethHeader.Number.Int64(),
		GasLimit:    ethHeader.GasLimit,
		GasUsed:     ethHeader.GasUsed,
		Time:        ethHeader.Time,
		Extra:       ethHeader.Extra,
		MixDigest:   bsc.BytesToHash(ethHeader.MixDigest[:]),
		Nonce:       blockNonce,
	}
}
