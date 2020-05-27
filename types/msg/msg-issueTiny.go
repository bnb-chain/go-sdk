package msg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/binance-chain/go-sdk/common/types"
)

const (
	IssueTinyMsgType = "tinyIssueMsg"
)

// TinyTokenIssueMsg def
type TinyTokenIssueMsg struct {
	From        types.AccAddress `json:"from"`
	Name        string           `json:"name"`
	Symbol      string           `json:"symbol"`
	TotalSupply int64            `json:"total_supply"`
	Mintable    bool             `json:"mintable"`
	TokenURI    string           `json:"token_uri"`
}

// NewTinyTokenIssueMsg for instance creation
func NewTinyTokenIssueMsg(from types.AccAddress, name, symbol string, supply int64, mintable bool, tokenURI string) TinyTokenIssueMsg {
	return TinyTokenIssueMsg{
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
func (msg TinyTokenIssueMsg) ValidateBasic() error {

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
func (msg TinyTokenIssueMsg) Route() string { return MiniRoute }

// Type part of Msg interface
func (msg TinyTokenIssueMsg) Type() string {
	return IssueTinyMsgType
}

// String part of Msg interface
func (msg TinyTokenIssueMsg) String() string { return fmt.Sprintf("IssueTinyMsg{%#v}", msg) }

// GetSigners part of Msg interface
func (msg TinyTokenIssueMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

// GetSignBytes part of Msg interface
func (msg TinyTokenIssueMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetInvolvedAddresses part of Msg interface
func (msg TinyTokenIssueMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}
