package query

import (
	"encoding/json"

	"github.com/binance-chain/go-sdk/common"
)

type MarketsQuery struct {
	Offset *uint32 `json:"offset,omitempty,string"` //Option
	Limit  *uint32 `json:"limit,omitempty,string"`  //Option
}

func NewMarketsQuery() *MarketsQuery {
	return &MarketsQuery{}
}

func (param *MarketsQuery) WithOffset(offset uint32) *MarketsQuery {
	param.Offset = &offset
	return param
}

func (param *MarketsQuery) WithLimit(limit uint32) *MarketsQuery {
	param.Limit = &limit
	return param
}

func (param *MarketsQuery) Check() error {
	if param.Limit != nil && *param.Limit <= 0 {
		return LimitOutOfRangeError
	}
	return nil
}

// SymbolPair definition
type SymbolPair struct {
	TradeAsset string `json:"base_asset_symbol"`
	QuoteAsset string `json:"quote_asset_symbol"`
	Price      string `json:"price"`
	TickSize   string `json:"tick_size"`
	LotSize    string `json:"lot_size"`
}

// GetMarkets returns list of trading pairs
func (c *client) GetMarkets(query *MarketsQuery) ([]SymbolPair, error) {
	err := query.Check()
	if err != nil {
		return nil, err
	}
	qp, err := common.QueryParamToMap(*query)
	if err != nil {
		return nil, err
	}
	resp, err := c.baseClient.Get("/markets", qp)
	if err != nil {
		return nil, err
	}
	var listOfPairs []SymbolPair
	if err := json.Unmarshal(resp, &listOfPairs); err != nil {
		return nil, err
	}

	return listOfPairs, nil
}
