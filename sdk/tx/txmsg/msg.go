package txmsg

import (
	"fmt"
	"strings"
)

// constants
const (
	DotBSuffix           = ".B"
	NativeToken          = "BNB"
	Decimals       int8  = 8
	MaxTotalSupply int64 = 9000000000000000000 // 90 billions with 8 decimal digits
)

// Msg interface must be fulfilled by transactions messages
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

// ValidateSymbol utility
func ValidateSymbol(symbol string) error {
	if len(symbol) == 0 {
		return fmt.Errorf("Token symbol cannot be empty")
	}

	if len(symbol) > 8 {
		return fmt.Errorf("Token symbol is too long")
	}

	if strings.HasSuffix(symbol, DotBSuffix) {
		symbol = strings.TrimSuffix(symbol, DotBSuffix)
	}

	// if !tx.IsAlphaNum(symbol) {
	// 	return fmt.Errorf("Token symbol should be alphanumeric")
	// }

	return nil
}
