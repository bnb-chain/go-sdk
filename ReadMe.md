# BNB Beacon Chain Go SDK

The BNB Beacon Chain GO SDK provides a thin wrapper around the BNB Beacon Chain API for readonly endpoints, in addition to creating and submitting different transactions.
It includes the following core components:

* **client** - implementations of BNB Beacon Chain transaction types and query, such as for transfers and trading.
* **common** - core cryptographic functions, uuid functions and other useful functions.
* **e2e** - end-to-end test package for go-sdk developer. For common users, it is also a good reference to use go-sdk. 
* **keys** - implement `KeyManage` to manage private key and accounts.
* **types** - core type of BNB Beacon Chain, such as `coin`, `account`, `tx` and `msg`.

## Install

### Requirement

Go version above 1.17

### Use go mod

Add "github.com/bnb-chain/go-sdk" dependency into your go.mod file. Example:
```go
require (
	github.com/bnb-chain/go-sdk latest
)

// Copy the same replace dep from https://github.com/bnb-chain/go-sdk/blob/master/go.mod
replace (
    github.com/cosmos/cosmos-sdk => github.com/bnb-chain/bnc-cosmos-sdk v0.25.4-0.20221221115251-f9e69ff1b273
    github.com/tendermint/go-amino => github.com/bnb-chain/bnc-go-amino v0.14.1-binance.2
    github.com/tendermint/iavl => github.com/bnb-chain/bnc-tendermint-iavl v0.12.0-binance.4
    github.com/tendermint/tendermint => github.com/bnb-chain/bnc-tendermint v0.32.3-binance.3.0.20221109023026-379ddbab19d1
    github.com/zondax/ledger-cosmos-go => github.com/bnb-chain/ledger-cosmos-go v0.9.9-binance.3
    github.com/zondax/ledger-go => github.com/bnb-chain/ledger-go v0.9.1
    golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20190823183015-45b1026d81ae
)
```

