package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/common/types"
)

type MarketDeltaEvent struct {
	EventType string           `json:"e"` // "e": "depthUpdate",
	EventTime int64            `json:"E"` // "E": 123456789,
	Symbol    string           `json:"s"` // "s": "ETH_BTC",
	Bids      [][]types.Fixed8 `json:"b"` // "b": [ [ "0.0024", "10" ] ]
	Asks      [][]types.Fixed8 `json:"a"` // "a": [ [ "0.0024", "10" ] ]
}

func (c *client) SubscribeMarketDiffEvent(baseAssetSymbol, quoteAssetSymbol string, quit chan struct{}, onReceive func(event *MarketDeltaEvent), onError func(err error), onClose func()) error {
	msgs, err := c.baseClient.WsGet(fmt.Sprintf("%s@%s", common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol), "marketDiff"), func(bz []byte) (interface{}, error) {
		var event MarketDeltaEvent
		err := json.Unmarshal(bz, &event)
		return &event, err
	}, quit)
	if err != nil {
		return err
	}
	go c.SubscribeEvent(quit, msgs, func(event interface{}) {
		if marketDeltaEvent, ok := event.(*MarketDeltaEvent); ok {
			onReceive(marketDeltaEvent)
		}
	}, onError, onClose)
	return nil
}

type MarketDepthEvent struct {
	LastUpdateID int64            `json:"lastUpdateId"` // "lastUpdateId": 160,
	Symbol       string           `json:"symbol"`       // "symbol": "BNB_BTC"
	Bids         [][]types.Fixed8 `json:"bids"`         // "bids": [ [ "0.0024", "10" ] ]
	Asks         [][]types.Fixed8 `json:"asks"`         // "asks": [ [ "0.0024", "10" ] ]
}

func (c *client) SubscribeMarketDepthEvent(baseAssetSymbol, quoteAssetSymbol string, quit chan struct{}, onReceive func(event *MarketDepthEvent), onError func(err error), onClose func()) error {
	msgs, err := c.baseClient.WsGet(fmt.Sprintf("%s@%s", common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol), "marketDepth"), func(bz []byte) (interface{}, error) {
		var event MarketDepthEvent
		err := json.Unmarshal(bz, &event)
		return &event, err
	}, quit)
	if err != nil {
		return err
	}
	go c.SubscribeEvent(quit, msgs, func(event interface{}) {
		if marketDepthEvent, ok := event.(*MarketDepthEvent); ok {
			onReceive(marketDepthEvent)
		}
	}, onError, onClose)
	return nil
}
