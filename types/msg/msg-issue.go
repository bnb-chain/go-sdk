package msg

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/binance-chain/go-sdk/common/types"

	"github.com/binance-chain/go-sdk/common"
)

// TokenIssueMsg def
type TokenIssueMsg struct {
	From        types.AccAddress `json:"from"`
	Name        string           `json:"name"`
	Symbol      string           `json:"symbol"`
	TotalSupply int64            `json:"total_supply"`
	Mintable    bool             `json:"mintable"`
}

// NewTokenIssueMsg for instance creation
func NewTokenIssueMsg(from types.AccAddress, name, symbol string, supply int64, mintable bool) TokenIssueMsg {
	return TokenIssueMsg{
		From:        from,
		Name:        name,
		Symbol:      symbol,
		TotalSupply: supply,
		Mintable:    mintable,
	}
}

// ValidateBasic does a simple validation check that
// doesn't require access to any other information.
func (msg TokenIssueMsg) ValidateBasic() error {
	if msg.From == nil {
		return fmt.Errorf("sender address cannot be empty")
	}

	if err := validateIssueMsgTokenSymbol(msg.Symbol); err != nil {
		return fmt.Errorf("Invalid symbol %v", msg.Symbol)
	}

	if len(msg.Name) == 0 || len(msg.Name) > 20 {
		return fmt.Errorf("Token name should have 1~20 characters")
	}

	if msg.TotalSupply <= 0 || msg.TotalSupply > MaxTotalSupply {
		return fmt.Errorf("Total supply should be <= " + string(MaxTotalSupply/int64(math.Pow10(int(Decimals)))))
	}

	return nil
}

// Route part of Msg interface
func (msg TokenIssueMsg) Route() string { return "tokenIssue" }

// Type part of Msg interface
func (msg TokenIssueMsg) Type() string { return "tokenIssue" }

// String part of Msg interface
func (msg TokenIssueMsg) String() string { return fmt.Sprintf("IssueMsg{%#v}", msg) }

// GetSigners part of Msg interface
func (msg TokenIssueMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

// GetSignBytes part of Msg interface
func (msg TokenIssueMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetInvolvedAddresses part of Msg interface
func (msg TokenIssueMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}

func validateIssueMsgTokenSymbol(symbol string) error {
	if len(symbol) == 0 {
		return errors.New("token symbol cannot be empty")
	}

	if strings.HasSuffix(symbol, DotBSuffix) {
		symbol = strings.TrimSuffix(symbol, DotBSuffix)
	}

	// check len without .B suffix
	if symbolLen := len(symbol); symbolLen > TokenSymbolMaxLen || symbolLen < TokenSymbolMinLen {
		return errors.New("length of token symbol is limited to 2~8")
	}

	if !common.IsAlphaNum(symbol) {
		return errors.New("token symbol should be alphanumeric")
	}

	return nil
}
