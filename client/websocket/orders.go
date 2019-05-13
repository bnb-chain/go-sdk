package websocket

import (
	"encoding/json"

	"github.com/binance-chain/go-sdk/common/types"
)

type OrderEvent struct {
	EventType            string       `json:"e"` // "e": "executionReport"
	EventTime            int64        `json:"E"` // "E": 1499405658658,
	Symbol               string       `json:"s"` // "s": "ETH_BTC",
	Side                 int8         `json:"S"` //"S": "BUY",
	OrderType            int8         `json:"o"` //"o": "LIMIT", always to `LIMIT`, `will be published`
	TimeInForce          int8         `json:"f"` //"f": "GTC",
	OrderQty             types.Fixed8 `json:"q"` //"q": "1.00000000", `will be published`
	OrderPrice           types.Fixed8 `json:"p"` //"p": "0.10264410", `will be published`
	CurrentExecutionType string       `json:"x"` //"x": "NEW", always `NEW` for now `will be published`
	CurrentOrderStatus   string       `json:"X"` //"X": "Ack", "Canceled", "Expired", "IocNoFill", "PartialFill", "FullyFill", "FailedBlocking", "FailedMatching", "Unknown"
	OrderID              string       `json:"i"` //"i": "917E1846D6B3C40B97465CCF52818471E2C1027C-466",
	LastExecutedQty      types.Fixed8 `json:"l"` // "l": "0.00000000",
	LastExecutedPrice    types.Fixed8 `json:"L"` // "L": "0.00000000",
	CommulativeFilledQty types.Fixed8 `json:"z"` // "z": "0.00000000",
	CommissionAmount     string       `json:"n"` // "n": "0",
	TransactionTime      int64        `json:"T"` //"T": 1499405658657, `this will be the BLockheight`
	TradeID              string       `json:"t"` //"t": -1, `will be published`
	OrderCreationTime    int64        `json:"O"` //"O": 1499405658657, `will be published`
}

func (c *client) SubscribeOrderEvent(userAddr string, quit chan struct{}, onReceive func(event []*OrderEvent), onError func(err error), onClose func()) error {
	msgs, err := c.baseClient.WsGet(userAddr, func(bz []byte) (interface{}, error) {
		events := make([]*OrderEvent, 0)
		err := json.Unmarshal(bz, &events)
		// Todo: the ws will return account data also. Ignore error now
		if err != nil {
			return nil, nil
		}
		return events, err
	}, quit)
	if err != nil {
		return err
	}
	go c.SubscribeEvent(quit, msgs, func(event interface{}) {
		if orderEvent, ok := event.([]*OrderEvent); ok {
			onReceive(orderEvent)
		}
	}, onError, onClose)
	return nil
}
