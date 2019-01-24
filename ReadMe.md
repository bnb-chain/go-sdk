# BNC Chain Go SDK

## Description

Bnc-Go-SDK provides a thin wrapper around the BNC Chain API for readonly endpoints, in addition to creating and submitting different transactions.

## Install

### Use go mod(recommend)

Add "github.com/binance-chain/go-sdk" dependency into your go.mod file. Example:
```go
require (
	github.com/binance-chain/go-sdk latest
)
```

### Use go get

Use go get to install sdk into your `GOPATH`:
```bash
go get github.com/binance-chain/go-sdk
```

## Use dep
Add dependency to your Gopkg.toml file. Example:
```bash
[[override]]
  name = "github.com/binance-chain/go-sdk"
```

## API 

### Key Manager

Before start using API, you should construct a Key Manager to help sign the transaction msg or verify signature.
Key Manager is an Identity Manger to define who you are in the bnbchain. It provide following interface:

```go
type KeyManager interface {
	Sign(tx.StdSignMsg) ([]byte, error)
	GetPrivKey() crypto.PrivKey
	GetAddr() txmsg.AccAddress
}
```

We provide three construct functions to generate Key Manger:
```go
NewMnemonicKeyManager(mnemonic string) (KeyManager, error)

NewKeyStoreKeyManager(file string, auth string) (KeyManager, error)

NewPrivateKeyManager(wifKey string) (KeyManager, error) 
```

- NewMnemonicKeyManager. You should provide your mnemonic, usually is a string of 24 words.
- NewKeyStoreKeyManager. You should provide a keybase json file and you password, you can download the key base json file when your create a wallet account.
- NewPrivateKeyManager. You should provide a Hex encoded string of your private key.

Examples:

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



### Init Client

```GO
import sdk "github.com/binance-chain/go-sdk/client"

mnemonic := "lock globe panda armed mandate fabric couple dove climb step stove price recall decrease fire sail ring media enhance excite deny valid ceiling arm"
//-----   Init KeyManager  -------------
keyManager, _ := keys.NewMnemonicKeyManager(mnemonic)

//-----   Init sdk  -------------
client, err := sdk.NewDexClient("https://testnet-dex.binance.org", types.TestNetwork, keyManager)

```
For sdk init, you should know the famous api address. Besides, you should know what kind of network the api gateway is in, since we have different configurations for 
test network and production network.
TestNetwork ChainNetwork = iota

|  ChainNetwork |  ApiAddr | 
|-------------- |----------------------------------|
|   TestNetwork | https://testnet-dex.binance.org  |  
|   ProdNetwork | https://dex.binance.org          |                                |

If you want broadcast some transactions, like send coins, create orders or cancel orders, you should construct a key manager.


### Read Operations

#### Get Account

```GO
account, err := client.GetAccount("Your address")
```
##### Parameters

- Address - **string** , The address of query account.

##### Returns

- Account - The account object with the following structure:

	- Number    **int64** , The account number of this user, which is a global unique number.
	- Address   **string** , the address of this account, which is hash of public key.
	- Balances  **[]Coin** , the balances of different kind of tokens.
	- PublicKey **[]uint8** , the public key of this user.
	- Sequence  **int64** , the next expected transaction sequence, which is used to prevent replay accack.

#### Get Markets
```
markets, err := client.GetMarkets(api.NewMarketsQuery().WithLimit(1))
```

##### Parameters
     
- **MarketsQuery**, The query object.
	- Offset **\*uint32** , optional, the offset of the first return symbol pair.
	- Limit  **\*uint32** , optional, the max length of return symbol pair.
##### Returns

- **[]SymbolPair**
  - **SymbolPair**, 
    - TradeAsset **string**  
    - QuoteAsset **string**
    - Price      **string**, the price of trade assert against quote assert.
    - TickSize   **string**, the minimum price movement of a trading instrument.
    - LotSize    **string**, refers to the quantity of an item ordered for delivery on a specific date or manufactured in a single production run. 

#### Get Depth
```go
depth, err := client.GetDepth(api.NewDepthQuery(tradeSymbol, nativeSymbol))
```
##### Parameters
     
- **DepthQuery**, The query object.
  - Symbol **string**, the combination of trade symbol and quote symbol.
  - Limit  **\*uint32**, optional, the max length of return depth.

