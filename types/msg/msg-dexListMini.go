package msg

import (
	"encoding/json"
	"fmt"

	"github.com/binance-chain/go-sdk/common/types"
)

// ListMiniMsg def
type ListMiniMsg struct {
	From             types.AccAddress `json:"from"`
	BaseAssetSymbol  string           `json:"base_asset_symbol"`
	QuoteAssetSymbol string           `json:"quote_asset_symbol"`
	InitPrice        int64            `json:"init_price"`
}

// NewListMiniMsg for instance creation
func NewListMiniMsg(from types.AccAddress, baseAssetSymbol string, quoteAssetSymbol string, initPrice int64) ListMiniMsg {
	return ListMiniMsg{
		From:             from,
		BaseAssetSymbol:  baseAssetSymbol,
		QuoteAssetSymbol: quoteAssetSymbol,
		InitPrice:        initPrice,
	}
}

// Route part of Msg interface
func (msg ListMiniMsg) Route() string { return "dexListMini" }

// Type part of Msg interface
func (msg ListMiniMsg) Type() string { return "dexListMini" }

// String part of Msg interface
func (msg ListMiniMsg) String() string { return fmt.Sprintf("MsgListMini{%#v}", msg) }

// GetSigners part of Msg interface
func (msg ListMiniMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

// GetSignBytes part of Msg interface
func (msg ListMiniMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetInvolvedAddresses part of Msg interface
func (msg ListMiniMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}

// ValidateBasic part of Msg interface
func (msg ListMiniMsg) ValidateBasic() error {
	err := ValidateMiniTokenSymbol(msg.BaseAssetSymbol)
	if err != nil {
		return fmt.Errorf("Invalid base asset token %v", msg.BaseAssetSymbol)
	}

	if msg.InitPrice <= 0 {
		return fmt.Errorf("Price should be positive")
	}

	return nil
}
