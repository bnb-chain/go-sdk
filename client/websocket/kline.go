package websocket

import (
	"encoding/json"
	"fmt"

	"github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/common/types"
)

type KlineInterval string

const (
	OneMinuteInterval KlineInterval = "1m"

	FiveMinutesInterval   KlineInterval = "5m"
	FifteenMinuteInterval KlineInterval = "15m"
	ThirtMinutesInterval  KlineInterval = "30m"
	OneHourInterval       KlineInterval = "1h"
	TwoHoursInterval      KlineInterval = "2h"
	FourHoursInterval     KlineInterval = "4h"
	SixHoursInterval      KlineInterval = "6h"
	EightHoursInterval    KlineInterval = "8h"
	TwelveHoursInterval   KlineInterval = "12h"
	OneDayInterval        KlineInterval = "1d"
	ThreeDaysInterval     KlineInterval = "3d"
	OneWeekInterval       KlineInterval = "1w"
	OneMonthInterval      KlineInterval = "1M"
)

type KlineEvent struct {
	EventType string           `json:"e"` // "e": "executionReport"
	EventTime int64            `json:"E"` // "E": 1499405658658,
	Symbol    string           `json:"s"` // "s": "ETH_BTC",
	Kline     KlineRecordEvent `json:"k"`
}

// KlineRecordEvent record structure as received from the kafka messages stream
type KlineRecordEvent struct {
	Timestamp        int64        `json:"-"`
	Symbol           string       `json:"s"` //"BNBBTC",  // Symbol
	OpenTime         int64        `json:"t"` //123400000, // Kline start time
	CloseTime        int64        `json:"T"` //123460000, // Kline close time
	Interval         string       `json:"i"` //"1m",      // Interval
	FirstTradeID     string       `json:"f"` //100,       // First trade ID
	LastTradeID      string       `json:"L"` //200,       // Last trade ID
	OpenPrice        types.Fixed8 `json:"o"` //"0.0010",  // Open price
	ClosePrice       types.Fixed8 `json:"c"` //"0.0020",  // Close price
	HighPrice        types.Fixed8 `json:"h"` //"0.0025",  // High price
	LowPrice         types.Fixed8 `json:"l"` //"0.0015",  // Low price
	Volume           types.Double `json:"v"` //"1000",    // Base asset volume
	QuoteAssetVolume types.Double `json:"q"` //"1.0000",  // Quote asset volume
	NumberOfTrades   int64        `json:"n"` //100,       // Number of trades
	Closed           bool         `json:"x"` //Is this kline closed?
}

func (c *client) SubscribeKlineEvent(baseAssetSymbol, quoteAssetSymbol string, interval KlineInterval, quit chan struct{}, onReceive func(event *KlineEvent), onError func(err error), onClose func()) error {
	msgs, err := c.baseClient.WsGet(fmt.Sprintf("%s@%s_%s", common.CombineSymbol(baseAssetSymbol, quoteAssetSymbol), "kline", interval), func(bz []byte) (interface{}, error) {
		var event KlineEvent
		err := json.Unmarshal(bz, &event)
		return &event, err
	}, quit)
	if err != nil {
		return err
	}
	go c.SubscribeEvent(quit, msgs, func(event interface{}) {
		if klineEvent, ok := event.(*KlineEvent); ok {
			onReceive(klineEvent)
		}
	}, onError, onClose)
	return nil
}