##### Returns
- **MarketDepth**
	- Bids   **[][]string**, each bid get two string element, the first one is the buy price, the second one is buy quantity. example: `[ [ "0.0024", "10" ] ]`.
	- Asks   **[][]string**, each ask get two string element, the first one is the sell price, the second one is sell quantity. example:` [ [ "0.0024", "10" ] ]`.
	- Height **int64**, the bids and asks is based on a certain height of the chain.
}

#### Get Kline
```go
kline, err := client.GetKlines(api.NewKlineQuery(tradeSymbol, nativeSymbol, "1h").WithLimit(1))

```
##### Parameters 
- **KlineQuery**, The query object.
  - Symbol    **string**  ,the combination of trade symbol and quote symbol.
  - Interval  **string**  interval like: (5m, 1h, 1d, 1w, etc.).
  - Limit     **\*uint32** , optional.
  - StartTime **\*int64** , optional, which is a nano time.
  - EndTime   **\*int64**  , optional, which is a nano time.


##### Returns
- **[]Kline** 
  - **Kline**
	- Close            **float64**, the close price .
	- CloseTime        **int64**, the close time.
	- High             **float64**, the highest price during the time.
	- Low              **float64**, the lowest price during the time.
	- NumberOfTrades   **int32**, the number of the trade transactions.
	- Open             **float64**, the open price.
	- OpenTime         **int64**,  the open time.
	- QuoteAssetVolume **float64**, the volume of the quote asset.
	- Volume           **float64**, the volume of trade asset.

#### Get Ticker 24h

```go
ticker24h, err := client.GetTicker24h(api.NewTicker24hQuery().WithSymbol(tradeSymbol, nativeSymbol))
```
##### Parameters 
- **Ticker24hQuery**, the query object.
  - Symbol **string**, the combination of trade symbol and quote symbol.

##### Returns
- **[]Ticker24h**
  - **Ticker24h**
  	- Symbol             **string** 
	- AskPrice           **string** , in decimal form, e.g. 1.00000000
	- AskQuantity        **string** in decimal form, e.g. 1.00000000
	- BidPrice           **string** in decimal form, e.g. 1.00000000
	- BidQuantity        **string** in decimal form, e.g. 1.00000000
	- CloseTime          **int64**  
	- Count              **int64**  
	- FirstID            **string** 
	- HighPrice          **string** in decimal form, e.g. 1.00000000
	- LastID             **string** 
	- LastPrice          **string** in decimal form, e.g. 1.00000000
	- LastQuantity       **string** in decimal form, e.g. 1.00000000
	- LowPrice           **string** in decimal form, e.g. 1.00000000
	- OpenPrice          **string** in decimal form, e.g. 1.00000000
	- OpenTime           **int64**  
	- PrevClosePrice     **string** in decimal form, e.g. 1.00000000
	- PriceChange        **string** in decimal form, e.g. 1.00000000
	- PriceChangePercent **string** 
	- QuoteVolume        **string** ,in decimal form, e.g. 1.00000000
	- Volume             **string** ,i n decimal form, e.g. 1.00000000
	- WeightedAvgPrice   **string** 

#### Get Tokens

```go
tokens, err := client.GetTokens()
```
##### Parameters 
- No parameters

##### Returns
- **[]Token**
  - **Token**
    - Name        **string**, the name of the token, which end with three random alphabets.
	- TotalSupply **string**, the total  supply of this token.
	- Owner       **string**, who issue this token, which is an address of an account.
	- OriginalSymbol **string**, the original symbol, which do not end with three random alphabet.
}

