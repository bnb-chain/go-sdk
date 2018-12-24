# BNC Chain Go SDK

## Description

Bnc-Go-SDK provides a thin wrapper around the BNC Chain API for readonly endpoints, in addition to creating and submitting different transactions.

## Usage

### Init

```GO
mnemonic := "lock globe panda armed mandate fabric couple dove climb step stove price recall decrease fire sail ring media enhance excite deny valid ceiling arm"
//-----   Init KeyManager  -------------
keyManager, _ := keys.NewMnemonicKeyManager(mnemonic)
//-----   Init sdk  -------------
bnc, _ := sdk.NewBncSDK("http://dex-api.fdgahl.cn", "chain-bnb", keyManager)

```
For sdk init, you should know the famous api address and chain id of bnbchain, it is "http://dex-api.fdgahl.cn" and "chain-bnb" in the above example.

If you want broadcast some transactions, like send coins, create orders or cancel orders, you should init a key manager to keep and use you private key.

There are three ways to get a key manager: from mnemonic, from key base file, from raw private key string.

From mnemonic:
```Go
mnemonic := "lock globe panda armed mandate fabric couple dove climb step stove price recall decrease fire sail ring media enhance excite deny valid ceiling arm"
keyManager, _ := keys.NewMnemonicKeyManager(mnemonic)
```

From key base file:
```GO
file := "testkeystore.json"
keyManager, err := NewKeyStoreKeyManager(file, "Zjubfd@123")

```

From raw private key string:
```GO
priv := "9579fff0cab07a4379e845a890105004ba4c8276f8ad9d22082b2acbf02d884b"
keyManager, err := NewPrivateKeyManager(priv)
```


### Read Operations

```GO
//----- Recover account ---------
mnemonic := "lock globe panda armed mandate fabric couple dove climb step stove price recall decrease fire sail ring media enhance excite deny valid ceiling arm"
keyManager, err := keys.NewMnemonicKeyManager(mnemonic)
assert.NoError(t, err)
testAccount1 := keyManager.GetAddr()
_, testAccount2 := PrivAndAddr()


//-----   Init sdk  -------------
bnc, _ := sdk.NewBncSDK("http://dex-api.fdgahl.cn", "chain-bnb", keyManager)
nativeSymbol := txmsg.NativeToken

//-----  Get account  -----------
account, err := bnc.GetAccount(testAccount1.String())
assert.NoError(t, err)
assert.NotNil(t, account)
assert.True(t, len(account.Balances) > 0)

//----- Get Markets  ------------
markets, err := bnc.GetMarkets(api.NewMarketsQuery().WithLimit(1))
assert.NoError(t, err)
assert.Equal(t, 1, len(markets))
tradeSymbol := markets[0].TradeAsset
if markets[0].QuoteAsset != "BNB" {
    tradeSymbol = markets[0].QuoteAsset
}

//-----  Get Depth  -----------
depth, err := bnc.GetDepth(api.NewDepthQuery(tradeSymbol, nativeSymbol))
assert.NoError(t, err)
assert.True(t, depth.Height > 0)

//----- Get Kline
kline, err := bnc.GetKlines(api.NewKlineQuery(tradeSymbol, nativeSymbol, "1h").WithLimit(1))
assert.NoError(t, err)
assert.Equal(t, 1, len(kline))

//-----  Get Ticker 24h  -----------
ticker24h, err := bnc.GetTicker24h(api.NewTicker24hQuery().WithSymbol(tradeSymbol, nativeSymbol))
assert.NoError(t, err)
assert.True(t, len(ticker24h) > 0)

//-----  Get Tokens  -----------
tokens, err := bnc.GetTokens()
assert.NoError(t, err)
fmt.Printf("GetTokens: %v \n", tokens)

//-----  Get Trades  -----------
fmt.Println(testAccount1.String())
trades, err := bnc.GetTrades(api.NewTradesQuery(testAccount1.String()).WithSymbol(tradeSymbol, nativeSymbol))
assert.NoError(t, err)
fmt.Printf("GetTrades: %v \n", trades)

//-----  Get Time    -----------
time, err := bnc.GetTime()
assert.NoError(t, err)
fmt.Printf("Get time: %v \n", time)

//---- Get Order    ------------
order, err := bnc.GetOrder("Your Order Id")
assert.NoError(t, err)

//---- Get Open Order ---------
openOrders, err := bnc.GetOpenOrders(api.NewOpenOrdersQuery(testAccount1.String()))
assert.NoError(t, err)
assert.True(t, len(openOrders.Order) > 0)

//---- Get Close Order---------
closedOrders, err := bnc.GetClosedOrders(api.NewClosedOrdersQuery(testAccount1.String()).WithSymbol(tradeSymbol, nativeSymbol))
assert.NoError(t, err)
assert.True(t, len(closedOrders.Order) > 0)
fmt.Printf("GetClosedOrders: %v \n", closedOrders)

//----    Get tx      ---------
tx, err := bnc.GetTx(openOrders.Order[0].TransactionHash)
assert.NoError(t, err)
```

