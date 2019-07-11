package types

// OrderSide enum
var OrderSide = struct {
	BUY  string
	SELL string
}{
	"BUY",
	"SELL",
}

// TimeInForce enum
var TimeInForce = struct {
	GTC string
	IOC string
}{"GTC", "IOC"}

// OrderStatus enum
var OrderStatus = struct {
	ACK              string
	PARTIALLY_FILLED string
	IOC_NO_FILL      string
	FULLY_FILLED     string
	CANCELED         string
	EXPIRED          string
	UNKNOWN          string
}{
	"ACK",
	"PARTIALLY_FILLED",
	"IOC_NO_FILL",
	"FULLY_FILLED",
	"CANCELED",
	"EXPIRED",
	"UNKNOWN",
}

// OrderType enum
var OrderType = struct {
	LIMIT             string
	MARKET            string
	STOP_LOSS         string
	STOP_LOSS_LIMIT   string
	TAKE_PROFIT       string
	TAKE_PROFIT_LIMIT string
	LIMIT_MAKER       string
}{
	"LIMIT",
	"MARKET",
	"STOP_LOSS",
	"STOP_LOSS_LIMIT",
	"TAKE_PROFIT",
	"TAKE_PROFIT_LIMIT",
	"LIMIT_MAKER",
}

type CloseOrders struct {
	Order []Order `json:"order"`
	Total int     `json:"total"`
}

type OpenOrders struct {
	Order []Order `json:"order"`
	Total int     `json:"total"`
}

type OpenOrder struct {
	Id                   string `json:"id"`
	Symbol               string `json:"symbol"`
	Price                Fixed8 `json:"price"`
	Quantity             Fixed8 `json:"quantity"`
	CumQty               Fixed8 `json:"cumQty"`
	CreatedHeight        int64  `json:"createdHeight"`
	CreatedTimestamp     int64  `json:"createdTimestamp"`
	LastUpdatedHeight    int64  `json:"lastUpdatedHeight"`
	LastUpdatedTimestamp int64  `json:"lastUpdatedTimestamp"`
}

type TradingPair struct {
	BaseAssetSymbol  string `json:"base_asset_symbol"`
	QuoteAssetSymbol string `json:"quote_asset_symbol"`
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

type Trades struct {
	Trade []Trade `json:"trade"`
	Total int     `json:"total"`
}

// Trade def
type Trade struct {
	BaseAsset     string `json:"baseAsset"`
	BlockHeight   int64  `json:"blockHeight"`
	BuyFee        string `json:"buyFee"`
	BuySingleFee  string `json:"buySingleFee"`
	BuyerId       string `json:"buyerId"`
	BuyerOrderID  string `json:"buyerOrderId"`
	Price         string `json:"price"`
	Quantity      string `json:"quantity"`
	QuoteAsset    string `json:"quoteAsset"`
	SellFee       string `json:"sellFee"`
	sellSingleFee string `json:"sellSingleFee"`
	SellerId      string `json:"sellerId"`
	SellerOrderID string `json:"sellerOrderId"`
	Symbol        string `json:"symbol"`
	Time          int64  `json:"time"`
	TradeID       string `json:"tradeId"`
	TickType      string `json:"tickType"`
}

type Order struct {
	ID                   string `json:"orderId"`
	Owner                string `json:"owner"`
	Symbol               string `json:"symbol"`
	Price                string `json:"price"`
	Quantity             string `json:"quantity"`
	CumulateQuantity     string `json:"cumulateQuantity"`
	Fee                  string `json:"fee"`
	Side                 int    `json:"side"` // BUY or SELL
	Status               string `json:"status"`
	TimeInForce          int    `json:"timeInForce"`
	Type                 int    `json:"type"`
	TradeId              string `json:"tradeId"`
	LastExecutedPrice    string `json:"last_executed_price"`
	LastExecutedQuantity string `json:"lastExecutedQuantity"`
	TransactionHash      string `json:"transactionHash"`
	OrderCreateTime      string `json:"orderCreateTime"`
	TransactionTime      string `json:"transactionTime"`
}