#### Get Trades
```go
trades, err := client.GetTrades(api.NewTradesQuery(testAccount1.String()).WithSymbol(tradeSymbol, nativeSymbol))
```
##### Parameters
- **TradesQuery**, the query object.
  - SenderAddress **string**, the address of the trade sender.
  - Symbol        **string**, the symbol of the trade, combination of trade symbol and quote symbol.
  -	Offset        **\*uint32**, optional.
  -	Limit         **\*uint32**, optional.
  -	Start         **\*int64**, optional.
  -	End           **\*int64**, optional.
  -	Side          **string**, the side of trades, options is ["BUY","SELl""].
##### Returns
- **Trades**
  - Trade **[]Trade**
    - BuyerOrderID  **string**, the order id of the buyer, which is combination of address and sequence.
    - BuyFee        **string**, the buy fee charged.
    - BuyerId       **string**, the buyer id.
    - Price         **string**, the trade price.
    - Quantity      **string**, the quantity of the trade.
    - SellFee       **string**,  the sell fee charged.
    - SellerId      **string**, the seller id.
    - SellerOrderID **string**, the order id of the buyer, which is combination of address and sequence.
    - Symbol        **string**, 
    - Time          **int64**, when the trade happened.
    - TradeID       **string**
    - BlockHeight   **int64**, in what height of the chain the trade happened.  
    - BaseAsset     **string**
    - QuoteAsset    **string**
  - Total **int**, the total num of trades.



#### Get Time
```go
time, err := client.GetTime()

```
##### Parameters
No parameters.

##### Returns
- **Time**
  - ApTime    **string**, the time of access point.
  - BlockTime **string**, the time of the block chain.


#### Get Order
```go
order, err := client.GetOrder("Your Order Id")
```
##### Parameters
- OrderId **string**, which is combination of account address and sequence.

##### Returns
- **Order**
  - ID                   **string**, the order id.
  -	Owner                **string**, the account address who set the order.
  -	Symbol               **string**, the combination of trade symbol and quote symbol.
  -	Price                **string**, the sell price or buy price.
  -	Quantity             **string**, the quantity of this order.
  -	CumulateQuantity     **string**, the total executed quantity.
  -	Fee                  **string**, the fee charged.
  -	Side                 **int**, 1 for buy and 2 for sell
  -	Status               **string**, options is [ ACK, PARTIALLY_FILLED, IOC_NO_FILL, FULLY_FILLED, CANCELED, EXPIRED, FAIL_BLOCKING, FAIL_MATCH, UNKNOWN ]
  -	TimeInForce          **int**, 1 for Good Till Expire(GTE) order and 3 for Immediate Or Cancel (IOC)
  -	Type                 **int**, only 2 is available for now, meaning limit order
  -	TradeId              **string**
  -	LastExecutedPrice    **string**, the price of last executed.
  -	LastExecutedQuantity **string**, the quantity of last execution.
  -	TransactionHash      **string** 
  -	TransactionTime      **string**


#### Get Open Orders
```go
openOrders, err := client.GetOpenOrders(api.NewOpenOrdersQuery(testAccount1.String()))

```
##### Parameters
- **OpenOrdersQuery**
  - SenderAddress **string**,the combination of trade symbol and quote symbol.
  - Symbol        **string**  
  - Offset        **\*uint32**, optional.
  - Limit         **\*uint32** , optional.
}
##### Returns
- **OpenOrders**
  - Order **[]Orde** 
  - Total **string**

#### Get Closed Orders

```go
closedOrders, err := client.GetClosedOrders(api.NewClosedOrdersQuery(testAccount1.String()).WithSymbol(tradeSymbol, nativeSymbol))
```
##### Parameters
- **ClosedOrdersQuery**
  - SenderAddress **string**,the combination of trade symbol and quote symbol.
  - Symbol        **string**  
  - Offset        **\*uint32**, optional.
  - Limit         **\*uint32** , optional.
}
##### Returns
- **OpenOrders**
  - Order **[]Orde** 
  - Total **string**

#### Get Tx
```go
tx, err := client.GetTx(openOrders.Order[0].TransactionHash)

```
##### Parameters
- TxHash **string**, the hash of the transaction.

##### Returns
- **TxResult**
  - Hash **string** 
  -	Log  **string**, log info if the transaction failed.
  -	Data **string**, the return result of different kind of transactions.
  -	Code **int32**, the result code of this transaction. Zero represent ok result.

#### Notice

For read option, each api will need a Query parameter. Each Query parameter we provide a construct function: `NewXXXQuery`.
We recommend you to use this construct function when you new a Query parameter since the construction only need required parameters,
and for the optional parameters, you can use `WithXXX` to add it.

### Create & Post Transaction

There is one most important point you should notice that we use int64 to represent a decimal.
The decimal length is fix 8, which means:
`100000000` is equal to `1`
`150000000` is equal to `1.5`
`1050000000` is equal to `10.5`

For a common transaction, the response is:

- **TxCommitResult**
  - Ok   **bool**, if the transaction accepted by chain.
  - Log  **string**, the error message of the transaction.
  - Hash **string**
  - Code **int32**, the result code. Zero represent fine.
  - Data **string**, different kind of transaction return different data message.

