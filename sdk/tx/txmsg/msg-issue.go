package txmsg

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/BiJie/BinanceChain/common/types"
)

// TokenIssueMsg def
type TokenIssueMsg struct {
	From        AccAddress `json:"from"`
	Name        string     `json:"name"`
	Symbol      string     `json:"symbol"`
	TotalSupply int64      `json:"total_supply"`
}

// NewTokenIssueMsg for instance creation
func NewTokenIssueMsg(from AccAddress, name, symbol string, supply int64) TokenIssueMsg {
	return TokenIssueMsg{
		From:        from,
		Name:        name,
		Symbol:      symbol,
		TotalSupply: supply,
	}
}

// ValidateBasic does a simple validation check that
// doesn't require access to any other information.
func (msg TokenIssueMsg) ValidateBasic() error {
	if msg.From == nil {
		return fmt.Errorf("sender address cannot be empty")
	}

	if err := ValidateSymbol(msg.Symbol); err != nil {
		return fmt.Errorf("Invalid symbol %v", msg.Symbol)
	}

	if len(msg.Name) == 0 || len(msg.Name) > 20 {
		return fmt.Errorf("Token name should have 1~20 characters")
	}

	if msg.TotalSupply <= 0 || msg.TotalSupply > types.MaxTotalSupply {
		return fmt.Errorf("Total supply should be <= " + string(types.MaxTotalSupply/int64(math.Pow10(int(types.Decimals)))))
	}

	return nil
}

// Type part of Msg interface
func (msg TokenIssueMsg) Type() string { return "tokenIssue" }

// String part of Msg interface
func (msg TokenIssueMsg) String() string { return fmt.Sprintf("IssueMsg{%#v}", msg) }

// Get part of Msg interface
func (msg TokenIssueMsg) Get(key interface{}) (value interface{}) { return nil }

// GetSigners part of Msg interface
func (msg TokenIssueMsg) GetSigners() []AccAddress { return []AccAddress{msg.From} }

// GetSignBytes part of Msg interface
func (msg TokenIssueMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}
