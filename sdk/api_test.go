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
      "coins": [
        {
          "denom": "BNB",
          "amount": "18975020177895000"
        },
        {
          "denom": "NNB",
          "amount": "1737120240518526"
        },
        {
          "denom": "ZCB",
          "amount": "1887578172962946"
        }
      ],
      "public_key": {
        "type": "tendermint/PubKeySecp256k1",
        "value": "A58TeSbC3MRQ1ig5heN/XPinu9kjZrK4gp60DD7czU8J"
      },
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

	return nil, nil
}

func (api *fDexAPI) Post(path string, qp map[string]string, body []byte) ([]byte, error) {
	return nil, nil
}
