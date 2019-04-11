package msg

import (
	"encoding/json"
	"fmt"

	"github.com/binance-chain/go-sdk/common/types"
)

// DexListMsg def
type DexListMsg struct {
	From             types.AccAddress `json:"from"`
	ProposalId       int64            `json:"proposal_id"`
	BaseAssetSymbol  string           `json:"base_asset_symbol"`
	QuoteAssetSymbol string           `json:"quote_asset_symbol"`
	InitPrice        int64            `json:"init_price"`
}

// NewDexListMsg for instance creation
func NewDexListMsg(from types.AccAddress, proposalId int64, baseAssetSymbol string, quoteAssetSymbol string, initPrice int64) DexListMsg {
	return DexListMsg{
		From:             from,
		ProposalId:       proposalId,
		BaseAssetSymbol:  baseAssetSymbol,
		QuoteAssetSymbol: quoteAssetSymbol,
		InitPrice:        initPrice,
	}
}

// Route part of Msg interface
func (msg DexListMsg) Route() string { return "dexList" }

// Type part of Msg interface
func (msg DexListMsg) Type() string { return "dexList" }

// String part of Msg interface
func (msg DexListMsg) String() string { return fmt.Sprintf("MsgList{%#v}", msg) }

// GetSigners part of Msg interface
func (msg DexListMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

// GetSignBytes part of Msg interface
func (msg DexListMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// GetInvolvedAddresses part of Msg interface
func (msg DexListMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
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
