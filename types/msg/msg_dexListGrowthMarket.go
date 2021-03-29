package msg

import (
	"encoding/json"
	"fmt"

	"github.com/binance-chain/go-sdk/common/types"
)

const ListGrowthMarketMsgType = "dexListGrowthMarket"

type ListGrowthMarketMsg struct {
	From             types.AccAddress `json:"from"`
	BaseAssetSymbol  string           `json:"base_asset_symbol"`
	QuoteAssetSymbol string           `json:"quote_asset_symbol"`
	InitPrice        int64            `json:"init_price"`
}

func NewListGrowthMarketMsg(from types.AccAddress, baseAssetSymbol string, quoteAssetSymbol string, initPrice int64) ListGrowthMarketMsg {
	return ListGrowthMarketMsg{
		From:             from,
		BaseAssetSymbol:  baseAssetSymbol,
		QuoteAssetSymbol: quoteAssetSymbol,
		InitPrice:        initPrice,
	}
}

func (msg ListGrowthMarketMsg) Route() string                  { return "dexList" }
func (msg ListGrowthMarketMsg) Type() string                   { return ListGrowthMarketMsgType }
func (msg ListGrowthMarketMsg) String() string                 { return fmt.Sprintf("MsgListGrowthMarket{%#v}", msg) }
func (msg ListGrowthMarketMsg) GetSigners() []types.AccAddress { return []types.AccAddress{msg.From} }

func (msg ListGrowthMarketMsg) ValidateBasic() error {
	if msg.BaseAssetSymbol == msg.QuoteAssetSymbol {
		return fmt.Errorf("base token and quote token should not be the same")
	}

	if !IsValidMiniTokenSymbol(msg.BaseAssetSymbol) {
		err := ValidateSymbol(msg.BaseAssetSymbol)
		if err != nil {
			return fmt.Errorf("Invalid base asset token %v", msg.BaseAssetSymbol)
		}
	}

	if !IsValidMiniTokenSymbol(msg.QuoteAssetSymbol) {
		err := ValidateSymbol(msg.QuoteAssetSymbol)
		if err != nil {
			return fmt.Errorf("Invalid quote asset token %v", msg.QuoteAssetSymbol)
		}
	}

	if msg.InitPrice <= 0 {
		return fmt.Errorf("price should be positive")
	}
	return nil
}

func (msg ListGrowthMarketMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ListGrowthMarketMsg) GetInvolvedAddresses() []types.AccAddress {
	return msg.GetSigners()
}
