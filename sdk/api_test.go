package sdk

type fDexAPI struct{}

func (api *fDexAPI) Get(path string, qp map[string]string) ([]byte, error) {

	// fmt.Println("qp: ", qp)

	if path == "/markets" {
		return []byte(`[
		{
			"base_asset_symbol": "BNB",
			"quote_asset_symbol": "NNB",
			"price": "100000000",
			"tick_size": "0.00005000",
			"lot_size": "0.10000000"
		}
	]`), nil
	}

	if path == "/depth" {
		return []byte(`{
			"lastUpdateId": 1000,
			"symbol": "BNB_NNB",
			"bids": [["0.00240000","50"],["0.00230000","100"]],
			"asks": [["0.00250000","90"],["0.00260000","120"],["0.0030000","350"]]
		}`), nil
	}

	if path == "/account/cosmosaccaddr1hy2e872rqtd675sn72ny87cyyaaanmqeuvwrpc" {
		return []byte(`{
      "address": "cosmosaccaddr1hy2e872rqtd675sn72ny87cyyaaanmqeuvwrpc",
      "balances": [
        {
          "symbol": "BNB",
					"free": "18975020177895000",
					"locked": "000000000",
					"frozen": "000000000"
        },
        {
          "symbol": "NNB",
					"free": "828912912928291",
					"locked": "000000000",
					"frozen": "75000000"
				}
      ],
      "public_key": [1,0,1,1,1,1,1,1,0,0,0,0,0,0,0,0,1],
      "account_number": "0",
      "sequence": "298113"
    }`), nil
	}

	if path == "/tx/52ECED0360605C1F3F336CA20B2C60535B0C72F0" {
		return []byte(`{
			"hash": "52ECED0360605C1F3F336CA20B2C60535B0C72F0",
			"log": "Msg 0: ",
			"data": "eyJ0eXBlIjoiZGV4L05ld09yZGVyUmVzcG9uc2UiLCJ2YWx1ZSI6eyJvcmRlcl9pZCI6ImNvc21vc2FjY2FkZHIxcTY4cGhxN3E2Znl1cDV4MjVtYWdsZjlzeGMydDRoeTQycGE2MjMtMjQwNDAyIn19",
			"code": 0
		}`), nil
	}

	if path == "/orders/cosmosaccaddr1q68phq7q6fyup5x25maglf9sxc2t4hy42pa623" {
		return []byte(`{
			"orderId": "cosmosaccaddr1q68phq7q6fyup5x25maglf9sxc2t4hy42pa623-240402",
			"owner": "cosmosaccaddr1q68phq7q6fyup5x25maglf9sxc2t4hy42pa623",
			"symbol": "NNB_BNB",
			"price": "1600000000",
			"quantity": "8900000000",
			"executedQuantity": "8900000000",
			"side": "SELL",
			"status": "FULLY_FILLED",
			"timeinforce": "GTC",
			"type": "LIMIT"
		}`), nil
	}

	if path == "/orders/open" {
		return []byte(`[{
			"orderId": "cosmosaccaddr1q68phq7q6fyup5x25maglf9sxc2t4hy42pa623-240402",
			"owner": "cosmosaccaddr1q68phq7q6fyup5x25maglf9sxc2t4hy42pa623",
			"symbol": "NNB_BNB",
			"price": "1600000000",
			"quantity": "8900000000",
			"executedQuantity": "8900000000",
			"side": "SELL",
			"status": "FULLY_FILLED",
			"timeinforce": "GTC",
			"type": "LIMIT"
		}]`), nil
	}

	if path == "/orders/closed" {
		return []byte(`[{
			"orderId": "cosmosaccaddr1q68phq7q6fyup5x25maglf9sxc2t4hy42pa623-240402",
			"owner": "cosmosaccaddr1q68phq7q6fyup5x25maglf9sxc2t4hy42pa623",
			"symbol": "NNB_BNB",
			"price": "1600000000",
			"quantity": "8900000000",
			"executedQuantity": "8900000000",
			"side": "SELL",
			"status": "FULLY_FILLED",
			"timeinforce": "GTC",
			"type": "LIMIT"
		}]`), nil
	}

	if path == "/trades" {
		return []byte(`[{
			"buyerOrderId": "order-buy-1",
			"buyFee": "0.50000000",
			"price": "0.75000000",
			"quantity": "1.00000000",
			"sellerOrderId": "order-sell-1",
			"sellFee": "0.50000000",
			"symbol": "BNB_NNB",
			"time": 1000000000,
			"tradeId": "trade-1"
		}]`), nil
	}

	if path == "/klines" {
		return []byte(`[{
			"close": 50000000,       
			"closeTime": 90000000,
			"high": 150000000,
			"low": 500000000,
			"numberOfTrades": 150,
			"open": 500000000,
			"openTime": 1000000,
			"quoteAssetVolume": 800000000,
			"volume": 2000000000
		}]`), nil
	}

	if path == "/ticker/24hr" {
		return []byte(`{
			"symbol":"BNB",
			"askPrice":"1.0000000",
			"askQuantity":"1.0000000",
			"bidPrice":"1.0000000",
			"bidQuantity":"1.0000000",
			"closeTime":10000000,
			"count":100,
			"firstId":"order-1",
			"highPrice":"1.0000000",
			"lastId":"order-100",
			"lastPrice":"1.0000000",
			"lastQuantity":"1.0000000",
			"lowPrice":"1.0000000",
			"openPrice":"1.0000000",
			"openTime":2000000,
			"prevClosePrice":"1.0000000",
			"priceChange":"1.0000000",
			"priceChangePercent":"15",
			"quoteVolume":"1.0000000",
			"volume":"1.0000000",
			"weightedAvgPrice":"1.0000000"
	 }`), nil
	}

	if path == "/ticker/ticker" {
		return []byte(`{
			"symbol":"BNB",
			"askPrice":"1.0000000",
			"askQty":"1.0000000",
			"bidPrice":"1.0000000",
			"bidQty":"1.0000000"
	 }`), nil
	}

	if path == "/time" {
		return []byte(`{
			"Time":"14:20:00T"
	 }`), nil
	}

	if path == "/tokens" {
		return []byte(`[
		{
			"name": "ABC Token",
			"symbol": "BNB",
			"total_supply": "100000000",
			"owner": "cosmosaccaddr1hy2e872rqtd675sn72ny87cyyaaanmqeuvwrpc"
		}
	]`), nil
	}

	return nil, nil
}

func (api *fDexAPI) Post(path string, qp map[string]string, body []byte) ([]byte, error) {
	return nil, nil
}
