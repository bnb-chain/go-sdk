package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/common/types"
)

type TickerEvent struct {
	EventType          string       `json:"e"` // "e": "24hrTicker",
	EventTime          int64        `json:"E"` // "E": 123456789,
	Symbol             string       `json:"s"` // "s": "BNBBTC",      // Symbol
	PriceChange        types.Fixed8 `json:"p"` // "p": "0.0015",      // Price change
	PriceChangePercent types.Fixed8 `json:"P"` // "P": "250.00",      // Price change percent
	WeightedAvgPrice   types.Fixed8 `json:"w"` // "w": "0.0018",      // Weighted average price
	PrevClosePrice     types.Fixed8 `json:"x"` // "x": "0.0009",      // Previous day's close price
	LastPrice          types.Fixed8 `json:"c"` // "c": "0.0025",      // Current day's close price
	LastQuantity       types.Fixed8 `json:"Q"` // "Q": "10",          // Close trade's quantity
	BidPrice           types.Fixed8 `json:"b"` // "b": "0.0024",      // Best bid price
	BidQuantity        types.Fixed8 `json:"B"` // "B": "10",          // Best bid quantity
	AskPrice           types.Fixed8 `json:"a"` // "a": "0.0026",      // Best ask price
	AskQuantity        types.Fixed8 `json:"A"` // "A": "100",         // Best ask quantity
	OpenPrice          types.Fixed8 `json:"o"` // "o": "0.0010",      // Open price
	HighPrice          types.Fixed8 `json:"h"` // "h": "0.0025",      // High price
	LowPrice           types.Fixed8 `json:"l"` // "l": "0.0010",      // Low price
	Volume             types.Double `json:"v"` // "v": "10000",       // Total traded base asset volume
	QuoteVolume        types.Double `json:"q"` // "q": "18",          // Total traded quote asset volume
	OpenTime           int64        `json:"O"` // "O": 0,             // Statistics open time
	CloseTime          int64        `json:"C"` // "C": 86400000,      // Statistics close time
	FirstID            string       `json:"F"` // "F": 0,             // First trade ID
	LastID             string       `json:"L"` // "L": 18150,         // Last trade Id
	Count              int64        `json:"n"` // "n": 18151          // Total number of trades
}

func (c *client) SubscribeTickerEvent(baseAssetSymbol, quoteAssetSymbol string, quit chan struct{}, onReceive func(event *TickerEvent), onError func(err error), onClose func()) error {
	msgs, err := c.baseClient.WsGet(fmt.Sprintf("%s@%s", common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol), "ticker"), func(bz []byte) (interface{}, error) {
		event := TickerEvent{}
		err := json.Unmarshal(bz, &event)
		return &event, err
	}, quit)
	if err != nil {
		return err
	}
	go c.SubscribeEvent(quit, msgs, func(event interface{}) {
		if tickerEvent, ok := event.(*TickerEvent); ok {
			onReceive(tickerEvent)
		}
	}, onError, onClose)
	return nil
}

func (c *client) SubscribeAllTickerEvent(quit chan struct{}, onReceiveHandler func(event []*TickerEvent), onError func(err error), onClose func()) error {
	msgs, err := c.baseClient.WsGet(fmt.Sprintf("%s@%s", "$all", "allTickers"), func(bz []byte) (interface{}, error) {
		events := make([]*TickerEvent, 0)
		err := json.Unmarshal(bz, &events)
		return events, err
	}, quit)
	if err != nil {
		return err
	}
	go c.SubscribeEvent(quit, msgs, func(event interface{}) {
		if tickersEvent, ok := event.([]*TickerEvent); ok {
			onReceiveHandler(tickersEvent)
		}
	}, onError, onClose)
	return nil
}

type MiniTickerEvent struct {
	EventType   string       `json:"e"` // "e": "24hrMiniTicker",
	EventTime   int64        `json:"E"` // "E": 123456789,
	Symbol      string       `json:"s"` // "s": "BNBBTC",      // Symbol
	LastPrice   types.Fixed8 `json:"c"` // "c": "0.0025",      // Current day's close price
	OpenPrice   types.Fixed8 `json:"o"` // "o": "0.0010",      // Open price
	HighPrice   types.Fixed8 `json:"h"` // "h": "0.0025",      // High price
	LowPrice    types.Fixed8 `json:"l"` // "l": "0.0010",      // Low price
	Volume      types.Double `json:"v"` // "v": "10000",       // Total traded base asset volume
	QuoteVolume types.Double `json:"q"` // "q": "18",          // Total traded quote asset volume
}

func (c *client) SubscribeMiniTickerEvent(baseAssetSymbol, quoteAssetSymbol string, quit chan struct{}, onReceive func(event *MiniTickerEvent), onError func(err error), onClose func()) error {
	msgs, err := c.baseClient.WsGet(fmt.Sprintf("%s@%s", common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol), "miniTicker"), func(bz []byte) (interface{}, error) {
		var event MiniTickerEvent
		err := json.Unmarshal(bz, &event)
		return &event, err
	}, quit)
	if err != nil {
		return err
	}
	go c.SubscribeEvent(quit, msgs, func(event interface{}) {
		if miniTickerEvent, ok := event.(*MiniTickerEvent); ok {
			onReceive(miniTickerEvent)
		}
	}, onError, onClose)
	return nil
}

func (c *client) SubscribeAllMiniTickersEvent(quit chan struct{}, onReceive func(events []*MiniTickerEvent), onError func(err error), onClose func()) error {
	msgs, err := c.baseClient.WsGet(fmt.Sprintf("%s@%s", "$all", "allMiniTickers"), func(bz []byte) (interface{}, error) {
		events := make([]*MiniTickerEvent, 0)
		err := json.Unmarshal(bz, &events)
		return events, err
	}, quit)
	if err != nil {
		return err
	}
	go c.SubscribeEvent(quit, msgs, func(event interface{}) {
		if miniTickers, ok := event.([]*MiniTickerEvent); ok {
			onReceive(miniTickers)
		}
	}, onError, onClose)
	return nil
}
