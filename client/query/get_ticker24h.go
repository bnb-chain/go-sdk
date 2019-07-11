package query

import (
	"encoding/json"
	"github.com/binance-chain/go-sdk/common/types"

	"github.com/binance-chain/go-sdk/common"
)

// GetTicker24h returns ticker 24h
func (c *client) GetTicker24h(query *types.Ticker24hQuery) ([]types.Ticker24h, error) {
	qp, err := common.QueryParamToMap(query)
	if err != nil {
		return nil, err
	}

	resp, _, err := c.baseClient.Get("/ticker/24hr", qp)
	if err != nil {
		return nil, err
	}

	tickers := []types.Ticker24h{}
	if err := json.Unmarshal(resp, &tickers); err != nil {
		return nil, err
	}

	return tickers, nil
}
