package tx

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"regexp"

	"github.com/BiJie/bnc-go-sdk/sdk/tx/txmsg"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
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

// PrivAndAddr is handy utility to generate and return secp256k1 keys
func PrivAndAddr() (tmcrypto.PrivKey, txmsg.AccAddress) {
	priv := secp256k1.GenPrivKey()
	addr := txmsg.AccAddress(priv.PubKey().Address())
	return priv, addr
}

// EncodeHex for converting bytes to hex formatted bytes
func EncodeHex(b []byte) []byte {
	txHex := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(txHex, b)
	return txHex
}

// DecodeHex for converting hex format to bytes
func DecodeHex(h []byte) []byte {
	b := make([]byte, hex.DecodedLen(len(h)))
	_, err := hex.Decode(b, h)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

var (
	isAlphaNumFunc = regexp.MustCompile(`^[[:alnum:]]+$`).MatchString
)

func IsAlphaNum(s string) bool {
	return isAlphaNumFunc(s)
}
