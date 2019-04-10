package types

type OpenOrder struct {
	Id                   string       `json:"id"`
	Symbol               string       `json:"symbol"`
	Price                Fixed8 `json:"price"`
	Quantity             Fixed8 `json:"quantity"`
	CumQty               Fixed8 `json:"cumQty"`
	CreatedHeight        int64        `json:"createdHeight"`
	CreatedTimestamp     int64        `json:"createdTimestamp"`
	LastUpdatedHeight    int64        `json:"lastUpdatedHeight"`
	LastUpdatedTimestamp int64        `json:"lastUpdatedTimestamp"`
}

type TradingPair struct {
	BaseAssetSymbol  string        `json:"base_asset_symbol"`
	QuoteAssetSymbol string        `json:"quote_asset_symbol"`
	ListPrice        Fixed8 `json:"list_price"`
	TickSize         Fixed8 `json:"tick_size"`
	LotSize          Fixed8 `json:"lot_size"`
}


type OrderBook struct {
	Height int64
	Levels []OrderBookLevel
}

// OrderBookLevel represents a single order book level.
type OrderBookLevel struct {
	BuyQty    Fixed8 `json:"buyQty"`
	BuyPrice  Fixed8 `json:"buyPrice"`
	SellQty   Fixed8 `json:"sellQty"`
	SellPrice Fixed8 `json:"sellPrice"`
}
