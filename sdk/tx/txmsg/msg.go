package txmsg

import (
	"fmt"
	"strings"

	"github.com/BiJie/BinanceChain/common/utils"
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

// ValidateSymbol utility
func ValidateSymbol(symbol string) error {
	DotBSuffix := ".B"

	if len(symbol) == 0 {
		return fmt.Errorf("Token symbol cannot be empty")
	}

	if len(symbol) > 8 {
		return fmt.Errorf("Token symbol is too long")
	}

	if strings.HasSuffix(symbol, DotBSuffix) {
		symbol = strings.TrimSuffix(symbol, DotBSuffix)
	}

	if !utils.IsAlphaNum(symbol) {
		return fmt.Errorf("Token symbol should be alphanumeric")
	}

	return nil
}
