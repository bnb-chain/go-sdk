package query

import (
	"github.com/binance-chain/go-sdk/client/basic"
	"github.com/binance-chain/go-sdk/common/types"
)

type QueryClient interface {
	GetClosedOrders(query *ClosedOrdersQuery) (*CloseOrders, error)
	GetDepth(query *DepthQuery) (*MarketDepth, error)
	GetKlines(query *KlineQuery) ([]Kline, error)
	GetMarkets(query *MarketsQuery) ([]types.TradingPair, error)
	GetOrder(orderID string) (*Order, error)
	GetOpenOrders(query *OpenOrdersQuery) (*OpenOrders, error)
	GetTicker24h(query *Ticker24hQuery) ([]Ticker24h, error)
	GetTrades(query *TradesQuery) (*Trades, error)
	GetAccount(string) (*Account, error)
	GetTime() (*Time, error)
	GetTokens() ([]types.Token, error)
	GetNodeInfo() (*ResultStatus, error)
}

type client struct {
	baseClient basic.BasicClient
}

func NewClient(c basic.BasicClient) QueryClient {
	return &client{baseClient: c}
}
