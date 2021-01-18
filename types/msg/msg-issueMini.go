package msg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/binance-chain/go-sdk/common/types"

	"github.com/binance-chain/go-sdk/common"
)

const (
	MiniRoute        = "miniTokensIssue"
	IssueMiniMsgType = "miniIssueMsg"
)

// MiniTokenIssueMsg def
type MiniTokenIssueMsg struct {
	From        types.AccAddress `json:"from"`
	Name        string           `json:"name"`
	Symbol      string           `json:"symbol"`
	TotalSupply int64            `json:"total_supply"`
	Mintable    bool             `json:"mintable"`
	TokenURI    string           `json:"token_uri"`
}

// NewMiniTokenIssueMsg for instance creation
func NewMiniTokenIssueMsg(from types.AccAddress, name, symbol string, supply int64, mintable bool, tokenURI string) MiniTokenIssueMsg {
	return MiniTokenIssueMsg{
		From:        from,
		Name:        name,
		Symbol:      symbol,
		TotalSupply: supply,
		Mintable:    mintable,
		TokenURI:    tokenURI,
	}
}

// ValidateBasic does a simple validation check that
// doesn't require access to any other information.
func (msg MiniTokenIssueMsg) ValidateBasic() error {

	if msg.From == nil {
		return errors.New("sender address cannot be empty")
	}

	if err := validateIssueMsgMiniTokenSymbol(msg.Symbol); err != nil {
		return fmt.Errorf("Invalid symbol %v", msg.Symbol)
	}

	if len(msg.Name) == 0 || len(msg.Name) > MaxMiniTokenNameLength {
		return fmt.Errorf("token name should have 1 ~ %v characters", MaxMiniTokenNameLength)
	}

	if len(msg.TokenURI) > MaxTokenURILength {
		return fmt.Errorf("token seturi should not exceed %v characters", MaxTokenURILength)
	}

	//if msg.MaxTotalSupply < MiniTokenMinTotalSupply || msg.MaxTotalSupply > MiniTokenMaxTotalSupplyUpperBound {
	//	return fmt.Errorf("max total supply should be between %d ~ %d", MiniTokenMinTotalSupply, MiniTokenMaxTotalSupplyUpperBound)
	//}
	//
	//if msg.TotalSupply < MiniTokenMinTotalSupply || msg.TotalSupply > msg.MaxTotalSupply {
	//	return fmt.Errorf("total supply should be between %d ~ %d", MiniTokenMinTotalSupply, msg.MaxTotalSupply)
	//}

	return nil
}

// Route part of Msg interface
func (msg MiniTokenIssueMsg) Route() string { return MiniRoute }

// Type part of Msg interface
func (msg MiniTokenIssueMsg) Type() string {
	return IssueMiniMsgType
}

// String part of Msg interface
func (msg MiniTokenIssueMsg) String() string { return fmt.Sprintf("IssueMiniMsg{%#v}", msg) }

// GetSigners part of Msg interface
func (msg MiniTokenIssueMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

// GetSignBytes part of Msg interface
func (msg MiniTokenIssueMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetInvolvedAddresses part of Msg interface
func (msg MiniTokenIssueMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}

func validateIssueMsgMiniTokenSymbol(symbol string) error {
	if len(symbol) == 0 {
		return errors.New("token symbol cannot be empty")
	}

	// check len without suffix
	if symbolLen := len(symbol); symbolLen > MiniTokenSymbolMaxLen || symbolLen < MiniTokenSymbolMinLen {
		return errors.New("length of token symbol is limited to 2~8")
	}

	if !common.IsAlphaNum(symbol) {
		return errors.New("token symbol should be alphanumeric")
	}

	return nil
}