For read option, each api will need a Query parameter. Each Query parameter we provide a construct function: `NewXXXQuery`.
We recommend you to use this construct function when you new a Query parameter since the construction only need required parameters,
and for the optional parameters, you can use `WithXXX` to add it.

### Create & Post Transaction

There is one most important point you should notice that we use int64 to represent a decimal.
The decimal length is fix 8, which means:
`100000000` is equal to `1`
`150000000` is equal to `1.5`
`1050000000` is equal to `10.5`

```GO
mnemonic := "lock globe panda armed mandate fabric couple dove climb step stove price recall decrease fire sail ring media enhance excite deny valid ceiling arm"
keyManager, err := keys.NewMnemonicKeyManager(mnemonic)
testAccount1 := keyManager.GetAddr()
_, testAccount2 := PrivAndAddr()

//-----   Init sdk  -------------
bnc, _ := sdk.NewBncSDK("http://dex-api.fdgahl.cn", "chain-bnb", keyManager)
nativeSymbol := txmsg.NativeToken

//----- Create order -----------
createOrderResult, err := bnc.CreateOrder(tradeSymbol, nativeSymbol, txmsg.OrderSide.BUY, 100000000, 100000000, true)
assert.NoError(t, err)
assert.True(t, true, createOrderResult.Ok)

//---- Cancle Order  ---------
cancleOrderResult, err := bnc.CancelOrder(tradeSymbol, nativeSymbol, createOrderResult.OrderId, createOrderResult.OrderId, true)
assert.NoError(t, err)
assert.True(t, cancleOrderResult.Ok)

//----   Send tx  -----------
send, err := bnc.SendToken(testAccount2, nativeSymbol, 10000000000, true)
assert.NoError(t, err)
assert.True(t, send.Ok)

//----   Freeze Token ---------
freeze, err := bnc.FreezeToken(nativeSymbol, 100000000, true)
assert.NoError(t, err)
assert.True(t, freeze.Ok)
fmt.Printf("freeze token: %v\n", freeze)

//----   Unfreeze Token ---------
unfreeze, err := bnc.UnfreezeToken(nativeSymbol, 100000000, true)
assert.NoError(t, err)
assert.True(t, unfreeze.Ok)

//----   issue token ---------
issue, err := bnc.IssueToken("SDK-Token", "sdk", 10000000000000000, true, false)
assert.NoError(t, err)

//---- Submit Proposal ------
listTradingProposal, err := bnc.SubmitListPairProposal("New trading pair", txmsg.ListTradingPairParams{issue.Symbol, nativeSymbol, 1000000000, "my trade", time2.Now().Add(1 * time2.Hour)}, 200000000000, true)
fmt.Println(err)
assert.NoError(t, err)
fmt.Printf("Submit list trading pair: %v\n", listTradingProposal)

//---- Vote Proposal  -------
time2.Sleep(10 * time2.Second)
vote, err := bnc.VoteProposal(listTradingProposal.ProposalId, txmsg.OptionYes, true)
assert.NoError(t, err)
```