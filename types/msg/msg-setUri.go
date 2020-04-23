package msg

import (
	"encoding/json"
	"fmt"
	"github.com/binance-chain/go-sdk/common/types"
)

// SetURIMsg def
type SetURIMsg struct {
	From     types.AccAddress `json:"from"`
	Symbol   string           `json:"symbol"`
	TokenURI string           `json:"token_uri"`
}

// NewSetUriMsg for instance creation
func NewSetUriMsg(from types.AccAddress, symbol string, tokenURI string) SetURIMsg {
	return SetURIMsg{
		From:     from,
		Symbol:   symbol,
		TokenURI: tokenURI,
	}
}

// ValidateBasic does a simple validation check that
// doesn't require access to any other information.
func (msg SetURIMsg) ValidateBasic() error {
	if msg.From == nil {
		return fmt.Errorf("sender address cannot be empty")
	}

	if len(msg.TokenURI) > MaxTokenURILength {
		return fmt.Errorf("token seturi should not exceed %v characters", MaxTokenURILength)
	}

	err := ValidateMiniTokenSymbol(msg.Symbol)
	if err != nil {
		return fmt.Errorf("Invalid symbol %v", msg.Symbol)
	}

	return nil
}

// Route part of Msg interface
func (msg SetURIMsg) Route() string { return "miniTokensSetURI" }

// Type part of Msg interface
func (msg SetURIMsg) Type() string { return "miniTokensSetURI" }

// String part of Msg interface
func (msg SetURIMsg) String() string { return fmt.Sprintf("MsgSetURI{%#v}", msg) }

// GetSigners part of Msg interface
func (msg SetURIMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

// GetSignBytes part of Msg interface
func (msg SetURIMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetInvolvedAddresses part of Msg interface
func (msg SetURIMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}
