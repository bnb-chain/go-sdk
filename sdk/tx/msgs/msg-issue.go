package msgs

import (
	"encoding/json"
	"fmt"
)

// var _ sdk.TokenIssueMsg = TokenIssueMsg{}

type TokenIssueMsg struct {
	From        AccAddress `json:"from"`
	Name        string     `json:"name"`
	Symbol      string     `json:"symbol"`
	TotalSupply int64      `json:"total_supply"`
}

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
	// if msg.From == nil {
	// 	return sdk.ErrInvalidAddress("sender address cannot be empty")
	// }

	// if err := types.ValidateSymbol(msg.Symbol); err != nil {
	// 	return sdk.ErrInvalidCoins(err.Error())
	// }

	// if len(msg.Name) == 0 || len(msg.Name) > 20 {
	// 	return sdk.ErrInvalidCoins("token name should have 1~20 characters")
	// }

	// if msg.TotalSupply <= 0 || msg.TotalSupply > types.MaxTotalSupply {
	// 	return sdk.ErrInvalidCoins("total supply should be <= " + string(types.MaxTotalSupply/int64(math.Pow10(int(types.Decimals)))))
	// }

	return nil
}

// Implements TokenIssueMsg.
func (msg TokenIssueMsg) Type() string                            { return "tokenIssue" }
func (msg TokenIssueMsg) String() string                          { return fmt.Sprintf("IssueMsg{%#v}", msg) }
func (msg TokenIssueMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg TokenIssueMsg) GetSigners() []AccAddress                { return []AccAddress{msg.From} }

func (msg TokenIssueMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}
