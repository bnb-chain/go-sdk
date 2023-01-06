module github.com/bnb-chain/go-sdk

go 1.17

require (
	github.com/bnb-chain/node v0.10.6
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v0.0.0-20190425235716-9e5f4b9a998d
	github.com/cosmos/cosmos-sdk v0.25.0
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/gorilla/websocket v1.4.1-0.20190629185528-ae1634f6a989
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/btcd v0.1.1
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.32.3
	github.com/zondax/ledger-cosmos-go v0.9.9
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	gopkg.in/resty.v1 v1.12.0
)

require (
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/cosmos/ledger-go v0.9.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/etcd-io/bbolt v1.3.3 // indirect
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/go-kit/kit v0.9.0 // indirect
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515 // indirect
	github.com/libp2p/go-buffer-pool v0.0.2 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.1.0 // indirect
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90 // indirect
	github.com/prometheus/common v0.6.0 // indirect
	github.com/prometheus/procfs v0.0.3 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563 // indirect
	github.com/rs/cors v1.6.0 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/cobra v0.0.5 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/spf13/viper v1.4.0 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20190923125748-758128399b1d // indirect
	github.com/tendermint/iavl v0.12.4 // indirect
	github.com/zondax/hid v0.9.0 // indirect
	golang.org/x/net v0.3.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20190425155659-357c62f0e4bb // indirect
	google.golang.org/grpc v1.23.0 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/bnb-chain/bnc-cosmos-sdk v0.25.4-0.20221221115251-f9e69ff1b273
	github.com/tendermint/go-amino => github.com/bnb-chain/bnc-go-amino v0.14.1-binance.2
	github.com/tendermint/iavl => github.com/bnb-chain/bnc-tendermint-iavl v0.12.0-binance.4
	github.com/tendermint/tendermint => github.com/bnb-chain/bnc-tendermint v0.32.3-binance.3.0.20221109023026-379ddbab19d1
	github.com/zondax/ledger-cosmos-go => github.com/bnb-chain/ledger-cosmos-go v0.9.9-binance.3
	github.com/zondax/ledger-go => github.com/bnb-chain/ledger-go v0.9.1
	golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20190823183015-45b1026d81ae
)
