package api

import (
	"encoding/json"
)

// DepthQuery
type DepthQuery struct {
	Symbol string  `json:"symbol"`
	Limit  *uint32 `json:"limit,omitempty,string"`
}

func NewDepthQuery(baseAssetSymbol, quoteAssetSymbol string) *DepthQuery {
	return &DepthQuery{Symbol: CombineSymbol(baseAssetSymbol, quoteAssetSymbol)}
}

func (param *DepthQuery) WithLimit(limit uint32) *DepthQuery {
	param.Limit = &limit
	return param
}

func (param *DepthQuery) Check() error {
	if param.Symbol == ""{
		return  SymbolMissingError
	}
	if param.Limit!=nil && *param.Limit<=0{
		return LimitOutOfRangeError
	}
	return nil
}


// MarketDepth broad caste to the user
type MarketDepth struct {
	Bids   [][]string `json:"bids"` // "bids": [ [ "0.0024", "10" ] ]
	Asks   [][]string `json:"asks"` // "asks": [ [ "0.0024", "10" ] ]
	Height int64      `json:"height"`
}

// GetDepth returns market depth records
func (dex *dexAPI) GetDepth(query *DepthQuery) (*MarketDepth, error) {
	err:=query.Check()
	if err!=nil{
		return nil,err
	}
	qp, err := QueryParamToMap(*query)
	if err != nil {
		return nil, err
	}
	resp, err := dex.Get("/depth", qp)
	if err != nil {
		return nil, err
	}

	var MarketDepth MarketDepth
	if err := json.Unmarshal(resp, &MarketDepth); err != nil {
		return nil, err
	}

	return &MarketDepth, nil
}
