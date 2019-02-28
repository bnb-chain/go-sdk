package websocket

import (
	"encoding/json"

	"github.com/binance-chain/go-sdk/common/types"
)

//
type AccountEvent struct {
	EventType string              `json:"e"` // "e": "outboundAccountInfo"
	EventTime int64               `json:"E"` // "E": 1499405658658,
	Balances  []EventAssetBalance `json:"B"`
}

// EventAssetBalance record structure as send to the user
type EventAssetBalance struct {
	Asset  string       `json:"a"` // "a": "LTC",               // Asset
	Free   types.Fixed8 `json:"f"` // "f": "17366.18538083",    // Free amount
	Frozen types.Fixed8 `json:"r"` // "r":  "0.00000000"        // Frozen amount
	Locked types.Fixed8 `json:"l"` // "l":  "0.00000000"        // Locked amount
}

func (c *client) SubscribeAccountEvent(userAddr string, quit chan struct{}, onReceive func(event *AccountEvent), onError func(err error), onClose func()) error {
	msgs, err := c.baseClient.WsGet(userAddr, func(bz []byte) (interface{}, error) {
		var event AccountEvent
		err := json.Unmarshal(bz, &event)
		// Todo: the ws will return order data also. Ignore error now
		if err != nil {
			return nil, nil
		}
		return &event, err
	}, quit)
	if err != nil {
		return err
	}
	go c.SubscribeEvent(quit, msgs, func(event interface{}) {
		if accountEvent, ok := event.(*AccountEvent); ok {
			onReceive(accountEvent)
		}
	}, onError, onClose)
	return nil
}