#### Create Order
```go
createOrderResult, err := client.CreateOrder(tradeSymbol, nativeSymbol, txmsg.OrderSide.BUY, 100000000, 100000000, true)
```
##### Parameters
- baseAssetSymbol **string** 
- quoteAssetSymbol **string**, 
- op **int8**, options is [1,2], 1 means "BUY", 2 means "SELL".
- price **int64**
- quantity **int64**
- sync **bool**, whether wait chain check this transaction. If true, in most case you will get `Data` field of `TxCommitResult`,
otherwise, `Data` field of `TxCommitResult` will be empty.

##### Return
- **CreateOrderResult**
  - **TxCommitResult**
  - OrderId **string**, the order id of this order.


#### Cancel Order
```go
cancelOrderResult, err := client.CancelOrder(tradeSymbol, nativeSymbol, createOrderResult.OrderId, createOrderResult.OrderId, true)
```
##### Parameters
- baseAssetSymbol **string**
- quoteAssetSymbol **string**
- id **string**, the order id.
- refId **string**, the order id will be fine.
- sync **bool**, whether wait chain check this transaction.

##### Return
- **CancelOrderResult**
  - **TxCommitResult**


#### Send token
```go
send, err := client.SendToken(testAccount2, nativeSymbol, 10000000000, true)
```
##### Parameters
- dest **txmsg.AccAddress**, the account address of user you want to send to.
- symbol **string**, the combination of trade symbol and quote symbol.
- quantity **int64** 
- sync **bool**, whether wait chain check this transaction.

##### Return
- **SendTokenResult**
  - **TxCommitResult**

#### Freeze token
```go
freeze, err := client.FreezeToken(nativeSymbol, 100000000, true)
```
##### Parameters
- symbol **string**, which kind of token you want to freeze.
- amount **int64** 
- sync **bool**, whether wait chain check this transaction.

##### Return
- **FreezeTokenResult**
  - **TxCommitResult**


#### UnFreeze token
```go
unFreeze, err := client.UnfreezeToken(nativeSymbol, 100000000, true)
```
##### Parameters
- symbol **string**, which kind of token you want to unfreeze.
- amount **int64** 
- sync **bool**, whether wait chain check this transaction.

##### Return
- **UnfreezeTokenResult**
  - **TxCommitResult**
  
#### Issue token
```go
issue, err := client.IssueToken("SDK-Token", "sdk", 10000000000000000, true, false)
```
##### Parameters
- name **string**, the name of your token.
- symbol **string**, a symbol of your token.
- supply **int64** 
- sync **bool**, whether wait chain check this transaction.
- mintable **bool**, whether you want mint token in the future.

##### Return
- **IssueTokenResult**
  - **TxCommitResult**
  - Symbol **string**, the actual symbol of you token.(which will end with three random alphabets).


#### Submit ListTrade Proposal
```go
listTradingProposal, err := client.SubmitListPairProposal("New trading pair", txmsg.ListTradingPairParams{issue.Symbol, nativeSymbol, 1000000000, "my trade", time2.Now().Add(1 * time2.Hour)}, 200000000000, true)
```
##### Parameters
- title **string**, 
- param **txmsg.ListTradingPairParams**
  - BaseAssetSymbol  **string**   
  - QuoteAssetSymbol **string**   
  - InitPrice        **int64**    
  - Description      **string**  
  - ExpireTime       **time.Time**, the expire time you active this list trading pair. 
- initialDeposit **int64**, the amount of BNB you want deposit for this proposal.
- sync **bool**, whether wait chain check this transaction.
##### Return
- **SubmitProposalResult**
  - *TxCommitResult*
  - ProposalId *int64*, the proposal id generated by chain. Useful when you vote or deposit for specified proposal.

#### Vote Proposal
```go
vote, err := client.VoteProposal(listTradingProposal.ProposalId, txmsg.OptionYes, true)
```
##### Parameters
- proposalID **int64**, the id of the proposal you want to vote.
- option **txmsg.VoteOption**, vote options: [OptionYes, OptionAbstain, OptionNo,O ptionNoWithVeto]
- sync bool, whether wait chain check this transaction.

##### Return
- **VoteProposalResult**
  - **TxCommitResult**
  
#### List Trade pair
```go
lp,err:=client.ListPair(listTradingProposal.ProposalId,  issue.Symbol,  nativeSymbol, 1000000000, true)
```
##### Parameters
- proposalId **int64**, the proposal id that propose by you or others that want to create a new list trading pair. Make sure the proposal is passed.
- baseAssetSymbol **string**
- quoteAssetSymbol **string** 
- initPrice **int64**
- sync **bool**,  whether wait chain check this transaction.

##### Return
- **ListPairResult**
  - **TxCommitResult**