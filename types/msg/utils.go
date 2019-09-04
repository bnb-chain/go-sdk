package msg

import (
	"encoding/binary"
	"encoding/json"
	"github.com/binance-chain/go-sdk/common/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"strings"
)

// SortJSON takes any JSON and returns it sorted by keys. Also, all white-spaces
// are removed.
// This method can be used to canonicalize JSON to be returned by GetSignBytes,
// e.g. for the ledger integration.
// If the passed JSON isn't valid it will return an error.
func SortJSON(toSortJSON []byte) ([]byte, error) {
	var c interface{}
	err := json.Unmarshal(toSortJSON, &c)
	if err != nil {
		return nil, err
	}
	js, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return js, nil
}

// MustSortJSON is like SortJSON but panic if an error occurs, e.g., if
// the passed JSON isn't valid.
func MustSortJSON(toSortJSON []byte) []byte {
	js, err := SortJSON(toSortJSON)
	if err != nil {
		panic(err)
	}
	return js
}

func CalculateRandomHash(randomNumber []byte, timestamp int64) []byte {
	data := make([]byte, RandomNumberLength+Int64Size)
	copy(data[:RandomNumberLength], randomNumber)
	binary.BigEndian.PutUint64(data[RandomNumberLength:], uint64(timestamp))
	return tmhash.Sum(data)
}

func CalculateSwapID(randomNumberHash []byte, sender types.AccAddress, senderOtherChain string) []byte {
	senderOtherChain = strings.ToLower(senderOtherChain)
	data := randomNumberHash
	data = append(data, []byte(sender)...)
	data = append(data, []byte(senderOtherChain)...)
	return tmhash.Sum(data)
}