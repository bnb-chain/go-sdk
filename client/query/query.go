package query

import (
	"github.com/binance-chain/go-sdk/client/basic"
	"github.com/binance-chain/go-sdk/common/types"
)

type QueryClient interface {
	GetClosedOrders(query *types.ClosedOrdersQuery) (*types.CloseOrders, error)
	GetDepth(query *types.DepthQuery) (*types.MarketDepth, error)
	GetKlines(query *types.KlineQuery) ([]types.Kline, error)
	GetMarkets(query *types.MarketsQuery) ([]types.TradingPair, error)
	GetOrder(orderID string) (*types.Order, error)
	GetOpenOrders(query *types.OpenOrdersQuery) (*types.OpenOrders, error)
	GetTicker24h(query *types.Ticker24hQuery) ([]types.Ticker24h, error)
	GetTrades(query *types.TradesQuery) (*types.Trades, error)
	GetAccount(string) (*types.BalanceAccount, error)
	GetTime() (*types.Time, error)
	GetTokens(query *types.TokensQuery) ([]types.Token, error)
	GetNodeInfo() (*types.ResultStatus, error)
}

type client struct {
	baseClient basic.BasicClient
}

func NewClient(c basic.BasicClient) QueryClient {
	return &client{baseClient: c}
}
