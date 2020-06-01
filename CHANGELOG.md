# Changelog
## 1.2.3
CHAIN UPGRADE
* [\#110](https://github.com/binance-chain/go-sdk/pull/110) [RPC] [API] Add pending_match flag
* [\#130](https://github.com/binance-chain/go-sdk/pull/130) [RPC] [API] Support Mini Token
## 1.2.2
* [\#106](https://github.com/binance-chain/go-sdk/pull/106) [RPC] fix nil point error in getBalance rpc call
* [\#103](https://github.com/binance-chain/go-sdk/pull/103) [RPC] change the default timeout of RPC client as 5 seconds
* [\#102](https://github.com/binance-chain/go-sdk/pull/102) [FIX] Some typos only (managr/manger) 

## 1.2.1
* [\#99](https://github.com/binance-chain/go-sdk/pull/99) [BUILD] upgrade version of btcd to avoid retag issue 

## v1.2.0
* [\#93](https://github.com/binance-chain/go-sdk/pull/93) [BREAKING] uprade to binance chain release 0.6.3

## v1.1.3
* [\#81](https://github.com/binance-chain/go-sdk/pull/81) [TX] support swap on a single chain 


## v1.1.2
* [\#88](https://github.com/binance-chain/go-sdk/pull/88) [RPC] wrap error for abci query when abci code is not 0

## v1.1.1
IMPROVEMENT
* [\#87](https://github.com/binance-chain/go-sdk/pull/87) [RPC] distinguish not found error for get timelock rpc
* [\#84](https://github.com/binance-chain/go-sdk/pull/84) [RPC] change interface of get timelock


## v1.1.0
IMPROVEMENT
* [\#82](https://github.com/binance-chain/go-sdk/pull/82) [RPC] refactor reconnection

## v1.0.9

FEATURES
* [\#71](https://github.com/binance-chain/go-sdk/pull/71) [RPC] add timelock query support 
* [\#73](https://github.com/binance-chain/go-sdk/pull/73) [RPC] add limit param to get depth api for RPC


## v1.0.8
IMPROVEMENTS
* [\#53](https://github.com/binance-chain/go-sdk/pull/53) [SOURCE] change the default source into 0
* [\#56](https://github.com/binance-chain/go-sdk/pull/56) [RPC] add reconnect strategy when timeout to receive response
* [\#61](https://github.com/binance-chain/go-sdk/pull/61) [KEY] support bip44 to derive many address from same seed phase

FEATURES
* [\#66](https://github.com/binance-chain/go-sdk/pull/66)  [API]  support set account flag transaction
* [\#70](https://github.com/binance-chain/go-sdk/pull/70)  [API]  support atomic swap transactions

BREAKING
* [\#57](https://github.com/binance-chain/go-sdk/pull/57) [API] add query option to getTokens api
