package sdk

type fDexAPI struct{}

func (api *fDexAPI) Get(path string, qp map[string]string) ([]byte, error) {

	if path == "/pairs" {

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

	return nil, nil
}

func (api *fDexAPI) Post(path string, qp map[string]string, body []byte) ([]byte, error) {
	return nil, nil
}
