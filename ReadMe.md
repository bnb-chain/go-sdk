# BNC Chain Go SDK

## Description

Bnc-Go-SDK provides a thin wrapper around the BNC Chain API for readonly endpoints, in addition to creating and submitting different transactions.

## Usage

### Init

```GO
// chain api url
bnc, _ := sdk.NewBncSDK("http://localhost:8080/api/v1")
```

### Read Operations

```GO
account, _ := bnc.GetAccount("cosmosaccaddr1hy2e872rqtd675sn72ny87cyyaaanmqeuvwrpc")
fmt.Println("GetAccount: ", account)

depth, _ := bnc.GetDepth(&sdk.DepthQuery{Symbol: "BNB", Limit: 10})
fmt.Println("GetDepth: ", depth)

kline, _ := bnc.GetKlines(&sdk.KlineQuery{Symbol: "BNB", Interval: "1h"})
fmt.Println("GetKlines: ", kline)

markets, _ := bnc.GetMarkets(100)
fmt.Println("GetMarkets: ", markets)

order, _ := bnc.GetOrder("order-id-abcd")
fmt.Println("GetOrder: ", order)

openOrders, _ := bnc.GetOpenOrders(&sdk.OpenOrdersQuery{SenderAddress: "cosmosaccaddr1hy2e872rqtd675sn72ny87cyyaaanmqeuvwrpc", Symbol: "BNB"})
fmt.Println("GetOpenOrders: ", openOrders)

closedOrders, _ := bnc.GetClosedOrders(&sdk.ClosedOrdersQuery{SenderAddress: "cosmosaccaddr1hy2e872rqtd675sn72ny87cyyaaanmqeuvwrpc", Symbol: "BNB"})
fmt.Println("GetClosedOrders: ", closedOrders)

ticker24h, _ := bnc.GetTicker24h("BNB")
fmt.Println("GetTicker24h: ", ticker24h)

tokens, _ := bnc.GetTokens()
fmt.Println("GetTokens: ", tokens)

trades, _ := bnc.GetTrades(&sdk.TradesQuery{SenderAddress: "cosmosaccaddr1hy2e872rqtd675sn72ny87cyyaaanmqeuvwrpc", Symbol: "BNB"})
fmt.Println("GetTrades: ", trades)

tx, _ := bnc.GetTx("tx-hash-abcd")
fmt.Println("GetTx: ", tx)

time, _ := bnc.GetTime()
fmt.Println("GetTime: ", time)
```

### Create & Post Transaction

```GO
// generate new tx sequence
var sequence int64 = 1

// a testing utility to generate a private key and account
priv, acc := tx.PrivAndAddr()

// transaction type: [CreateOrderMsg, CancelOrderMsg, IssueToken, BurnToken, FreezeToken, UnfreezeToken or DexList]
newOrderMsg := txmsg.NewCreateOrderMsg(
  acc,
  txmsg.GenerateOrderID(sequence, acc),
  txmsg.OrderSide.BUY,
  "BNB_NNB",
  100000000,
  500000000,
)

// prepare message to sign
signMsg := tx.StdSignMsg{
  ChainID:       "bnc-chain-1",
  AccountNumber: 100,
  Sequence:      sequence,
  Memo:          "",
  Fee:           tx.NewStdFee(5000, tx.Coin{Denom: "BNB", Amount: 100000000}),
  Msgs:          []txmsg.Msg{newOrderMsg},
}

// Hex encoded signed transaction, ready to be posted to BncChain API
hexTx, _ := tx.Sign(priv.Bytes(), signMsg)
txResult, _ := bnc.PostTx(hexTx)
fmt.Println("PostTx: ", txResult)
```

### Transaction Types (Messages)

#### CreateOrderMsg

```GO
type CreateOrderMsg struct {
  Sender      AccAddress `json:"sender"`
  ID          string     `json:"id"`
  Symbol      string     `json:"symbol"`
  OrderType   int8       `json:"ordertype"`
  OrderSide   int8       `json:"side"`
  Price       int64      `json:"price"`
  Quantity    int64      `json:"quantity"`
  TimeInForce int8       `json:"timeinforce"`
}
```

#### CancelOrderMsg

```GO
type CancelOrderMsg struct {
  Sender AccAddress
  Symbol string `json:"symbol"`
  ID     string `json:"id"`
  RefID  string `json:"refid"`
}
```

#### TokenIssueMsg

```GO
type TokenIssueMsg struct {
  From        AccAddress `json:"from"`
  Name        string     `json:"name"`
  Symbol      string     `json:"symbol"`
  TotalSupply int64      `json:"total_supply"`
}
```

#### TokenBurnMsg

```GO
type TokenBurnMsg struct {
  From   AccAddress `json:"from"`
  Symbol string     `json:"symbol"`
  Amount int64      `json:"amount"`
}
```

#### TokenFreezeMsg

```GO
type TokenFreezeMsg struct {
  From   AccAddress `json:"from"`
  Symbol string     `json:"symbol"`
  Amount int64      `json:"amount"`
}
```

#### TokenUnfreezeMsg

```GO
type TokenUnfreezeMsg struct {
  From   AccAddress `json:"from"`
  Symbol string     `json:"symbol"`
  Amount int64      `json:"amount"`
}
```

#### DexListMsg

```GO
type DexListMsg struct {
  From             AccAddress `json:"from"`
  BaseAssetSymbol  string     `json:"base_asset_symbol"`
  QuoteAssetSymbol string     `json:"quote_asset_symbol"`
  InitPrice        int64      `json:"init_price"`
}
```