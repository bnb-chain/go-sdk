package query

import (
	"encoding/json"
	"fmt"

	"github.com/binance-chain/go-sdk/common"
	"github.com/binance-chain/go-sdk/common/types"
)

// GetKlines returns transaction details
func (c *client) GetKlines(query *types.KlineQuery) ([]types.Kline, error) {
	err := query.Check()
	if err != nil {
		return nil, err
	}
	qp, err := common.QueryParamToMap(*query)
	if err != nil {
		return nil, err
	}

	resp, _, err := c.baseClient.Get("/klines", qp)
	if err != nil {
		return nil, err
	}

	iklines := [][]interface{}{}
	if err := json.Unmarshal(resp, &iklines); err != nil {
		return nil, err
	}
	klines := make([]types.Kline, len(iklines))
	// Todo
	for index, ikline := range iklines {
		kl := types.Kline{}
		imap := make(map[string]interface{}, 9)
		if len(ikline) >= 9 {
			imap["openTime"] = ikline[0]
			imap["open"] = ikline[1]
			imap["high"] = ikline[2]
			imap["low"] = ikline[3]
			imap["close"] = ikline[4]
			imap["volume"] = ikline[5]
			imap["closeTime"] = ikline[6]
			imap["quoteAssetVolume"] = ikline[7]
			imap["NumberOfTrades"] = ikline[8]
		} else {
			return nil, fmt.Errorf("Receive kline scheme is unexpected ")
		}
		bz, err := json.Marshal(imap)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bz, &kl)
		if err != nil {
			return nil, err
		}
		klines[index] = kl
	}
	return klines, nil
}
