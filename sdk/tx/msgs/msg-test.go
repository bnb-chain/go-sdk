package msgs

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tendermint/tendermint/libs/bech32"
)

// Transactions messages must fulfill the Msg
type Msg interface {

	// Return the message type.
	// Must be alphanumeric or empty.
	Type() string

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() error

	// Get the canonical byte representation of the Msg.
	GetSignBytes() []byte

	// Signers returns the addrs of signers that must sign.
	// CONTRACT: All signatures must be present to be valid.
	// CONTRACT: Returns addrs in some deterministic order.
	GetSigners() []AccAddress
}

// AccAddress a wrapper around bytes meant to represent an account address.
// When marshaled to a string or JSON, it uses Bech32.
type AccAddress []byte

// Bech32 prefixes
const (
	// expected address length
	AddrLen = 20

	// Bech32 prefixes
	Bech32PrefixAccAddr = "cosmosaccaddr"
	Bech32PrefixAccPub  = "cosmosaccpub"
	Bech32PrefixValAddr = "cosmosvaladdr"
	Bech32PrefixValPub  = "cosmosvalpub"
)

// Marshal needed for protobuf compatibility
func (bz AccAddress) Marshal() ([]byte, error) {
	return bz, nil
}

// Unmarshal needed for protobuf compatibility
func (bz *AccAddress) Unmarshal(data []byte) error {
	*bz = data
	return nil
}

// Marshals to JSON using Bech32
func (bz AccAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(bz.String())
}

// Unmarshals from JSON assuming Bech32 encoding
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

// create an AccAddress from a hex string
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

// create an AccAddress from a bech32 string
func AccAddressFromBech32(address string) (addr AccAddress, err error) {
	bz, err := GetFromBech32(address, Bech32PrefixAccAddr)
	if err != nil {
		return nil, err
	}
	return AccAddress(bz), nil
}

// decode a bytestring from a bech32-encoded string
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

// Allow it to fulfill various interfaces in light-client, etc...
func (bz AccAddress) Bytes() []byte {
	return bz
}

func (bz AccAddress) String() string {
	bech32Addr, err := bech32.ConvertAndEncode(Bech32PrefixAccAddr, bz.Bytes())
	if err != nil {
		panic(err)
	}
	return bech32Addr
}

var _ Msg = (*TestMsg)(nil)

// msg type for testing
type TestMsg struct {
	signers []AccAddress
}

func NewTestMsg(addrs ...AccAddress) *TestMsg {
	return &TestMsg{
		signers: addrs,
	}
}

//nolint
func (msg *TestMsg) Type() string { return "TestMsg" }
func (msg *TestMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg.signers)
	if err != nil {
		panic(err)
	}
	return bz
}
func (msg *TestMsg) ValidateBasic() error { return nil }
func (msg *TestMsg) GetSigners() []AccAddress {
	return msg.signers
}
