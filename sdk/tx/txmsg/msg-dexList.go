package txmsg

import (
	"encoding/json"
	"fmt"
)

type DexListMsg struct {
	From             AccAddress `json:"from"`
	BaseAssetSymbol  string     `json:"base_asset_symbol"`
	QuoteAssetSymbol string     `json:"quote_asset_symbol"`
	InitPrice        int64      `json:"init_price"`
}

func NewDexListMsg(from AccAddress, baseAssetSymbol string, quoteAssetSymbol string, initPrice int64) DexListMsg {
	return DexListMsg{
		From:             from,
		BaseAssetSymbol:  baseAssetSymbol,
		QuoteAssetSymbol: quoteAssetSymbol,
		InitPrice:        initPrice,
	}
}

func (msg DexListMsg) Type() string                            { return "dexList" }
func (msg DexListMsg) String() string                          { return fmt.Sprintf("MsgList{%#v}", msg) }
func (msg DexListMsg) Get(key interface{}) (value interface{}) { return nil }
func (msg DexListMsg) GetSigners() []AccAddress                { return []AccAddress{msg.From} }

func (msg DexListMsg) ValidateBasic() error {
	// err := types.ValidateSymbol(msg.BaseAssetSymbol)
	// if err != nil {
	// 	return sdk.ErrInvalidCoins("base token: " + err.Error())
	// }

	// err = types.ValidateSymbol(msg.QuoteAssetSymbol)
	// if err != nil {
	// 	return sdk.ErrInvalidCoins("quote token: " + err.Error())
	// }

	// if msg.InitPrice <= 0 {
	// 	return sdk.ErrInvalidCoins("price should be positive")
	// }

	return nil
}

func (msg DexListMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}
