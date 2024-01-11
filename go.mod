module github.com/bnb-chain/go-sdk

go 1.17

require (
	github.com/bnb-chain/node v0.10.7
	github.com/btcsuite/btcd/btcec/v2 v2.3.2
	github.com/btcsuite/btcd/btcutil v1.1.3 // indirect; indirect	github.com/btcsuite/btcutil v0.0.0-20190425235716-9e5f4b9a998d
	github.com/cosmos/cosmos-sdk v0.25.0
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/gorilla/websocket v1.5.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.2
	github.com/tendermint/btcd v0.1.1
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.35.9
	github.com/zondax/ledger-cosmos-go v0.9.9
	golang.org/x/crypto v0.7.0
	gopkg.in/resty.v1 v1.12.0
)

require (
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/bnb-chain/ics23 v0.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cosmos/ledger-go v0.9.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/etcd-io/bbolt v1.3.3 // indirect
	github.com/ethereum/go-ethereum v1.11.3 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-kit/kit v0.10.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/herumi/bls-eth-go-binary v0.0.0-20210917013441-d37c07cfda4e // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/libp2p/go-buffer-pool v0.1.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pelletier/go-toml v1.9.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/prysmaticlabs/fastssz v0.0.0-20220628121656-93dfe28febab // indirect
	github.com/prysmaticlabs/gohashtree v0.0.3-alpha // indirect
	github.com/prysmaticlabs/prysm/v4 v4.0.2 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.8.1 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/supranational/blst v0.3.8-0.20220526154634-513d2456b344 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d // indirect
	github.com/tendermint/iavl v0.12.4 // indirect
	github.com/thomaso-mirodin/intmath v0.0.0-20160323211736-5dc6d854e46e // indirect
	github.com/zondax/hid v0.9.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20210624195500-8bfb893ecb84 // indirect
	google.golang.org/grpc v1.40.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/cosmos/cosmos-sdk => github.com/bnb-chain/bnc-cosmos-sdk v0.26.5-0.20240110102103-166bea30edf3
	github.com/tendermint/go-amino => github.com/bnb-chain/bnc-go-amino v0.14.1-binance.2
	github.com/tendermint/iavl => github.com/bnb-chain/bnc-tendermint-iavl v0.12.0-binance.5
	github.com/tendermint/tendermint => github.com/bnb-chain/bnc-tendermint v0.32.3-bc.10
	github.com/zondax/ledger-cosmos-go => github.com/bnb-chain/ledger-cosmos-go v0.9.10-0.20230201065744-d644bede1667
	github.com/zondax/ledger-go => github.com/bnb-chain/ledger-go v0.9.1
	golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20190823183015-45b1026d81ae
)
