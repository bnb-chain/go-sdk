module github.com/binance-chain/go-sdk

require (
	github.com/binance-chain/ledger-cosmos-go v0.9.9-binance.1
	github.com/btcsuite/btcd v0.0.0-20181013004428-67e573d211ac
	github.com/btcsuite/btcutil v0.0.0-20180706230648-ab6388e0c60a
	github.com/coreos/go-iptables v0.4.0
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-kit/kit v0.8.0 // indirect
	github.com/go-logfmt/logfmt v0.3.0 // indirect
	github.com/gogo/protobuf v1.1.1 // indirect
	github.com/gorilla/websocket v1.4.0
	github.com/pkg/errors v0.8.0
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v0.9.2 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a
	github.com/rs/cors v1.6.0 // indirect
	github.com/stretchr/testify v1.2.2
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tendermint/btcd v0.0.0-20180816174608-e5840949ff4f
	github.com/tendermint/ed25519 v0.0.0-20171027050219-d8387025d2b9 // indirect
	github.com/tendermint/go-amino v0.14.1
	github.com/tendermint/tendermint v0.31.2-rc0
	github.com/zondax/hid v0.9.0 // indirect
	github.com/zondax/ledger-go v0.9.0 // indirect
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	google.golang.org/grpc v1.19.1 // indirect
	gopkg.in/resty.v1 v1.10.3
)

replace github.com/tendermint/go-amino => github.com/binance-chain/bnc-go-amino v0.14.1-binance.1

replace github.com/zondax/ledger-go => github.com/binance-chain/ledger-go v0.9.1