**NOTE**: Please make sure you have the same replace dep as [go.mod](https://github.com/bnb-chain/go-sdk/blob/master/go.mod).

## Usage 

### Key Manager

Before start using API, you should construct a Key Manager to help sign the transaction msg or verify signature.
Key Manager is an Identity Manager to define who you are in the bnbchain. It provide following interface:

```go
type KeyManager interface {
	Sign(tx.StdSignMsg) ([]byte, error)
	GetPrivKey() crypto.PrivKey
	GetAddr() txmsg.AccAddress
	
	ExportAsMnemonic() (string, error)
	ExportAsPrivateKey() (string, error)
	ExportAsKeyStore(password string) (*EncryptedKeyJSON, error)
}
```

We provide four construct functions to generate Key Manager:
```go
NewKeyManager() (KeyManager, error)

NewMnemonicKeyManager(mnemonic string) (KeyManager, error)

NewMnemonicPathKeyManager(mnemonic, keyPath string) (KeyManager, error) 

NewKeyStoreKeyManager(file string, auth string) (KeyManager, error)

NewPrivateKeyManager(priKey string) (KeyManager, error) 

NewLedgerKeyManager(path ledger.DerivationPath) (KeyManager, error)

```
- NewKeyManager. You will get a new private key without provide anything, you can export and save this `KeyManager`.
- NewMnemonicKeyManager. You should provide your mnemonic, usually is a string of 24 words.
- NewMnemonicPathKeyManager. The difference between `NewMnemonicKeyManager` is that you can use custom keypath to generate different `keyManager` while using the same mnemonic. 5 levels in BIP44 path: "purpose' / coin_type' / account' / change / address_index", "purpose' / coin_type'" is fixed as "44'/714'/", you can customize the rest part. 
- NewKeyStoreKeyManager. You should provide a keybase json file and you password, you can download the key base json file when your create a wallet account.
- NewPrivateKeyManager. You should provide a Hex encoded string of your private key.
- NewLedgerKeyManager. You must have a ledger device with BNB Beacon Chain ledger app and connect it to your machine.

Examples:

From mnemonic:
```Go
mnemonic := "lock globe panda armed mandate fabric couple dove climb step stove price recall decrease fire sail ring media enhance excite deny valid ceiling arm"
keyManager, _ := keys.NewMnemonicKeyManager(mnemonic)
```

From key base file:
```GO
file := "testkeystore.json"
keyManager, err := NewKeyStoreKeyManager(file, "your password")

```

From raw private key string:
```GO
priv := "9579fff0cab07a4379e845a890105004ba4c8276f8ad9d22082b2acbf02d884b"
keyManager, err := NewPrivateKeyManager(priv)
```

From ledger device:
```GO
bip44Params := keys.NewBinanceBIP44Params(0, 0)
keyManager, err := NewLedgerKeyManager(bip44Params.DerivationPath())
```

We provide three export functions to persistent a Key Manager:

```go
ExportAsMnemonic() (string, error)

ExportAsPrivateKey() (string, error)

ExportAsKeyStore(password string) (*EncryptedKeyJSON, error)
``` 

Examples:
```go
km, _ := NewKeyManager()
encryPlain1, _ := km.GetPrivKey().Sign([]byte("test plain"))
keyJSONV1, err := km.ExportAsKeyStore("testpassword")
bz, _ := json.Marshal(keyJSONV1)
ioutil.WriteFile("TestGenerateKeyStoreNoError.json", bz, 0660)
newkm, _ := NewKeyStoreKeyManager("TestGenerateKeyStoreNoError.json", "testpassword")
encryPlain2, _ := newkm.GetPrivKey().Sign([]byte("test plain"))
assert.True(t, bytes.Equal(encryPlain1, encryPlain2))
```
**As for ledger key, it can't be exported. Because its private key is saved on ledger device and no one can directly access it outside.** 

### Init Client

```GO
import sdk "https://github.com/bnb-chain/go-sdk/tree/master/client"

mnemonic := "lock globe panda armed mandate fabric couple dove climb step stove price recall decrease fire sail ring media enhance excite deny valid ceiling arm"
//-----   Init KeyManager  -------------
keyManager, _ := keys.NewMnemonicKeyManager(mnemonic)

//-----   Init sdk  -------------
client, err := sdk.NewDexClient("testnet-dex.binance.org", types.TestNetwork, keyManager)

```
For sdk init, you should know the famous api address. Besides, you should know what kind of network the api gateway is in, since we have different configurations for 
test network and production network.

|  ChainNetwork |  ApiAddr | 
|-------------- |----------------------------------|
|   TestNetwork | testnet-dex.binance.org  |  
|   ProdNetwork | dex.binance.org          |                                |

If you want broadcast some transactions, like send coins, create orders or cancel orders, you should construct a key manager.


### Example

Create a `SendToken` transaction: 
```go
client.SendToken([]msg.Transfer{{testAccount, []ctypes.Coin{{nativeSymbol, 100000000}}}}, true)
```

If want to attach memo or source to the transaction, more `WithSource` and `WithMemo` options are required:
```go
client.SendToken([]msg.Transfer{{testAccount, []ctypes.Coin{{nativeSymbol, 100000000}}}}, true, transaction.WithSource(100),transaction.WithMemo("test memo"))
```

In some scenarios, continuously send multi transactions very fast. Before the previous transaction being included in the chain, the next transaction is being sent, to avoid sequence mismatch error, option `WithAcNumAndSequence` is required:
```
acc,err:=client.GetAccount(client.GetKeyManager().GetAddr().String())
_, err = client.SendToken([]msg.Transfer{{testAccount, []ctypes.Coin{{nativeSymbol, 100000000}}}}, true, transaction.WithAcNumAndSequence(acc.Number,acc.Sequence))
_, err = client.SendToken([]msg.Transfer{{testAccount, []ctypes.Coin{{nativeSymbol, 100000000}}}}, true, transaction.WithAcNumAndSequence(acc.Number,acc.Sequence+1))
_, err = client.SendToken([]msg.Transfer{{testAccount, []ctypes.Coin{{nativeSymbol, 100000000}}}}, true, transaction.WithAcNumAndSequence(acc.Number,acc.Sequence+2))
```

For more API usage documentation, please check the [wiki](https://github.com/bnb-chain/go-sdk/wiki)..

## RPC Client
RPC endpoints may be used to interact with a node directly over HTTP or websockets. Using RPC, you may perform low-level 
operations like executing ABCI queries, viewing network/consensus state or broadcasting a transaction against full node or
light client.

### Example
```go
nodeAddr := "tcp://127.0.0.1:27147"
testClientInstance := rpc.NewRPCClient(nodeAddr,types.TestNetwork)
status, err := c.Status()
```
