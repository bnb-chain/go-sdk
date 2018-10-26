package txmsg

import (
	"encoding/json"
	"fmt"
)

// DexListMsg def
type DexListMsg struct {
	From             AccAddress `json:"from"`
	BaseAssetSymbol  string     `json:"base_asset_symbol"`
	QuoteAssetSymbol string     `json:"quote_asset_symbol"`
	InitPrice        int64      `json:"init_price"`
}

// NewDexListMsg for instance creation
func NewDexListMsg(from AccAddress, baseAssetSymbol string, quoteAssetSymbol string, initPrice int64) DexListMsg {
	return DexListMsg{
		From:             from,
		BaseAssetSymbol:  baseAssetSymbol,
		QuoteAssetSymbol: quoteAssetSymbol,
		InitPrice:        initPrice,
	}
}

// Type part of Msg interface
func (msg DexListMsg) Type() string { return "dexList" }

// String part of Msg interface
func (msg DexListMsg) String() string { return fmt.Sprintf("MsgList{%#v}", msg) }

// Get part of Msg interface
func (msg DexListMsg) Get(key interface{}) (value interface{}) { return nil }

// GetSigners part of Msg interface
func (msg DexListMsg) GetSigners() []AccAddress { return []AccAddress{msg.From} }

// GetSignBytes part of Msg interface
func (msg DexListMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic part of Msg interface
func (msg DexListMsg) ValidateBasic() error {
	err := ValidateSymbol(msg.BaseAssetSymbol)
	if err != nil {
		return fmt.Errorf("Invalid base asset token %v", msg.BaseAssetSymbol)
	}

	err = ValidateSymbol(msg.QuoteAssetSymbol)
	if err != nil {
		return fmt.Errorf("Invalid quote asset token %v", msg.QuoteAssetSymbol)
	}

	if msg.InitPrice <= 0 {
		return fmt.Errorf("Price should be positive")
	}

	return nil
}
