# Changelog

## latest
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
