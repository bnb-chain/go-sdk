package websocket

import "encoding/json"

type BlockHeightEvent struct {
	BlockHeight int64 `json:"h"` // "h": 1499405658658,
}

func (c *client) SubscribeBlockHeightEvent(quit chan struct{}, onReceive func(event *BlockHeightEvent), onError func(err error), onClose func()) error {
	msgs, err := c.baseClient.WsGet("$all@blockheight", func(bz []byte) (interface{}, error) {
		var event BlockHeightEvent
		err := json.Unmarshal(bz, &event)
		return &event, err
	}, quit)
	if err != nil {
		return err
	}
	go c.SubscribeEvent(quit, msgs, func(event interface{}) {
		if blockHeightEvent, ok := event.(*BlockHeightEvent); ok {
			onReceive(blockHeightEvent)
		}
	}, onError, onClose)
	return nil
}
