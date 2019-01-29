package types

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/binance-chain/go-sdk/common/bech32"
)

// AccAddress a wrapper around bytes meant to represent an account address.
// When marshaled to a string or JSON, it uses Bech32.
type AccAddress []byte

type ChainNetwork uint8

const (
	TestNetwork ChainNetwork = iota
	ProdNetwork
)

var Network = ProdNetwork

func (this ChainNetwork) Bech32Prefixes() string {
	switch this {
	case TestNetwork:
		return "tbnb"
	case ProdNetwork:
		return "bnb"
	default:
		panic("Unknown network type")
	}
}

// Marshal needed for protobuf compatibility
func (bz AccAddress) Marshal() ([]byte, error) {
	return bz, nil
}

// Unmarshal needed for protobuf compatibility
func (bz *AccAddress) Unmarshal(data []byte) error {
	*bz = data
	return nil
}

// MarshalJSON to Marshals to JSON using Bech32
func (bz AccAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(bz.String())
}

// UnmarshalJSON to Unmarshal from JSON assuming Bech32 encoding
func (bz *AccAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := AccAddressFromBech32(s)
	if err != nil {
		return err
	}
	*bz = bz2
	return nil
}

// AccAddressFromHex to create an AccAddress from a hex string
func AccAddressFromHex(address string) (addr AccAddress, err error) {
	if len(address) == 0 {
		return addr, errors.New("decoding bech32 address failed: must provide an address")
	}
	bz, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}
	return AccAddress(bz), nil
}

// AccAddressFromBech32 to create an AccAddress from a bech32 string
func AccAddressFromBech32(address string) (addr AccAddress, err error) {
	bz, err := GetFromBech32(address, Network.Bech32Prefixes())
	if err != nil {
		return nil, err
	}
	return AccAddress(bz), nil
}

// GetFromBech32 to decode a bytestring from a bech32-encoded string
func GetFromBech32(bech32str, prefix string) ([]byte, error) {
	if len(bech32str) == 0 {
		return nil, errors.New("decoding bech32 address failed: must provide an address")
	}
	hrp, bz, err := bech32.DecodeAndConvert(bech32str)
	if err != nil {
		return nil, err
	}

	if hrp != prefix {
		return nil, fmt.Errorf("invalid bech32 prefix. Expected %s, Got %s", prefix, hrp)
	}

	return bz, nil
}

func (bz AccAddress) Bytes() []byte {
	return bz
}

// String representation
func (bz AccAddress) String() string {
	bech32Addr, err := bech32.ConvertAndEncode(Network.Bech32Prefixes(), bz.Bytes())
	if err != nil {
		panic(err)
	}
	return bech32Addr
}
