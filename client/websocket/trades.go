package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/common/types"
)

type TradeEvent struct {
	EventType     string       `json:"e"`  // "e": "trade",
	EventTime     int64        `json:"E"`  // "E": 123456789,
	Symbol        string       `json:"s"`  // "s": "ETH_BTC",
	TradeID       string       `json:"t"`  // "t": 4293153,
	Price         types.Fixed8 `json:"p"`  // "p": "0.001",
	Qty           types.Fixed8 `json:"q"`  // "q": "88",
	BuyerOrderID  string       `json:"b"`  // "f": "50",
	SellerOrderID string       `json:"a"`  // "q": "100",
	TradeTime     int64        `json:"T"`  // "p": 123456785
	SellerAddress string       `json:"sa"` // "sa": 0x4092678e4e78230f46a1534c0fbc8fa39780892b
	BuyerAddress  string       `json:"ba"` // "ba": 0x4092778e4e78230f46a1534c0fbc8fa39780892c
}

func (c *client) SubscribeTradeEvent(baseAssetSymbol, quoteAssetSymbol string, quit chan struct{}, onReceive func(events []*TradeEvent), onError func(err error), onClose func()) error {
	msgs, err := c.baseClient.WsGet(fmt.Sprintf("%s@%s", common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol), "trades"), func(bz []byte) (interface{}, error) {
		events := make([]*TradeEvent, 0)
		err := json.Unmarshal(bz, &events)
		return events, err
	}, quit)
	if err != nil {
		return err
	}
	go c.SubscribeEvent(quit, msgs, func(event interface{}) {
		if tradeEvents, ok := event.([]*TradeEvent); ok {
			onReceive(tradeEvents)
		}
	}, onError, onClose)
	return nil
}
