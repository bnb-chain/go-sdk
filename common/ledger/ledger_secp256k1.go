package ledger

import (
	"github.com/binance-chain/go-sdk/common/types"
	ledgergo "github.com/binance-chain/ledger-cosmos-go"
	"github.com/btcsuite/btcd/btcec"
	tmbtcec "github.com/tendermint/btcd/btcec"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var (
	// discoverLedger defines a function to be invoked at runtime for discovering
	// a connected Ledger device.
	DiscoverLedger discoverLedgerFn
)

type (

	// discoverLedgerFn defines a Ledger discovery function that returns a
	// connected device or an error upon failure. Its allows a method to avoid CGO
	// dependencies when Ledger support is potentially not enabled.
	discoverLedgerFn func() (LedgerSecp256k1, error)

	// DerivationPath represents a Ledger derivation path.
	DerivationPath []uint32

	// LedgerSecp256k1 reflects an interface a Ledger API must implement for
	// the SECP256K1 scheme.
	LedgerSecp256k1 interface {
		GetPublicKeySECP256K1([]uint32) ([]byte, error)
		ShowAddressSECP256K1([]uint32, string) error
		SignSECP256K1([]uint32, []byte) ([]byte, error)
		GetVersion() (*ledgergo.VersionInfo, error)
		Close() error
	}

	// PrivKeyLedgerSecp256k1 implements PrivKey, calling the ledger nano we
	// cache the PubKey from the first call to use it later.
	PrivKeyLedgerSecp256k1 struct {
		crypto.PrivKey
		pubkey secp256k1.PubKeySecp256k1
		path   DerivationPath
		ledger LedgerSecp256k1
	}
)

func GenLedgerSecp256k1Key(path DerivationPath, device LedgerSecp256k1) (*PrivKeyLedgerSecp256k1, error) {
	var pk secp256k1.PubKeySecp256k1
	pubkey, err := device.GetPublicKeySECP256K1(path)
	if err != nil {
		return nil, err
	}
	// re-serialize in the 33-byte compressed format
	cmp, err := btcec.ParsePubKey(pubkey[:], btcec.S256())
	if err != nil {
		return nil, err
	}
	copy(pk[:], cmp.SerializeCompressed())

	privKey := PrivKeyLedgerSecp256k1{path: path, ledger: device, pubkey: pk}
	return &privKey, nil
}

func (pkl PrivKeyLedgerSecp256k1) Bytes() []byte {
	return nil
}

func (pkl PrivKeyLedgerSecp256k1) ShowSignAddr() error {
	return pkl.ledger.ShowAddressSECP256K1(pkl.path, types.Network.Bech32Prefixes())
}

func (pkl PrivKeyLedgerSecp256k1) Sign(msg []byte) ([]byte, error) {
	err := pkl.ShowSignAddr()
	if err != nil {
		return nil, err
	}
	sig, err := pkl.ledger.SignSECP256K1(pkl.path, msg)
	if err != nil {
		return nil, err
	}
	return convertDERtoBER(sig)
}

func (pkl PrivKeyLedgerSecp256k1) PubKey() crypto.PubKey {
	return pkl.pubkey
}

func (pkl PrivKeyLedgerSecp256k1) Equals(other crypto.PrivKey) bool {
	if ledger, ok := other.(*PrivKeyLedgerSecp256k1); ok {
		return pkl.PubKey().Equals(ledger.PubKey())
	}

	return false
}

func convertDERtoBER(signatureDER []byte) ([]byte, error) {
	sigDER, err := btcec.ParseSignature(signatureDER[:], btcec.S256())
	if err != nil {
		return nil, err
	}
	sigBER := tmbtcec.Signature{R: sigDER.R, S: sigDER.S}
	return sigBER.Serialize(), nil
}
