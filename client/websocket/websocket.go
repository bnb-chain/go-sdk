package websocket

import (
	"github.com/binance-chain/go-sdk/client/basic"
)

type WSClient interface {
	SubscribeAccountEvent(userAddr string, quit chan struct{}, onReceive func(event *AccountEvent), onError func(err error), onClose func()) error
	SubscribeBlockHeightEvent(quit chan struct{}, onReceive func(event *BlockHeightEvent), onError func(err error), onClose func()) error
	SubscribeKlineEvent(baseAssetSymbol, quoteAssetSymbol string, interval KlineInterval, quit chan struct{}, onReceive func(event *KlineEvent), onError func(err error), onClose func()) error
	SubscribeMarketDiffEvent(baseAssetSymbol, quoteAssetSymbol string, quit chan struct{}, onReceive func(event *MarketDeltaEvent), onError func(err error), onClose func()) error
	SubscribeMarketDepthEvent(baseAssetSymbol, quoteAssetSymbol string, quit chan struct{}, onReceive func(event *MarketDepthEvent), onError func(err error), onClose func()) error
	SubscribeOrderEvent(userAddr string, quit chan struct{}, onReceive func(event []*OrderEvent), onError func(err error), onClose func()) error
	SubscribeTickerEvent(baseAssetSymbol, quoteAssetSymbol string, quit chan struct{}, onReceive func(event *TickerEvent), onError func(err error), onClose func()) error
	SubscribeAllTickerEvent(quit chan struct{}, onReceive func(event []*TickerEvent), onError func(err error), onClose func()) error
	SubscribeMiniTickerEvent(baseAssetSymbol, quoteAssetSymbol string, quit chan struct{}, onReceive func(event *MiniTickerEvent), onError func(err error), onClose func()) error
	SubscribeAllMiniTickersEvent(quit chan struct{}, onReceive func(events []*MiniTickerEvent), onError func(err error), onClose func()) error
	SubscribeTradeEvent(baseAssetSymbol, quoteAssetSymbol string, quit chan struct{}, onReceive func(events []*TradeEvent), onError func(err error), onClose func()) error
}

type client struct {
	baseClient basic.BasicClient
}

func NewClient(c basic.BasicClient) WSClient {
	return &client{baseClient: c}
}

func (c *client) SubscribeEvent(quit chan struct{}, msgs <-chan interface{}, onReceive func(event interface{}), onError func(err error), onClose func()) {
	for {
		select {
		case <-quit:
			return
		case m, ok := <-msgs:
			if !ok {
				if onClose != nil {
					onClose()
				}
				return
			} else {
				switch o := m.(type) {
				case error:
					if onError != nil {
						onError(o)
					}
					return
				default:
					onReceive(o)
				}
			}
		}
	}
}
